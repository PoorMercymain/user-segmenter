// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/PoorMercymain/user-segmenter/internal/domain (interfaces: SegmentRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSegmentRepository is a mock of SegmentRepository interface.
type MockSegmentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSegmentRepositoryMockRecorder
}

// MockSegmentRepositoryMockRecorder is the mock recorder for MockSegmentRepository.
type MockSegmentRepositoryMockRecorder struct {
	mock *MockSegmentRepository
}

// NewMockSegmentRepository creates a new mock instance.
func NewMockSegmentRepository(ctrl *gomock.Controller) *MockSegmentRepository {
	mock := &MockSegmentRepository{ctrl: ctrl}
	mock.recorder = &MockSegmentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSegmentRepository) EXPECT() *MockSegmentRepositoryMockRecorder {
	return m.recorder
}

// CreateSegment mocks base method.
func (m *MockSegmentRepository) CreateSegment(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSegment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSegment indicates an expected call of CreateSegment.
func (mr *MockSegmentRepositoryMockRecorder) CreateSegment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSegment", reflect.TypeOf((*MockSegmentRepository)(nil).CreateSegment), arg0, arg1)
}

// DeleteSegment mocks base method.
func (m *MockSegmentRepository) DeleteSegment(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSegment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSegment indicates an expected call of DeleteSegment.
func (mr *MockSegmentRepositoryMockRecorder) DeleteSegment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSegment", reflect.TypeOf((*MockSegmentRepository)(nil).DeleteSegment), arg0, arg1)
}

// ReadUserSegments mocks base method.
func (m *MockSegmentRepository) ReadUserSegments(arg0 context.Context, arg1 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadUserSegments", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadUserSegments indicates an expected call of ReadUserSegments.
func (mr *MockSegmentRepositoryMockRecorder) ReadUserSegments(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadUserSegments", reflect.TypeOf((*MockSegmentRepository)(nil).ReadUserSegments), arg0, arg1)
}

// UpdateUserSegments mocks base method.
func (m *MockSegmentRepository) UpdateUserSegments(arg0 context.Context, arg1 string, arg2, arg3 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserSegments", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserSegments indicates an expected call of UpdateUserSegments.
func (mr *MockSegmentRepositoryMockRecorder) UpdateUserSegments(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserSegments", reflect.TypeOf((*MockSegmentRepository)(nil).UpdateUserSegments), arg0, arg1, arg2, arg3)
}