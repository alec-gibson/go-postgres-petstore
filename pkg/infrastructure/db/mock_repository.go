// Code generated by MockGen. DO NOT EDIT.
// Source: alecgibson.ca/go-postgres-petstore/pkg/infrastructure/db (interfaces: Querier)

// Package db is a generated GoMock package.
package db

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockQuerier is a mock of Querier interface.
type MockQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockQuerierMockRecorder
}

// MockQuerierMockRecorder is the mock recorder for MockQuerier.
type MockQuerierMockRecorder struct {
	mock *MockQuerier
}

// NewMockQuerier creates a new mock instance.
func NewMockQuerier(ctrl *gomock.Controller) *MockQuerier {
	mock := &MockQuerier{ctrl: ctrl}
	mock.recorder = &MockQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuerier) EXPECT() *MockQuerierMockRecorder {
	return m.recorder
}

// CreatePet mocks base method.
func (m *MockQuerier) CreatePet(arg0 context.Context, arg1 CreatePetParams) (PetstorePet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePet", arg0, arg1)
	ret0, _ := ret[0].(PetstorePet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePet indicates an expected call of CreatePet.
func (mr *MockQuerierMockRecorder) CreatePet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePet", reflect.TypeOf((*MockQuerier)(nil).CreatePet), arg0, arg1)
}

// DeletePet mocks base method.
func (m *MockQuerier) DeletePet(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePet", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePet indicates an expected call of DeletePet.
func (mr *MockQuerierMockRecorder) DeletePet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePet", reflect.TypeOf((*MockQuerier)(nil).DeletePet), arg0, arg1)
}

// FindPetByID mocks base method.
func (m *MockQuerier) FindPetByID(arg0 context.Context, arg1 int64) (PetstorePet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPetByID", arg0, arg1)
	ret0, _ := ret[0].(PetstorePet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindPetByID indicates an expected call of FindPetByID.
func (mr *MockQuerierMockRecorder) FindPetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPetByID", reflect.TypeOf((*MockQuerier)(nil).FindPetByID), arg0, arg1)
}

// ListPets mocks base method.
func (m *MockQuerier) ListPets(arg0 context.Context, arg1 []string) ([]PetstorePet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPets", arg0, arg1)
	ret0, _ := ret[0].([]PetstorePet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPets indicates an expected call of ListPets.
func (mr *MockQuerierMockRecorder) ListPets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPets", reflect.TypeOf((*MockQuerier)(nil).ListPets), arg0, arg1)
}

// ListPetsWithLimit mocks base method.
func (m *MockQuerier) ListPetsWithLimit(arg0 context.Context, arg1 ListPetsWithLimitParams) ([]PetstorePet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPetsWithLimit", arg0, arg1)
	ret0, _ := ret[0].([]PetstorePet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPetsWithLimit indicates an expected call of ListPetsWithLimit.
func (mr *MockQuerierMockRecorder) ListPetsWithLimit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPetsWithLimit", reflect.TypeOf((*MockQuerier)(nil).ListPetsWithLimit), arg0, arg1)
}