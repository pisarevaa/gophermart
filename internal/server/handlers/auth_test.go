package handlers_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/pisarevaa/gophermart/internal/server"
	"github.com/pisarevaa/gophermart/internal/server/configs"
	mock "github.com/pisarevaa/gophermart/internal/server/mocks"
	"github.com/pisarevaa/gophermart/internal/server/storage"
	"github.com/pisarevaa/gophermart/internal/server/utils"
)

type ServerTestSuite struct {
	suite.Suite
	cfg    configs.Config
	logger *zap.SugaredLogger
	client *resty.Client
	token  string
}

const login = "test"

func (suite *ServerTestSuite) SetupSuite() {
	suite.cfg = configs.NewConfig()
	suite.logger = server.NewLogger()
	suite.client = resty.New()
	token, err := utils.GenerateJWTString(suite.cfg.TokenExpSec, suite.cfg.SecretKey, login)
	suite.Require().NoError(err)
	suite.token = token
}

func TestAgentSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (suite *ServerTestSuite) TestRegisterUserMockDB() {
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

	user := storage.RegisterUser{
		Login:    "test",
		Password: "123",
	}
	resp, err := suite.client.R().
		SetBody(user).
		SetHeader("Content-Type", "application/json").
		Post(ts.URL + "/api/user/register")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
}

func (suite *ServerTestSuite) TestLoginMockDB() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	passwordHash, err := utils.GetPasswordHash("123", suite.cfg.SecretKey)
	suite.Require().NoError(err)
	dbUser := storage.User{
		Login:    "test",
		Password: passwordHash,
		Balance:  500,
	}

	m.EXPECT().
		GetUser(gomock.Any(), gomock.Any()).
		Return(dbUser, nil)

	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, m))
	defer ts.Close()

	user := storage.RegisterUser{
		Login:    "test",
		Password: "123",
	}

	resp, err := suite.client.R().
		SetBody(user).
		SetHeader("Content-Type", "application/json").
		Post(ts.URL + "/api/user/login")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
}

func (suite *ServerTestSuite) TestRegisterUserAndLoginInMemory() {
	m := storage.NewMemory()

	ts := httptest.NewServer(server.NewRouter(suite.cfg, suite.logger, m))
	defer ts.Close()

	user := storage.RegisterUser{
		Login:    "test",
		Password: "123",
	}
	resp, err := suite.client.R().
		SetBody(user).
		SetHeader("Content-Type", "application/json").
		Post(ts.URL + "/api/user/register")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())

	resp, err = suite.client.R().
		SetBody(user).
		SetHeader("Content-Type", "application/json").
		Post(ts.URL + "/api/user/login")
	suite.Require().NoError(err)
	suite.Require().Equal(200, resp.StatusCode())
}
