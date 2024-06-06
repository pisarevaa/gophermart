// Code generated by MockGen. DO NOT EDIT.
// Source: internal/server/storage/types.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storage "github.com/pisarevaa/gophermart/internal/server/storage"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AccrualOrderBalance mocks base method.
func (m *MockStorage) AccrualOrderBalance(ctx context.Context, number string, withdraw int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccrualOrderBalance", ctx, number, withdraw)
	ret0, _ := ret[0].(error)
	return ret0
}

// AccrualOrderBalance indicates an expected call of AccrualOrderBalance.
func (mr *MockStorageMockRecorder) AccrualOrderBalance(ctx, number, withdraw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccrualOrderBalance", reflect.TypeOf((*MockStorage)(nil).AccrualOrderBalance), ctx, number, withdraw)
}

// CloseConnection mocks base method.
func (m *MockStorage) CloseConnection() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CloseConnection")
}

// CloseConnection indicates an expected call of CloseConnection.
func (mr *MockStorageMockRecorder) CloseConnection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseConnection", reflect.TypeOf((*MockStorage)(nil).CloseConnection))
}

// GetOrder mocks base method.
func (m *MockStorage) GetOrder(ctx context.Context, number string) (storage.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrder", ctx, number)
	ret0, _ := ret[0].(storage.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrder indicates an expected call of GetOrder.
func (mr *MockStorageMockRecorder) GetOrder(ctx, number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrder", reflect.TypeOf((*MockStorage)(nil).GetOrder), ctx, number)
}

// GetOrders mocks base method.
func (m *MockStorage) GetOrders(ctx context.Context, login string, onlyAccrual bool) ([]storage.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrders", ctx, login, onlyAccrual)
	ret0, _ := ret[0].([]storage.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrders indicates an expected call of GetOrders.
func (mr *MockStorageMockRecorder) GetOrders(ctx, login, onlyAccrual interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrders", reflect.TypeOf((*MockStorage)(nil).GetOrders), ctx, login, onlyAccrual)
}

// GetUser mocks base method.
func (m *MockStorage) GetUser(ctx context.Context, login string) (storage.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, login)
	ret0, _ := ret[0].(storage.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStorageMockRecorder) GetUser(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStorage)(nil).GetUser), ctx, login)
}

// GetUserWithdrawals mocks base method.
func (m *MockStorage) GetUserWithdrawals(ctx context.Context, login string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserWithdrawals", ctx, login)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserWithdrawals indicates an expected call of GetUserWithdrawals.
func (mr *MockStorageMockRecorder) GetUserWithdrawals(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserWithdrawals", reflect.TypeOf((*MockStorage)(nil).GetUserWithdrawals), ctx, login)
}

// StoreOrder mocks base method.
func (m *MockStorage) StoreOrder(ctx context.Context, number, login string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreOrder", ctx, number, login)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreOrder indicates an expected call of StoreOrder.
func (mr *MockStorageMockRecorder) StoreOrder(ctx, number, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreOrder", reflect.TypeOf((*MockStorage)(nil).StoreOrder), ctx, number, login)
}

// StoreUser mocks base method.
func (m *MockStorage) StoreUser(ctx context.Context, login, passwordHash string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreUser", ctx, login, passwordHash)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreUser indicates an expected call of StoreUser.
func (mr *MockStorageMockRecorder) StoreUser(ctx, login, passwordHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreUser", reflect.TypeOf((*MockStorage)(nil).StoreUser), ctx, login, passwordHash)
}

// WithdrawUserBalance mocks base method.
func (m *MockStorage) WithdrawUserBalance(ctx context.Context, login string, withdraw int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithdrawUserBalance", ctx, login, withdraw)
	ret0, _ := ret[0].(error)
	return ret0
}

// WithdrawUserBalance indicates an expected call of WithdrawUserBalance.
func (mr *MockStorageMockRecorder) WithdrawUserBalance(ctx, login, withdraw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithdrawUserBalance", reflect.TypeOf((*MockStorage)(nil).WithdrawUserBalance), ctx, login, withdraw)
}
