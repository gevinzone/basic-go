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

package async

import (
	"context"
	"errors"
	retry "github.com/ecodeclub/ekit/retry"
	"github.com/gevinzone/basic-go/week4/webook/internal/domain"
	"github.com/gevinzone/basic-go/week4/webook/internal/repository"
	repositorymock "github.com/gevinzone/basic-go/week4/webook/internal/repository/mock"
	svcmocks "github.com/gevinzone/basic-go/week4/webook/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestSimpleWorkShop_Consume(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) *SimpleWorkShop
		ctx     context.Context
		wantErr error
	}{
		{
			name: "正常发送",
			mock: func(ctrl *gomock.Controller) *SimpleWorkShop {
				repo := repositorymock.NewMockSmsRepository(ctrl)
				repo.EXPECT().GetFirst(gomock.Any()).Return(createSms(), nil)
				repo.EXPECT().UpdateStatusAsProcessed(gomock.Any(), gomock.Any()).Return(int64(1), nil)
				svc := svcmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				workshop := NewSimpleWorkShop(1, repo, svc)
				s := workshop.(*SimpleWorkShop)
				return s
			},
			ctx:     context.Background(),
			wantErr: nil,
		},
		{
			name: "未抢到数据，不处理",
			mock: func(ctrl *gomock.Controller) *SimpleWorkShop {
				repo := repositorymock.NewMockSmsRepository(ctrl)
				repo.EXPECT().GetFirst(gomock.Any()).Return(domain.Sms{}, repository.ErrCompetitionFailed)
				svc := svcmocks.NewMockService(ctrl)
				workshop := NewSimpleWorkShop(1, repo, svc)
				s := workshop.(*SimpleWorkShop)
				return s
			},
			ctx:     context.Background(),
			wantErr: nil,
		},
		{
			name: "预期外错误",
			mock: func(ctrl *gomock.Controller) *SimpleWorkShop {
				repo := repositorymock.NewMockSmsRepository(ctrl)
				repo.EXPECT().GetFirst(gomock.Any()).Return(domain.Sms{}, errors.New("error not expected"))
				svc := svcmocks.NewMockService(ctrl)
				workshop := NewSimpleWorkShop(1, repo, svc)
				s := workshop.(*SimpleWorkShop)
				return s
			},
			ctx:     context.Background(),
			wantErr: errors.New("error not expected"),
		},
		{
			name: "UpdateStatusAsProcessed",
			mock: func(ctrl *gomock.Controller) *SimpleWorkShop {
				repo := repositorymock.NewMockSmsRepository(ctrl)
				repo.EXPECT().GetFirst(gomock.Any()).Return(createSms(), nil)
				repo.EXPECT().UpdateStatusAsProcessed(gomock.Any(), gomock.Any()).Return(int64(1), errors.New("not expected"))
				svc := svcmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				workshop := NewSimpleWorkShop(1, repo, svc)
				s := workshop.(*SimpleWorkShop)
				return s
			},
			ctx:     context.Background(),
			wantErr: nil,
		},
		{
			name: "重试成功",
			mock: func(ctrl *gomock.Controller) *SimpleWorkShop {
				repo := repositorymock.NewMockSmsRepository(ctrl)
				repo.EXPECT().GetFirst(gomock.Any()).Return(createSms(), nil)
				repo.EXPECT().UpdateStatusAsProcessed(gomock.Any(), gomock.Any()).Return(int64(1), errors.New("not expected"))
				svc := svcmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("send 失败"))
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				strategy, err := retry.NewFixedIntervalRetryStrategy(time.Second, 3)
				require.NoError(t, err)
				workshop := NewSimpleWorkShop(1, repo, svc, WithRetry(strategy))
				s := workshop.(*SimpleWorkShop)
				return s
			},
			ctx:     context.Background(),
			wantErr: nil,
		},
		{
			name: "重试失败",
			mock: func(ctrl *gomock.Controller) *SimpleWorkShop {
				repo := repositorymock.NewMockSmsRepository(ctrl)
				repo.EXPECT().GetFirst(gomock.Any()).Return(createSms(), nil)
				repo.EXPECT().UpdateStatusAsProcessFailed(gomock.Any(), gomock.Any()).Return(int64(1), nil)
				svc := svcmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("send 失败"))
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("send 失败"))
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("send 失败"))
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("send 失败"))

				strategy, err := retry.NewFixedIntervalRetryStrategy(time.Second, 3)
				require.NoError(t, err)
				workshop := NewSimpleWorkShop(1, repo, svc, WithRetry(strategy))
				s := workshop.(*SimpleWorkShop)
				return s
			},
			ctx:     context.Background(),
			wantErr: errors.New("未成功处理任务"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := tc.mock(ctrl)
			err := s.consume(tc.ctx)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func createSms() domain.Sms {
	return domain.Sms{
		Id:         10,
		Tpl:        "tmp",
		Args:       []string{"a", "b"},
		Numbers:    []string{"123", "456"},
		Processing: 0,
		Retry:      0,
		Ctime:      time.Now().Add(-time.Minute * 6),
		Utime:      time.Now().Add(-time.Minute * 6),
	}
}
