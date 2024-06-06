package handlers_test

import (
	"encoding/json"
	"net/http/httptest"

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
		Login:    "test",
		Password: "123",
		Balance:  int64(500),
	}

	m.EXPECT().
		GetUser(gomock.Any(), gomock.Any()).
		Return(user, nil)

	m.EXPECT().
		GetUserWithdrawals(gomock.Any(), gomock.Any()).
		Return(int64(300), nil)

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
