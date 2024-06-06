package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/golang/mock/gomock"

	"github.com/pisarevaa/gophermart/internal/server"
	"github.com/pisarevaa/gophermart/internal/server/handlers"
	mock "github.com/pisarevaa/gophermart/internal/server/mocks"
	"github.com/pisarevaa/gophermart/internal/server/storage"
)

func MakeAuthRequest(
	suite *ServerTestSuite,
	ts *httptest.Server,
	method string,
	url string,
	body []byte,
	isJson bool,
	token string,
) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+url, bytes.NewBuffer(body))
	suite.Require().NoError(err)
	if isJson {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "text/plain")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept-Encoding", "")
	resp, err := ts.Client().Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	suite.Require().NoError(err)
	return resp, string(respBody)
}

func (suite *ServerTestSuite) TestAddOrder() {
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

	resp, _ := MakeAuthRequest(suite, ts, "POST", "/api/user/orders", []byte(number), false, suite.token)

	defer resp.Body.Close()
	suite.Require().Equal(202, resp.StatusCode)
}

func (suite *ServerTestSuite) TestGetOrders() {
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

	resp, bodyResp := MakeAuthRequest(suite, ts, "GET", "/api/user/orders", nil, true, suite.token)

	var ordersResponse []handlers.OrderReponse

	err := json.Unmarshal([]byte(bodyResp), &ordersResponse)
	suite.Require().NoError(err)

	defer resp.Body.Close()
	suite.Require().Equal(200, resp.StatusCode)
}
