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

type FailureRateFailOverService struct {
	svcs       []sms.Service
	idx        int64
	failureCnt int64
	successCnt int64
	rate       float64
}

var _ sms.Service = (*FailureRateFailOverService)(nil)

func NewFailureRateFailOverService(svcs []sms.Service, rate float64) sms.Service {
	if svcs == nil || len(svcs) == 0 || rate > 1 || rate < 0 {
		panic("预期之外的参数错误")
	}
	return &FailureRateFailOverService{
		svcs: svcs,
		rate: rate,
	}
}

func (f *FailureRateFailOverService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt64(&f.idx)
	failureCnt := atomic.LoadInt64(&f.failureCnt)
	successCnt := atomic.LoadInt64(&f.successCnt)
	if float64(failureCnt)/float64(failureCnt+successCnt) >= f.rate {
		newIdx := (idx + 1) % int64(len(f.svcs))
		if atomic.CompareAndSwapInt64(&f.idx, idx, newIdx) {
			atomic.StoreInt64(&f.failureCnt, 0)
			atomic.StoreInt64(&f.successCnt, 0)
		} else {
			// cas 不成功的并发请求，等上面重置完成再继续
			// 如果没有这个分支，虽然有并发不安全隐患，但几乎不可能发生，且即便发生，对本业务影响比较小，可以忽略本分支
			time.Sleep(time.Millisecond)
		}
		idx = atomic.LoadInt64(&f.idx)
	}
	err := f.svcs[idx].Send(ctx, tpl, args, numbers...)
	switch {
	case err == nil:
		atomic.AddInt64(&f.successCnt, 1)
		return nil
	case errors.Is(err, context.DeadlineExceeded):
		atomic.AddInt64(&f.failureCnt, 1)
		return err
	default:
		return err
	}
}
