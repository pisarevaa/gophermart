package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/golang/mock/gomock"

	"github.com/pisarevaa/gophermart/internal/server"
	"github.com/pisarevaa/gophermart/internal/server/configs"
	mock "github.com/pisarevaa/gophermart/internal/server/mocks"
	"github.com/pisarevaa/gophermart/internal/server/storage"
)

type ServerTestSuite struct {
	suite.Suite
	cfg    configs.Config
	logger *zap.SugaredLogger
}

func (suite *ServerTestSuite) SetupSuite() {
	suite.cfg = configs.NewConfig()
	suite.logger = server.NewLogger()
}

func TestAgentSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func MakeRequest(
	suite *ServerTestSuite,
	ts *httptest.Server,
	method string,
	url string,
	body []byte,
) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+url, bytes.NewBuffer(body))
	suite.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "")
	resp, err := ts.Client().Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	suite.Require().NoError(err)
	return resp, string(respBody)
}

func (suite *ServerTestSuite) TestServerUpdateAndGetMetricsJSONBatch() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	m.EXPECT().
	GetUser(gomock.Any(), gomock.Any()).
		Return(storage.User{}, errors.New("user exists"))

	m.EXPECT().
	StoreUser(gomock.Any(), gomock.Any(), gomock.Any()).
	Return(nil)

	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, m))
	defer ts.Close()

	user := storage.User{
		Login:    "test",
		Password: "123",
	}
	userJson, _ := json.Marshal(user)

	resp, _ := MakeRequest(suite, ts, "POST", "/api/user/register", userJson)

	defer resp.Body.Close()
	suite.Require().Equal(200, resp.StatusCode)
}
