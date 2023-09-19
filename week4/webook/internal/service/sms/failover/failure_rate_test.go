// Copyright 2023 igevin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package failover

import (
	"context"
	"errors"
	svcmocks "github.com/gevinzone/basic-go/week4/webook/internal/service/mocks"
	"github.com/gevinzone/basic-go/week4/webook/internal/service/sms"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestFailureRateFailOverService_Send_Normal(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) sms.Service
		ctx     context.Context
		tpl     string
		args    []string
		numbers []string
		wantErr error
	}{
		{
			name: "正常，从第0个开始",
			mock: func(ctrl *gomock.Controller) sms.Service {
				svc := svcmocks.NewMockService(ctrl)
				svc.EXPECT().
					Send(context.Background(), "tmp", []string{"a", "b", "c"}, []string{"1", "2", "3"}).
					Return(nil)
				return NewFailureRateFailOverService([]sms.Service{svc}, 0.5)
			},
			ctx:     context.Background(),
			tpl:     "tmp",
			args:    []string{"a", "b", "c"},
			numbers: []string{"1", "2", "3"},
			wantErr: nil,
		},
		{
			name: "报错，从第0个开始",
			mock: func(ctrl *gomock.Controller) sms.Service {
				svc := svcmocks.NewMockService(ctrl)
				svc.EXPECT().
					Send(context.Background(), "tmp", []string{"a", "b", "c"}, []string{"1", "2", "3"}).
					Return(errors.New("error"))
				return NewFailureRateFailOverService([]sms.Service{svc}, 0.5)
			},
			ctx:     context.Background(),
			tpl:     "tmp",
			args:    []string{"a", "b", "c"},
			numbers: []string{"1", "2", "3"},
			wantErr: errors.New("error"),
		},
		{
			name: "超时，从第0个开始",
			mock: func(ctrl *gomock.Controller) sms.Service {
				svc := svcmocks.NewMockService(ctrl)
				svc.EXPECT().
					Send(context.Background(), "tmp", []string{"a", "b", "c"}, []string{"1", "2", "3"}).
					Return(context.DeadlineExceeded)
				return NewFailureRateFailOverService([]sms.Service{svc}, 0.5)
			},
			ctx:     context.Background(),
			tpl:     "tmp",
			args:    []string{"a", "b", "c"},
			numbers: []string{"1", "2", "3"},
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "正常，从第1个开始",
			mock: func(ctrl *gomock.Controller) sms.Service {
				svc := svcmocks.NewMockService(ctrl)
				svc.EXPECT().
					Send(context.Background(), "tmp", []string{"a", "b", "c"}, []string{"1", "2", "3"}).
					Return(nil)
				res := NewFailureRateFailOverService([]sms.Service{nil, svc}, 0.5)
				instance := res.(*FailureRateFailOverService)
				instance.idx = 1
				return instance
			},
			ctx:     context.Background(),
			tpl:     "tmp",
			args:    []string{"a", "b", "c"},
			numbers: []string{"1", "2", "3"},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := tc.mock(ctrl)
			err := svc.Send(tc.ctx, tc.tpl, tc.args, tc.numbers...)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestFailureRateFailOverService_Send_FailOver(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) *FailureRateFailOverService
		ctx        context.Context
		tpl        string
		args       []string
		numbers    []string
		successCnt int64
		failureCnt int64
		wantErr    error
	}{
		{
			name: "无报错，不需要fail over",
			mock: func(ctrl *gomock.Controller) *FailureRateFailOverService {
				svc := svcmocks.NewMockService(ctrl)
				svc.EXPECT().
					Send(context.Background(), "tmp", []string{"a", "b", "c"}, []string{"1", "2", "3"}).
					Return(nil)
				s := NewFailureRateFailOverService([]sms.Service{svc}, 0.5)
				instance := s.(*FailureRateFailOverService)
				instance.successCnt = 5
				return instance
			},
			ctx:        context.Background(),
			tpl:        "tmp",
			args:       []string{"a", "b", "c"},
			numbers:    []string{"1", "2", "3"},
			successCnt: int64(6),
			wantErr:    nil,
		},
		{
			name: "前面有报错，不需要fail over，当前无报错",
			mock: func(ctrl *gomock.Controller) *FailureRateFailOverService {
				svc := svcmocks.NewMockService(ctrl)
				svc.EXPECT().
					Send(context.Background(), "tmp", []string{"a", "b", "c"}, []string{"1", "2", "3"}).
					Return(nil)
				s := NewFailureRateFailOverService([]sms.Service{svc}, 0.5)
				instance := s.(*FailureRateFailOverService)
				instance.successCnt = 5
				instance.failureCnt = 1
				return instance
			},
			ctx:        context.Background(),
			tpl:        "tmp",
			args:       []string{"a", "b", "c"},
			numbers:    []string{"1", "2", "3"},
			successCnt: 6,
			failureCnt: 1,
			wantErr:    nil,
		},
		{
			name: "前面有报错，不需要fail over，此次调用超时报错",
			mock: func(ctrl *gomock.Controller) *FailureRateFailOverService {
				svc := svcmocks.NewMockService(ctrl)
				svc.EXPECT().
					Send(context.Background(), "tmp", []string{"a", "b", "c"}, []string{"1", "2", "3"}).
					Return(context.DeadlineExceeded)
				s := NewFailureRateFailOverService([]sms.Service{svc}, 0.5)
				instance := s.(*FailureRateFailOverService)
				instance.successCnt = 5
				instance.failureCnt = 1
				return instance
			},
			ctx:        context.Background(),
			tpl:        "tmp",
			args:       []string{"a", "b", "c"},
			numbers:    []string{"1", "2", "3"},
			successCnt: 5,
			failureCnt: 2,
			wantErr:    context.DeadlineExceeded,
		},
		{
			name: "前面有报错，不需要fail over，此次调用有未知错误",
			mock: func(ctrl *gomock.Controller) *FailureRateFailOverService {
				svc := svcmocks.NewMockService(ctrl)
				svc.EXPECT().
					Send(context.Background(), "tmp", []string{"a", "b", "c"}, []string{"1", "2", "3"}).
					Return(errors.New("unknown"))
				s := NewFailureRateFailOverService([]sms.Service{svc}, 0.5)
				instance := s.(*FailureRateFailOverService)
				instance.successCnt = 5
				instance.failureCnt = 1
				return instance
			},
			ctx:        context.Background(),
			tpl:        "tmp",
			args:       []string{"a", "b", "c"},
			numbers:    []string{"1", "2", "3"},
			successCnt: 5,
			failureCnt: 1,
			wantErr:    errors.New("unknown"),
		},
		{
			name: "fail over，此次调用成功",
			mock: func(ctrl *gomock.Controller) *FailureRateFailOverService {
				svc := svcmocks.NewMockService(ctrl)
				svc.EXPECT().
					Send(context.Background(), "tmp", []string{"a", "b", "c"}, []string{"1", "2", "3"}).
					Return(nil)
				s := NewFailureRateFailOverService([]sms.Service{nil, svc}, 0.5)
				instance := s.(*FailureRateFailOverService)
				instance.successCnt = 5
				instance.failureCnt = 5
				return instance
			},
			ctx:        context.Background(),
			tpl:        "tmp",
			args:       []string{"a", "b", "c"},
			numbers:    []string{"1", "2", "3"},
			successCnt: 1,
			failureCnt: 0,
			wantErr:    nil,
		},
		{
			name: "fail over，此次调用也超时",
			mock: func(ctrl *gomock.Controller) *FailureRateFailOverService {
				svc := svcmocks.NewMockService(ctrl)
				svc.EXPECT().
					Send(context.Background(), "tmp", []string{"a", "b", "c"}, []string{"1", "2", "3"}).
					Return(context.DeadlineExceeded)
				s := NewFailureRateFailOverService([]sms.Service{nil, svc}, 0.5)
				instance := s.(*FailureRateFailOverService)
				instance.successCnt = 5
				instance.failureCnt = 5
				return instance
			},
			ctx:        context.Background(),
			tpl:        "tmp",
			args:       []string{"a", "b", "c"},
			numbers:    []string{"1", "2", "3"},
			successCnt: 0,
			failureCnt: 1,
			wantErr:    context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := tc.mock(ctrl)
			err := svc.Send(tc.ctx, tc.tpl, tc.args, tc.numbers...)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.successCnt, svc.successCnt)
			assert.Equal(t, tc.failureCnt, svc.failureCnt)
		})
	}
}
