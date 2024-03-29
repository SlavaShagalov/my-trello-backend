// Code generated by MockGen. DO NOT EDIT.
// Source: internal/boards/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	boards "github.com/SlavaShagalov/my-trello-backend/internal/boards"
	models "github.com/SlavaShagalov/my-trello-backend/internal/models"
	gomock "github.com/golang/mock/gomock"
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

// Create mocks base method.
func (m *MockRepository) Create(ctx context.Context, params *boards.CreateParams) (models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, params)
	ret0, _ := ret[0].(models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), ctx, params)
}

// Delete mocks base method.
func (m *MockRepository) Delete(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), ctx, id)
}

// FullUpdate mocks base method.
func (m *MockRepository) FullUpdate(ctx context.Context, params *boards.FullUpdateParams) (models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FullUpdate", ctx, params)
	ret0, _ := ret[0].(models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FullUpdate indicates an expected call of FullUpdate.
func (mr *MockRepositoryMockRecorder) FullUpdate(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FullUpdate", reflect.TypeOf((*MockRepository)(nil).FullUpdate), ctx, params)
}

// Get mocks base method.
func (m *MockRepository) Get(ctx context.Context, id int) (models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepositoryMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepository)(nil).Get), ctx, id)
}

// List mocks base method.
func (m *MockRepository) List(ctx context.Context, workspaceID int) ([]models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, workspaceID)
	ret0, _ := ret[0].([]models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockRepositoryMockRecorder) List(ctx, workspaceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRepository)(nil).List), ctx, workspaceID)
}

// ListByTitle mocks base method.
func (m *MockRepository) ListByTitle(ctx context.Context, title string, userID int) ([]models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByTitle", ctx, title, userID)
	ret0, _ := ret[0].([]models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByTitle indicates an expected call of ListByTitle.
func (mr *MockRepositoryMockRecorder) ListByTitle(ctx, title, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByTitle", reflect.TypeOf((*MockRepository)(nil).ListByTitle), ctx, title, userID)
}

// PartialUpdate mocks base method.
func (m *MockRepository) PartialUpdate(ctx context.Context, params *boards.PartialUpdateParams) (models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PartialUpdate", ctx, params)
	ret0, _ := ret[0].(models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PartialUpdate indicates an expected call of PartialUpdate.
func (mr *MockRepositoryMockRecorder) PartialUpdate(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PartialUpdate", reflect.TypeOf((*MockRepository)(nil).PartialUpdate), ctx, params)
}

// UpdateBackground mocks base method.
func (m *MockRepository) UpdateBackground(ctx context.Context, id int, background string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBackground", ctx, id, background)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBackground indicates an expected call of UpdateBackground.
func (mr *MockRepositoryMockRecorder) UpdateBackground(ctx, id, background interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBackground", reflect.TypeOf((*MockRepository)(nil).UpdateBackground), ctx, id, background)
}
