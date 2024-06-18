package handlers_test

import (
	"net/http/httptest"
	"time"

	"github.com/golang/mock/gomock"

	server "github.com/pisarevaa/gophermart/internal"
	"github.com/pisarevaa/gophermart/internal/handlers"
	mock "github.com/pisarevaa/gophermart/internal/mocks"
	"github.com/pisarevaa/gophermart/internal/storage"
)

type WithdrawalsReponse struct {
	Order       string    `json:"order"        binding:"required"`
	Sum         float32   `json:"sum"          binding:"required"`
	ProcessedAt time.Time `json:"processed_at" binding:"required"`
}

func (suite *ServerTestSuite) TestGetBalanceMockDB() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	user := storage.User{
		Login:     "test",
		Password:  "123",
		Balance:   float32(500),
		Withdrawn: float32(300),
	}

	m.EXPECT().
		GetUser(gomock.Any(), gomock.Any()).
		Return(user, nil)

	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, m))
	defer ts.Close()

	var userBalanceResponse handlers.UserBalanceInfo
	resp, err := suite.client.R().
		SetResult(&userBalanceResponse).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+suite.token).
		Get(ts.URL + "/api/user/balance")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
	suite.Require().Equal(float32(300), userBalanceResponse.Withdrawn)
	suite.Require().Equal(float32(500), userBalanceResponse.Current)
}

func (suite *ServerTestSuite) TestWithdrawBalanceMockDB() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)
	tx := mock.NewMockTransaction(ctrl)

	withdraw := handlers.Withdraw{
		Order: "123",
		Sum:   float32(200),
	}

	user := storage.User{
		Login:    "test",
		Password: "123",
		Balance:  float32(500),
	}

	order := storage.Order{
		Number:     "123",
		Status:     "PROCESSED",
		Accrual:    float32(100),
		Login:      "test",
		UploadedAt: time.Now(),
	}

	m.EXPECT().
		BeginTransaction(gomock.Any()).
		Return(tx, nil)

	tx.EXPECT().
		Rollback(gomock.Any()).
		Return(nil)

	tx.EXPECT().GetUserWithLock(gomock.Any(), gomock.Any()).
		Return(user, nil)

	tx.EXPECT().
		GetOrderWithLock(gomock.Any(), gomock.Any()).
		Return(order, nil)

	tx.EXPECT().
		WithdrawUserBalance(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	tx.EXPECT().
		WithdrawOrderBalance(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	tx.EXPECT().
		Commit(gomock.Any()).
		Return(nil)

	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, m))
	defer ts.Close()

	resp, err := suite.client.R().
		SetBody(withdraw).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+suite.token).
		Post(ts.URL + "/api/user/balance/withdraw")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
}

func (suite *ServerTestSuite) TestWithdrawlsMockDB() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)
	now := time.Now()
	orders := []storage.Order{{
		Number:      "123",
		Status:      "PROCESSED",
		Accrual:     float32(100),
		Login:       "test",
		UploadedAt:  now,
		ProcessedAt: &now,
	}}

	m.EXPECT().
		GetOrders(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(orders, nil)

	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, m))
	defer ts.Close()

	var withdrawalsResponse []WithdrawalsReponse
	resp, err := suite.client.R().
		SetResult(&withdrawalsResponse).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+suite.token).
		Get(ts.URL + "/api/user/withdrawals")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
	suite.Require().Len(withdrawalsResponse, 1)
}

func (suite *ServerTestSuite) TestGetBalanceInMemory() {
	m := storage.NewMemory()

	m.Users[login] = storage.User{
		Login:     "test",
		Password:  "123",
		Balance:   float32(500),
		Withdrawn: float32(300),
	}

	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, m))
	defer ts.Close()

	var userBalanceResponse handlers.UserBalanceInfo
	resp, err := suite.client.R().
		SetResult(&userBalanceResponse).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+suite.token).
		Get(ts.URL + "/api/user/balance")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
	suite.Require().Equal(float32(300), userBalanceResponse.Withdrawn)
	suite.Require().Equal(float32(500), userBalanceResponse.Current)
}

func (suite *ServerTestSuite) TestWithdrawInMemory() {
	m := storage.NewMemory()

	now := time.Now()
	m.Orders["123"] = storage.Order{
		Number:      "123",
		Status:      "PROCESSED",
		Accrual:     float32(100),
		Withdrawn:   float32(50),
		Login:       "test",
		UploadedAt:  now,
		ProcessedAt: &now,
	}

	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, m))
	defer ts.Close()

	var withdrawalsResponse []WithdrawalsReponse
	resp, err := suite.client.R().
		SetResult(&withdrawalsResponse).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+suite.token).
		Get(ts.URL + "/api/user/withdrawals")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
	suite.Require().Len(withdrawalsResponse, 1)
}

func (suite *ServerTestSuite) TestWithdrawBalanceInMemory() {
	m := storage.NewMemory()

	m.Users[login] = storage.User{
		Login:    "test",
		Password: "123",
		Balance:  float32(500),
	}

	now := time.Now()
	m.Orders["123"] = storage.Order{
		Number:      "123",
		Status:      "PROCESSED",
		Accrual:     float32(100),
		Login:       "test",
		UploadedAt:  now,
		ProcessedAt: &now,
	}

	withdraw := handlers.Withdraw{
		Order: "123",
		Sum:   float32(200),
	}

	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, m))
	defer ts.Close()

	resp, err := suite.client.R().
		SetBody(withdraw).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+suite.token).
		Post(ts.URL + "/api/user/balance/withdraw")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
}
