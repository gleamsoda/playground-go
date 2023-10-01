// Code generated by MockGen. DO NOT EDIT.
// Source: playground/app (interfaces: Repository)

// Package mock_app is a generated GoMock package.
package mock_app

import (
	context "context"
	app "playground/app"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockRepository) CreateAccount(arg0 context.Context, arg1 *app.Account) (*app.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0, arg1)
	ret0, _ := ret[0].(*app.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockRepositoryMockRecorder) CreateAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockRepository)(nil).CreateAccount), arg0, arg1)
}

// CreateSession mocks base method.
func (m *MockRepository) CreateSession(arg0 context.Context, arg1 *app.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockRepositoryMockRecorder) CreateSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockRepository)(nil).CreateSession), arg0, arg1)
}

// CreateTransfer mocks base method.
func (m *MockRepository) CreateTransfer(arg0 context.Context, arg1 *app.Transfer) (*app.Transfer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransfer", arg0, arg1)
	ret0, _ := ret[0].(*app.Transfer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransfer indicates an expected call of CreateTransfer.
func (mr *MockRepositoryMockRecorder) CreateTransfer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransfer", reflect.TypeOf((*MockRepository)(nil).CreateTransfer), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockRepository) CreateUser(arg0 context.Context, arg1 *app.User) (*app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(*app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepositoryMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepository)(nil).CreateUser), arg0, arg1)
}

// CreateVerifyEmail mocks base method.
func (m *MockRepository) CreateVerifyEmail(arg0 context.Context, arg1 *app.VerifyEmail) (*app.VerifyEmail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVerifyEmail", arg0, arg1)
	ret0, _ := ret[0].(*app.VerifyEmail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateVerifyEmail indicates an expected call of CreateVerifyEmail.
func (mr *MockRepositoryMockRecorder) CreateVerifyEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVerifyEmail", reflect.TypeOf((*MockRepository)(nil).CreateVerifyEmail), arg0, arg1)
}

// DeleteAccount mocks base method.
func (m *MockRepository) DeleteAccount(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAccount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAccount indicates an expected call of DeleteAccount.
func (mr *MockRepositoryMockRecorder) DeleteAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccount", reflect.TypeOf((*MockRepository)(nil).DeleteAccount), arg0, arg1)
}

// GetAccount mocks base method.
func (m *MockRepository) GetAccount(arg0 context.Context, arg1 int64) (*app.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", arg0, arg1)
	ret0, _ := ret[0].(*app.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockRepositoryMockRecorder) GetAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockRepository)(nil).GetAccount), arg0, arg1)
}

// GetSession mocks base method.
func (m *MockRepository) GetSession(arg0 context.Context, arg1 uuid.UUID) (*app.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", arg0, arg1)
	ret0, _ := ret[0].(*app.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSession indicates an expected call of GetSession.
func (mr *MockRepositoryMockRecorder) GetSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockRepository)(nil).GetSession), arg0, arg1)
}

// GetTransfer mocks base method.
func (m *MockRepository) GetTransfer(arg0 context.Context, arg1 int64) (*app.Transfer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransfer", arg0, arg1)
	ret0, _ := ret[0].(*app.Transfer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransfer indicates an expected call of GetTransfer.
func (mr *MockRepositoryMockRecorder) GetTransfer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransfer", reflect.TypeOf((*MockRepository)(nil).GetTransfer), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockRepository) GetUser(arg0 context.Context, arg1 string) (*app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(*app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockRepositoryMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockRepository)(nil).GetUser), arg0, arg1)
}

// ListAccounts mocks base method.
func (m *MockRepository) ListAccounts(arg0 context.Context, arg1 *app.ListAccountsParams) ([]app.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAccounts", arg0, arg1)
	ret0, _ := ret[0].([]app.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAccounts indicates an expected call of ListAccounts.
func (mr *MockRepositoryMockRecorder) ListAccounts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAccounts", reflect.TypeOf((*MockRepository)(nil).ListAccounts), arg0, arg1)
}

// ListTransfers mocks base method.
func (m *MockRepository) ListTransfers(arg0 context.Context, arg1 *app.ListTransfersParams) ([]app.Transfer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTransfers", arg0, arg1)
	ret0, _ := ret[0].([]app.Transfer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTransfers indicates an expected call of ListTransfers.
func (mr *MockRepositoryMockRecorder) ListTransfers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTransfers", reflect.TypeOf((*MockRepository)(nil).ListTransfers), arg0, arg1)
}

// UpdateAccount mocks base method.
func (m *MockRepository) UpdateAccount(arg0 context.Context, arg1 *app.Account) (*app.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccount", arg0, arg1)
	ret0, _ := ret[0].(*app.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAccount indicates an expected call of UpdateAccount.
func (mr *MockRepositoryMockRecorder) UpdateAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccount", reflect.TypeOf((*MockRepository)(nil).UpdateAccount), arg0, arg1)
}

// UpdateUser mocks base method.
func (m *MockRepository) UpdateUser(arg0 context.Context, arg1 *app.User) (*app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1)
	ret0, _ := ret[0].(*app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockRepositoryMockRecorder) UpdateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockRepository)(nil).UpdateUser), arg0, arg1)
}

// UpdateUserEmailVerified mocks base method.
func (m *MockRepository) UpdateUserEmailVerified(arg0 context.Context, arg1 *app.VerifyEmailParams) (*app.User, *app.VerifyEmail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserEmailVerified", arg0, arg1)
	ret0, _ := ret[0].(*app.User)
	ret1, _ := ret[1].(*app.VerifyEmail)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UpdateUserEmailVerified indicates an expected call of UpdateUserEmailVerified.
func (mr *MockRepositoryMockRecorder) UpdateUserEmailVerified(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserEmailVerified", reflect.TypeOf((*MockRepository)(nil).UpdateUserEmailVerified), arg0, arg1)
}

// UpdateVerifyEmail mocks base method.
func (m *MockRepository) UpdateVerifyEmail(arg0 context.Context, arg1 *app.VerifyEmail) (*app.VerifyEmail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVerifyEmail", arg0, arg1)
	ret0, _ := ret[0].(*app.VerifyEmail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateVerifyEmail indicates an expected call of UpdateVerifyEmail.
func (mr *MockRepositoryMockRecorder) UpdateVerifyEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVerifyEmail", reflect.TypeOf((*MockRepository)(nil).UpdateVerifyEmail), arg0, arg1)
}
