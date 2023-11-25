// Code generated by MockGen. DO NOT EDIT.
// Source: ./cron_job.go
//
// Generated by this command:
//
//	mockgen -source=./cron_job.go -package=repomocks -destination=mocks/cron_job.mock.go CronJobRepository
//
// Package repomocks is a generated GoMock package.
package repomocks

import (
	context "context"
	reflect "reflect"
	time "time"

	domain "github.com/gevinzone/basic-go/week9/webook/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockCronJobRepository is a mock of CronJobRepository interface.
type MockCronJobRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCronJobRepositoryMockRecorder
}

// MockCronJobRepositoryMockRecorder is the mock recorder for MockCronJobRepository.
type MockCronJobRepositoryMockRecorder struct {
	mock *MockCronJobRepository
}

// NewMockCronJobRepository creates a new mock instance.
func NewMockCronJobRepository(ctrl *gomock.Controller) *MockCronJobRepository {
	mock := &MockCronJobRepository{ctrl: ctrl}
	mock.recorder = &MockCronJobRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCronJobRepository) EXPECT() *MockCronJobRepositoryMockRecorder {
	return m.recorder
}

// Preempt mocks base method.
func (m *MockCronJobRepository) Preempt(ctx context.Context) (domain.CronJob, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Preempt", ctx)
	ret0, _ := ret[0].(domain.CronJob)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Preempt indicates an expected call of Preempt.
func (mr *MockCronJobRepositoryMockRecorder) Preempt(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Preempt", reflect.TypeOf((*MockCronJobRepository)(nil).Preempt), ctx)
}

// Release mocks base method.
func (m *MockCronJobRepository) Release(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Release", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Release indicates an expected call of Release.
func (mr *MockCronJobRepositoryMockRecorder) Release(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Release", reflect.TypeOf((*MockCronJobRepository)(nil).Release), ctx, id)
}

// UpdateNextTime mocks base method.
func (m *MockCronJobRepository) UpdateNextTime(ctx context.Context, id int64, t time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNextTime", ctx, id, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNextTime indicates an expected call of UpdateNextTime.
func (mr *MockCronJobRepositoryMockRecorder) UpdateNextTime(ctx, id, t any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNextTime", reflect.TypeOf((*MockCronJobRepository)(nil).UpdateNextTime), ctx, id, t)
}

// UpdateUtime mocks base method.
func (m *MockCronJobRepository) UpdateUtime(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUtime", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUtime indicates an expected call of UpdateUtime.
func (mr *MockCronJobRepositoryMockRecorder) UpdateUtime(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUtime", reflect.TypeOf((*MockCronJobRepository)(nil).UpdateUtime), ctx, id)
}
