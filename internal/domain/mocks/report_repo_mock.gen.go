// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/PoorMercymain/user-segmenter/internal/domain (interfaces: ReportRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	io "io"
	reflect "reflect"
	time "time"

	domain "github.com/PoorMercymain/user-segmenter/internal/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockReportRepository is a mock of ReportRepository interface.
type MockReportRepository struct {
	ctrl     *gomock.Controller
	recorder *MockReportRepositoryMockRecorder
}

// MockReportRepositoryMockRecorder is the mock recorder for MockReportRepository.
type MockReportRepositoryMockRecorder struct {
	mock *MockReportRepository
}

// NewMockReportRepository creates a new mock instance.
func NewMockReportRepository(ctrl *gomock.Controller) *MockReportRepository {
	mock := &MockReportRepository{ctrl: ctrl}
	mock.recorder = &MockReportRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReportRepository) EXPECT() *MockReportRepositoryMockRecorder {
	return m.recorder
}

// CreateCSV mocks base method.
func (m *MockReportRepository) CreateCSV(arg0 context.Context, arg1 string, arg2, arg3 time.Time) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCSV", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCSV indicates an expected call of CreateCSV.
func (mr *MockReportRepositoryMockRecorder) CreateCSV(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCSV", reflect.TypeOf((*MockReportRepository)(nil).CreateCSV), arg0, arg1, arg2, arg3)
}

// ReadUserSegmentsHistory mocks base method.
func (m *MockReportRepository) ReadUserSegmentsHistory(arg0 context.Context, arg1 string, arg2, arg3 time.Time, arg4, arg5 int) ([]domain.HistoryElem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadUserSegmentsHistory", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].([]domain.HistoryElem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadUserSegmentsHistory indicates an expected call of ReadUserSegmentsHistory.
func (mr *MockReportRepositoryMockRecorder) ReadUserSegmentsHistory(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadUserSegmentsHistory", reflect.TypeOf((*MockReportRepository)(nil).ReadUserSegmentsHistory), arg0, arg1, arg2, arg3, arg4, arg5)
}

// SendCSVReportFile mocks base method.
func (m *MockReportRepository) SendCSVReportFile(arg0 string, arg1 io.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCSVReportFile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCSVReportFile indicates an expected call of SendCSVReportFile.
func (mr *MockReportRepositoryMockRecorder) SendCSVReportFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCSVReportFile", reflect.TypeOf((*MockReportRepository)(nil).SendCSVReportFile), arg0, arg1)
}
