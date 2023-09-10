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
	"time"
)

type TimeoutFailOverService struct {
	svcs      []sms.Service
	threshold int64
	cnt       int64
	idx       int64
}

var _ sms.Service = (*TimeoutFailOverService)(nil)

func NewTimeoutFailOverService(svcs []sms.Service, threshold int64) sms.Service {
	return &TimeoutFailOverService{
		svcs:      svcs,
		threshold: threshold,
	}
}

func (s *TimeoutFailOverService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt64(&s.idx)
	cnt := atomic.LoadInt64(&s.cnt)
	if cnt > s.threshold {
		newIdx := (idx + 1) % int64(len(s.svcs))
		if atomic.CompareAndSwapInt64(&s.idx, idx, newIdx) {
			atomic.StoreInt64(&s.cnt, 0)
		} else {
			// cas 不成功的并发请求，等上面重置完成再继续
			// 如果没有这个分支，虽然并发不安全，但对本业务影响比较小，也可以忽略本分支
			time.Sleep(time.Millisecond)
		}
		idx = atomic.LoadInt64(&s.idx)
	}
	err := s.svcs[idx].Send(ctx, tpl, args, numbers...)
	switch {
	case err == nil:
		atomic.StoreInt64(&s.cnt, 0)
		return nil
	case errors.Is(err, context.DeadlineExceeded):
		atomic.AddInt64(&s.cnt, 1)
		return err
	default:
		return err
	}
}
