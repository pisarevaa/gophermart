package integration_test

import (
	"context"
	"errors"
	"math/rand/v2"
	"net/http/httptest"
	"os"
	"os/signal"
	"strconv"
	"testing"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/pisarevaa/gophermart/internal/server"
	"github.com/pisarevaa/gophermart/internal/server/configs"
	"github.com/pisarevaa/gophermart/internal/server/handlers"
	"github.com/pisarevaa/gophermart/internal/server/storage"
	"github.com/pisarevaa/gophermart/internal/server/tasks"
)

type OrderReponse struct {
	Number     string    `json:"number"     binding:"required"`
	Status     string    `json:"status"     binding:"required"`
	Accrual    float32   `json:"accrual"    binding:"required"`
	UploadedAt time.Time `json:"uploadedAt" binding:"required"`
}

type WithdrawalsReponse struct {
	Order       string    `json:"order"        binding:"required"`
	Sum         float32   `json:"sum"          binding:"required"`
	ProcessedAt time.Time `json:"processed_at" binding:"required"`
}

type Good struct {
	Description string `json:"description" binding:"required"`
	Price       int64  `json:"price"       binding:"required"`
}

type AddAccrualOrder struct {
	Order string `json:"order" binding:"required"`
	Goods []Good `json:"goods" binding:"required"`
}

type RewardSchema struct {
	Match      string `json:"match"       binding:"required"`
	Reward     int64  `json:"reward"      binding:"required"`
	RewardType string `json:"reward_type" binding:"required"`
}

type ServerTestSuite struct {
	suite.Suite
	cfg    configs.Config
	logger *zap.SugaredLogger
	repo   storage.Storage
	client *resty.Client
}

func (suite *ServerTestSuite) SetupSuite() {
	suite.cfg = configs.NewConfig()
	suite.logger = server.NewLogger()
	suite.repo = storage.NewDB(suite.cfg.DatabaseURI, suite.logger)
	suite.client = resty.New()
}

func TestAgentSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (suite *ServerTestSuite) TestFullProccess() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exit, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, suite.repo))
	defer ts.Close()

	task := tasks.NewTask(suite.cfg, suite.logger, suite.repo, suite.client)
	go task.RunUpdateOrderStatuses(exit)

	password := "123"
	randomLogin := "test_user_" + strconv.Itoa(rand.IntN(100000))
	user := storage.RegisterUser{
		Login:    randomLogin,
		Password: password,
	}

	suite.Run("Регистрация пользователя", func() {
		var successRegister storage.Success
		resp, err := suite.client.R().
			SetResult(&successRegister).
			SetBody(user).
			SetHeader("Content-Type", "application/json").
			Post(ts.URL + "/api/user/register")
		suite.Require().NoError(err)
		suite.Require().Equal(200, resp.StatusCode())
		suite.Require().True(successRegister.Success)
	})

	var successLogin handlers.SuccessLogin

	suite.Run("Логин пользователя", func() {
		resp, err := suite.client.R().
			SetResult(&successLogin).
			SetBody(user).
			SetHeader("Content-Type", "application/json").
			Post(ts.URL + "/api/user/login")
		suite.Require().NoError(err)
		suite.Require().Equal(200, resp.StatusCode())
		suite.Require().True(successLogin.Success)
	})

	number := goluhn.Generate(9)

	suite.Run("Загрузка заказа", func() {
		var successAddOrder storage.Success
		resp, err := suite.client.R().
			SetResult(&successAddOrder).
			SetBody(number).
			SetHeader("Content-Type", "text/plain").
			SetHeader("Authorization", "Bearer "+successLogin.Token).
			Post(ts.URL + "/api/user/orders")
		suite.Require().NoError(err)
		suite.Require().Equal(202, resp.StatusCode())
		suite.Require().True(successAddOrder.Success)
	})

	match := "match_" + strconv.Itoa(rand.IntN(100000))

	suite.Run("Загрузка схемы вознаграждения в сервис Accraul", func() {
		rewardSchema := RewardSchema{
			Match:      match,
			Reward:     20,
			RewardType: "%",
		}
		resp, err := suite.client.R().
			SetBody(rewardSchema).
			SetHeader("Content-Type", "application/json").
			Post(suite.cfg.AccrualSystemAddress + "/api/goods")
		suite.Require().NoError(err)
		suite.Require().Equal(200, resp.StatusCode())
	})

	suite.Run("Загрузка заказа в сервис Accraul", func() {
		accraulOrder := AddAccrualOrder{
			Order: number,
			Goods: []Good{{
				Description: match,
				Price:       100,
			}},
		}
		resp, err := suite.client.R().
			SetBody(accraulOrder).
			SetHeader("Content-Type", "application/json").
			Post(suite.cfg.AccrualSystemAddress + "/api/orders")
		suite.Require().NoError(err)
		suite.Require().Equal(202, resp.StatusCode())
	})

	var orderAccrual float32

	suite.Run("Получение списка заказов и проверка что заказ успешен", func() {
		ticker := time.NewTicker(time.Duration(5) * time.Second)
		defer ticker.Stop()
		var attempt int64
		var maxAttempts int64 = 20
		for {
			<-ticker.C
			var orders []OrderReponse
			resp, err := suite.client.R().
				SetResult(&orders).
				SetHeader("Content-Type", "application/json").
				SetHeader("Authorization", "Bearer "+successLogin.Token).
				Get(ts.URL + "/api/user/orders")
			suite.Require().NoError(err)
			suite.Require().Equal(200, resp.StatusCode())
			suite.Require().Len(orders, 1)
			suite.Require().Equal(orders[0].Number, number)
			status := orders[0].Status
			if attempt > maxAttempts {
				suite.Require().NoError(errors.New("too many attempts"))
				break
			}
			if status == "NEW" || status == "PROCESSING" || status == "REGISTERED" {
				attempt++
				continue
			}
			suite.Require().Positive(orders[0].Accrual)
			orderAccrual = orders[0].Accrual
			break
		}
	})

	sumToWidraw := float32(rand.Int64N(int64(orderAccrual)))

	suite.Run("Списание средств со счета пользователя", func() {
		withdrawOrder := handlers.Withdraw{
			Order: number,
			Sum:   sumToWidraw,
		}
		var successWithdraw storage.Success
		resp, err := suite.client.R().
			SetResult(&successWithdraw).
			SetBody(withdrawOrder).
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", "Bearer "+successLogin.Token).
			Post(ts.URL + "/api/user/balance/withdraw")
		suite.Require().NoError(err)
		suite.Require().Equal(200, resp.StatusCode())
		suite.Require().True(successWithdraw.Success)
	})

	suite.Run("Проверка баланса пользователя", func() {
		var userBalance handlers.UserBalanceInfo
		resp, err := suite.client.R().
			SetResult(&userBalance).
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", "Bearer "+successLogin.Token).
			Get(ts.URL + "/api/user/balance")
		suite.Require().NoError(err)
		suite.Require().Equal(200, resp.StatusCode())
		suite.Require().Equal(int64(sumToWidraw), int64(userBalance.Withdrawn))
		suite.Require().Equal(int64(orderAccrual-sumToWidraw), int64(userBalance.Current))
	})

	suite.Run("Проверка баланса пользователя", func() {
		var withdrawals []WithdrawalsReponse
		resp, err := suite.client.R().
			SetResult(&withdrawals).
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", "Bearer "+successLogin.Token).
			Get(ts.URL + "/api/user/withdrawals")
		suite.Require().NoError(err)
		suite.Require().Equal(200, resp.StatusCode())
		suite.Require().Len(withdrawals, 1)
		suite.Require().Equal(withdrawals[0].Order, number)
		suite.Require().Equal(int64(withdrawals[0].Sum), int64(sumToWidraw))
	})
}
