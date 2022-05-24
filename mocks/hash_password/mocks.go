// Code generated by MockGen. DO NOT EDIT.
// Source: ./hash_password.go

// Package mockHashPassword is a generated GoMock package.
package mockHashPassword

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHashPassword is a mock of HashPassword interface.
type MockHashPassword struct {
	ctrl     *gomock.Controller
	recorder *MockHashPasswordMockRecorder
}

// MockHashPasswordMockRecorder is the mock recorder for MockHashPassword.
type MockHashPasswordMockRecorder struct {
	mock *MockHashPassword
}

// NewMockHashPassword creates a new mock instance.
func NewMockHashPassword(ctrl *gomock.Controller) *MockHashPassword {
	mock := &MockHashPassword{ctrl: ctrl}
	mock.recorder = &MockHashPasswordMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHashPassword) EXPECT() *MockHashPasswordMockRecorder {
	return m.recorder
}

// HashPassword mocks base method.
func (m *MockHashPassword) HashPassword(pwd string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HashPassword", pwd)
	ret0, _ := ret[0].(string)
	return ret0
}

// HashPassword indicates an expected call of HashPassword.
func (mr *MockHashPasswordMockRecorder) HashPassword(pwd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HashPassword", reflect.TypeOf((*MockHashPassword)(nil).HashPassword), pwd)
}
