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
	"github.com/gevinzone/basic-go/week4/webook/internal/service/sms"
	"sync/atomic"
)

type InorderService struct {
	svcs []sms.Service
	idx  int64
}

var _ sms.Service = (*InorderService)(nil)

func NewInorderService(svcs []sms.Service) sms.Service {
	return &StartAnewService{svcs: svcs}
}

func (s *InorderService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	length := int64(len(s.svcs))
	old := atomic.LoadInt64(&s.idx)
	idx := atomic.AddInt64(&s.idx, 1)
	for err := s.svcs[idx].Send(ctx, tpl, args, numbers...); ; {
		if err == nil {
			return nil
		}
		idx = atomic.AddInt64(&s.idx, 1)
		if idx >= old+length {
			return errors.New("全部服务商都失败了")
		}
	}
}
