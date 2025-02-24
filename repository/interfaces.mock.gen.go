// Code generated by MockGen. DO NOT EDIT.
// Source: repository/interfaces.go
//
// Generated by this command:
//
//	mockgen -source=repository/interfaces.go -destination=repository/interfaces.mock.gen.go -package=repository
//

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockRepositoryInterface is a mock of RepositoryInterface interface.
type MockRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryInterfaceMockRecorder
	isgomock struct{}
}

// MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.
type MockRepositoryInterfaceMockRecorder struct {
	mock *MockRepositoryInterface
}

// NewMockRepositoryInterface creates a new mock instance.
func NewMockRepositoryInterface(ctrl *gomock.Controller) *MockRepositoryInterface {
	mock := &MockRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryInterface) EXPECT() *MockRepositoryInterfaceMockRecorder {
	return m.recorder
}

// CreateEstate mocks base method.
func (m *MockRepositoryInterface) CreateEstate(ctx context.Context, input CreateEstateInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEstate", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateEstate indicates an expected call of CreateEstate.
func (mr *MockRepositoryInterfaceMockRecorder) CreateEstate(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEstate", reflect.TypeOf((*MockRepositoryInterface)(nil).CreateEstate), ctx, input)
}

// CreateTree mocks base method.
func (m *MockRepositoryInterface) CreateTree(ctx context.Context, input CreateTreeInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTree", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTree indicates an expected call of CreateTree.
func (mr *MockRepositoryInterfaceMockRecorder) CreateTree(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTree", reflect.TypeOf((*MockRepositoryInterface)(nil).CreateTree), ctx, input)
}

// GetCalculatedEstateStats mocks base method.
func (m *MockRepositoryInterface) GetCalculatedEstateStats(ctx context.Context, estateId uuid.UUID) (*EstateStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCalculatedEstateStats", ctx, estateId)
	ret0, _ := ret[0].(*EstateStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCalculatedEstateStats indicates an expected call of GetCalculatedEstateStats.
func (mr *MockRepositoryInterfaceMockRecorder) GetCalculatedEstateStats(ctx, estateId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCalculatedEstateStats", reflect.TypeOf((*MockRepositoryInterface)(nil).GetCalculatedEstateStats), ctx, estateId)
}

// GetEstateWithAllDetails mocks base method.
func (m *MockRepositoryInterface) GetEstateWithAllDetails(ctx context.Context, id uuid.UUID, exludeRelations ...Relation) (*Estate, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, id}
	for _, a := range exludeRelations {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetEstateWithAllDetails", varargs...)
	ret0, _ := ret[0].(*Estate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEstateWithAllDetails indicates an expected call of GetEstateWithAllDetails.
func (mr *MockRepositoryInterfaceMockRecorder) GetEstateWithAllDetails(ctx, id any, exludeRelations ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, id}, exludeRelations...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEstateWithAllDetails", reflect.TypeOf((*MockRepositoryInterface)(nil).GetEstateWithAllDetails), varargs...)
}

// UpsertEstateStats mocks base method.
func (m *MockRepositoryInterface) UpsertEstateStats(ctx context.Context, estateID uuid.UUID, stats *EstateStats) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertEstateStats", ctx, estateID, stats)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertEstateStats indicates an expected call of UpsertEstateStats.
func (mr *MockRepositoryInterfaceMockRecorder) UpsertEstateStats(ctx, estateID, stats any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertEstateStats", reflect.TypeOf((*MockRepositoryInterface)(nil).UpsertEstateStats), ctx, estateID, stats)
}
