package integration_test

import (
	"context"
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
	Accrual    int64     `json:"accrual"    binding:"required"`
	UploadedAt time.Time `json:"uploadedAt" binding:"required"`
}

type WithdrawalsReponse struct {
	Order       string    `json:"order"        binding:"required"`
	Sum         int64     `json:"sum"          binding:"required"`
	ProcessedAt time.Time `json:"processed_at" binding:"required"`
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
	suite.repo = storage.NewDB(suite.cfg.DatabaseUri, suite.logger)
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

	// Регистрация пользователя
	password := "123"
	randomLogin := "test_user_" + strconv.Itoa(rand.IntN(100000))
	suite.logger.Info("randomLogin: ", randomLogin)
	user := storage.RegisterUser{
		Login:    randomLogin,
		Password: password,
	}
	var successRegister storage.Success
	resp, err := suite.client.R().
		SetResult(&successRegister).
		SetBody(user).
		SetHeader("Content-Type", "application/json").
		Post(ts.URL + "/api/user/register")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
	suite.Require().True(successRegister.Success)

	// Логин пользователя
	var successLogin handlers.SuccessLogin
	resp, err = suite.client.R().
		SetResult(&successLogin).
		SetBody(user).
		SetHeader("Content-Type", "application/json").
		Post(ts.URL + "/api/user/login")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
	suite.Require().True(successLogin.Success)

	// Загрузка заказа
	number := goluhn.Generate(9)
	suite.logger.Info("randomOrder: ", number)
	var successAddOrder storage.Success
	resp, err = suite.client.R().
		SetResult(&successAddOrder).
		SetBody(number).
		SetHeader("Content-Type", "text/plain").
		SetHeader("Authorization", "Bearer "+successLogin.Token).
		Post(ts.URL + "/api/user/orders")
	suite.Require().NoError(err)
	suite.Require().Equal(202, resp.StatusCode())
	suite.Require().True(successAddOrder.Success)

	// Получение списка заказов и проверка что заказ успешен
	var orderAccrual int64
	ticker := time.NewTicker(time.Duration(5) * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		var orders []OrderReponse
		resp, err = suite.client.R().
			SetResult(&orders).
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", "Bearer "+successLogin.Token).
			Get(ts.URL + "/api/user/orders")
		suite.Require().NoError(err)
		suite.Require().Equal(200, resp.StatusCode())
		suite.Require().Len(orders, 1)
		suite.Require().Equal(orders[0].Number, number)
		status := orders[0].Status
		if status == "NEW" || status == "PROCESSING" || status == "REGISTERED" {
			continue
		}
		suite.Require().Positive(orders[0].Accrual)
		orderAccrual = orders[0].Accrual
		break
	}

	// Списание средств со счета пользователя
	sumToWidraw := rand.Int64N(orderAccrual)
	suite.logger.Info("sumToWidraw: ", sumToWidraw)
	withdrawOrder := handlers.Withdraw{
		Order: number,
		Sum:   sumToWidraw,
	}
	var successWithdraw storage.Success
	resp, err = suite.client.R().
		SetResult(&successWithdraw).
		SetBody(withdrawOrder).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+successLogin.Token).
		Post(ts.URL + "/api/user/balance/withdraw")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
	suite.Require().True(successWithdraw.Success)

	// Проверка баланса пользователя
	var userBalance handlers.UserBalanceInfo
	resp, err = suite.client.R().
		SetResult(&userBalance).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+successLogin.Token).
		Get(ts.URL + "/api/user/balance")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
	suite.Require().Equal(sumToWidraw, userBalance.Withdrawn)
	suite.Require().Equal(orderAccrual-sumToWidraw, userBalance.Current)

	// Список заказов со списанием
	var withdrawals []WithdrawalsReponse
	resp, err = suite.client.R().
		SetResult(&withdrawals).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+successLogin.Token).
		Get(ts.URL + "/api/user/withdrawals")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
	suite.Require().Len(withdrawals, 1)
	suite.Require().Equal(withdrawals[0].Order, number)
	suite.Require().Equal(withdrawals[0].Sum, sumToWidraw)
}
