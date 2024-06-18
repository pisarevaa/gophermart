package handlers_test

import (
	"errors"
	"net/http/httptest"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/golang/mock/gomock"

	"github.com/pisarevaa/gophermart/internal/server"
	mock "github.com/pisarevaa/gophermart/internal/server/mocks"
	"github.com/pisarevaa/gophermart/internal/server/storage"
)

type OrderReponse struct {
	Number     string    `json:"number"     binding:"required"`
	Status     string    `json:"status"     binding:"required"`
	Accrual    int64     `json:"accrual"    binding:"required"`
	UploadedAt time.Time `json:"uploadedAt" binding:"required"`
}

func (suite *ServerTestSuite) TestAddOrderMockDB() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	m.EXPECT().
		GetOrder(gomock.Any(), gomock.Any()).
		Return(storage.Order{}, errors.New("not found"))

	m.EXPECT().
		StoreOrder(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	number := goluhn.Generate(9)

	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, m))
	defer ts.Close()

	resp, err := suite.client.R().
		SetBody(number).
		SetHeader("Content-Type", "text/plain").
		SetHeader("Authorization", "Bearer "+suite.token).
		Post(ts.URL + "/api/user/orders")
	suite.Require().NoError(err)
	suite.Require().Equal(202, resp.StatusCode())
}

func (suite *ServerTestSuite) TestGetOrdersMockDB() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	number := goluhn.Generate(9)

	orders := []storage.Order{{
		Number:     number,
		Status:     "NEW",
		Accrual:    0,
		Login:      login,
		UploadedAt: time.Now(),
	}}

	m.EXPECT().
		GetOrders(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(orders, nil)

	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, m))
	defer ts.Close()

	var ordersResponse []OrderReponse
	resp, err := suite.client.R().
		SetResult(&ordersResponse).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+suite.token).
		Get(ts.URL + "/api/user/orders")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
	suite.Require().Len(orders, 1)
}

func (suite *ServerTestSuite) TestAddAndGetOrdersInMemory() {
	m := storage.NewMemory()

	number := goluhn.Generate(9)

	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, m))
	defer ts.Close()

	resp, err := suite.client.R().
		SetBody(number).
		SetHeader("Content-Type", "text/plain").
		SetHeader("Authorization", "Bearer "+suite.token).
		Post(ts.URL + "/api/user/orders")
	suite.Require().NoError(err)
	suite.Require().Equal(202, resp.StatusCode())

	var ordersResponse []OrderReponse
	resp, err = suite.client.R().
		SetResult(&ordersResponse).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+suite.token).
		Get(ts.URL + "/api/user/orders")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
	suite.Require().Len(ordersResponse, 1)
}
