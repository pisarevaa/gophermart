package handlers_test

import (
	"net/http/httptest"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/pisarevaa/gophermart/internal/server"
	"github.com/pisarevaa/gophermart/internal/server/handlers"
	mock "github.com/pisarevaa/gophermart/internal/server/mocks"
	"github.com/pisarevaa/gophermart/internal/server/storage"
)

type WithdrawalsReponse struct {
	Order       string    `json:"order"        binding:"required"`
	Sum         float32   `json:"sum"          binding:"required"`
	ProcessedAt time.Time `json:"processed_at" binding:"required"`
}

func (suite *ServerTestSuite) TestGetBalance() {
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

// Пока не понимаю как замокать DBTransaction, вернусь позже
// func (suite *ServerTestSuite) TestWithdrawBalance() {
// 	ctrl := gomock.NewController(suite.T())
// 	defer ctrl.Finish()

// 	m := mock.NewMockStorage(ctrl)
// 	tx := mock.NewMockTransaction(ctrl)

// 	withdraw := handlers.Withdraw{
// 		Order: "123",
// 		Sum:   float32(200),
// 	}

// 	user := storage.User{
// 		Login:    "test",
// 		Password: "123",
// 		Balance:  float32(500),
// 	}

// 	order := storage.Order{
// 		Number:     "123",
// 		Status:     "PROCESSED",
// 		Accrual:    float32(100),
// 		Login:      "test",
// 		UploadedAt: time.Now(),
// 	}

// 	transaction := &storage.DBTransaction{}

// 	m.EXPECT().
// 		BeginTransaction(gomock.Any()).
// 		Return(transaction, nil)

// 	tx.EXPECT().GetUserWithLock(gomock.Any(), gomock.Any()).
// 		Return(user, nil)

// 	tx.EXPECT().
// 		GetOrderWithLock(gomock.Any(), gomock.Any()).
// 		Return(order, nil)

// 	tx.EXPECT().
// 		WithdrawUserBalance(gomock.Any(), gomock.Any(), gomock.Any()).
// 		Return(nil)

// 	tx.EXPECT().
// 		WithdrawOrderBalance(gomock.Any(), gomock.Any(), gomock.Any()).
// 		Return(nil)

// 	tx.EXPECT().
// 		Commit(gomock.Any()).
// 		Return(nil)

// 	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, m))
// 	defer ts.Close()

// 	resp, err := suite.client.R().
// 		SetBody(withdraw).
// 		SetHeader("Content-Type", "application/json").
// 		SetHeader("Authorization", "Bearer "+suite.token).
// 		Post(ts.URL + "/api/user/balance/withdraw")
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(200, resp.StatusCode())
// }

func (suite *ServerTestSuite) TestWithdrawls() {
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
