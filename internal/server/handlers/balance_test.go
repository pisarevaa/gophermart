package handlers_test

import (
	"encoding/json"
	"net/http/httptest"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/pisarevaa/gophermart/internal/server"
	"github.com/pisarevaa/gophermart/internal/server/handlers"
	mock "github.com/pisarevaa/gophermart/internal/server/mocks"
	"github.com/pisarevaa/gophermart/internal/server/storage"
)

func (suite *ServerTestSuite) TestGetBalance() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	user := storage.User{
		Login:     "test",
		Password:  "123",
		Balance:   int64(500),
		Withdrawn: int64(300),
	}

	m.EXPECT().
		GetUser(gomock.Any(), gomock.Any()).
		Return(user, nil)

	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, m))
	defer ts.Close()

	resp, bodyResp := MakeAuthRequest(suite, ts, "GET", "/api/user/balance", nil, false, suite.token)

	var userBalanceResponse handlers.UserBalanceInfo

	err := json.Unmarshal([]byte(bodyResp), &userBalanceResponse)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	suite.Require().Equal(200, resp.StatusCode)
	suite.Require().Equal(int64(500), userBalanceResponse.Current)
	suite.Require().Equal(int64(300), userBalanceResponse.Withdrawn)
}

func (suite *ServerTestSuite) TestWithdrawBalance() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)
	tx := mock.NewMockTransaction(ctrl)

	withdraw := handlers.Withdraw{
		Order: "123",
		Sum:   int64(200),
	}
	withdrawJson, err := json.Marshal(withdraw)
	suite.Require().NoError(err)

	user := storage.User{
		Login:    "test",
		Password: "123",
		Balance:  int64(500),
	}

	order := storage.Order{
		Number:     "123",
		Status:     "NEW",
		Accrual:    int64(0),
		Login:      "test",
		UploadedAt: time.Now(),
	}

	m.EXPECT().
		BeginTransaction(gomock.Any()).
		Return(tx, nil)

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

	resp, _ := MakeAuthRequest(suite, ts, "POST", "/api/user/balance/withdraw", withdrawJson, true, suite.token)

	defer resp.Body.Close()
	suite.Require().Equal(200, resp.StatusCode)
}

func (suite *ServerTestSuite) TestWithdrawls() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	orders := []storage.Order{{
		Number:     "123",
		Status:     "NEW",
		Accrual:    int64(0),
		Login:      "test",
		UploadedAt: time.Now(),
	}}

	m.EXPECT().
		GetOrders(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(orders, nil)

	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, m))
	defer ts.Close()

	resp, bodyResp := MakeAuthRequest(suite, ts, "GET", "/api/user/withdrawals", nil, true, suite.token)

	var withdrawalsResponse []handlers.WithdrawalsReponse

	err := json.Unmarshal([]byte(bodyResp), &withdrawalsResponse)
	suite.Require().NoError(err)

	defer resp.Body.Close()
	suite.Require().Equal(200, resp.StatusCode)
	suite.Require().Len(withdrawalsResponse, 1)
}
