// Code generated by MockGen. DO NOT EDIT.
// Source: internal/pkg/hasher/hasher.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHasher is a mock of Hasher interface.
type MockHasher struct {
	ctrl     *gomock.Controller
	recorder *MockHasherMockRecorder
}

// MockHasherMockRecorder is the mock recorder for MockHasher.
type MockHasherMockRecorder struct {
	mock *MockHasher
}

// NewMockHasher creates a new mock instance.
func NewMockHasher(ctrl *gomock.Controller) *MockHasher {
	mock := &MockHasher{ctrl: ctrl}
	mock.recorder = &MockHasherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHasher) EXPECT() *MockHasherMockRecorder {
	return m.recorder
}

// CompareHashAndPassword mocks base method.
func (m *MockHasher) CompareHashAndPassword(hashedPassword, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompareHashAndPassword", hashedPassword, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompareHashAndPassword indicates an expected call of CompareHashAndPassword.
func (mr *MockHasherMockRecorder) CompareHashAndPassword(hashedPassword, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompareHashAndPassword", reflect.TypeOf((*MockHasher)(nil).CompareHashAndPassword), hashedPassword, password)
}

// GetHashedPassword mocks base method.
func (m *MockHasher) GetHashedPassword(password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHashedPassword", password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHashedPassword indicates an expected call of GetHashedPassword.
func (mr *MockHasherMockRecorder) GetHashedPassword(password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHashedPassword", reflect.TypeOf((*MockHasher)(nil).GetHashedPassword), password)
}
