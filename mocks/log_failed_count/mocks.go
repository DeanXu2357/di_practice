// Code generated by MockGen. DO NOT EDIT.
// Source: ./log_failed_count.go

// Package mockLogFailedCount is a generated GoMock package.
package mockLogFailedCount

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLogFailedCount is a mock of LogFailedCount interface.
type MockLogFailedCount struct {
	ctrl     *gomock.Controller
	recorder *MockLogFailedCountMockRecorder
}

// MockLogFailedCountMockRecorder is the mock recorder for MockLogFailedCount.
type MockLogFailedCountMockRecorder struct {
	mock *MockLogFailedCount
}

// NewMockLogFailedCount creates a new mock instance.
func NewMockLogFailedCount(ctrl *gomock.Controller) *MockLogFailedCount {
	mock := &MockLogFailedCount{ctrl: ctrl}
	mock.recorder = &MockLogFailedCountMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogFailedCount) EXPECT() *MockLogFailedCountMockRecorder {
	return m.recorder
}

// LogFailedCount mocks base method.
func (m *MockLogFailedCount) LogFailedCount(accountID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogFailedCount", accountID)
	ret0, _ := ret[0].(error)
	return ret0
}

// LogFailedCount indicates an expected call of LogFailedCount.
func (mr *MockLogFailedCountMockRecorder) LogFailedCount(accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogFailedCount", reflect.TypeOf((*MockLogFailedCount)(nil).LogFailedCount), accountID)
}
