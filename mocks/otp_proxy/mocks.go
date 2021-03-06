// Code generated by MockGen. DO NOT EDIT.
// Source: ./otp_proxy.go

// Package mockOtpProxy is a generated GoMock package.
package mockOtpProxy

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockOtpProxy is a mock of OtpProxy interface.
type MockOtpProxy struct {
	ctrl     *gomock.Controller
	recorder *MockOtpProxyMockRecorder
}

// MockOtpProxyMockRecorder is the mock recorder for MockOtpProxy.
type MockOtpProxyMockRecorder struct {
	mock *MockOtpProxy
}

// NewMockOtpProxy creates a new mock instance.
func NewMockOtpProxy(ctrl *gomock.Controller) *MockOtpProxy {
	mock := &MockOtpProxy{ctrl: ctrl}
	mock.recorder = &MockOtpProxyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOtpProxy) EXPECT() *MockOtpProxyMockRecorder {
	return m.recorder
}

// GetOtp mocks base method.
func (m *MockOtpProxy) GetOtp(accountID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOtp", accountID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOtp indicates an expected call of GetOtp.
func (mr *MockOtpProxyMockRecorder) GetOtp(accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOtp", reflect.TypeOf((*MockOtpProxy)(nil).GetOtp), accountID)
}
