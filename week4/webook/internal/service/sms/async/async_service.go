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
	"fmt"
	"github.com/gevinzone/basic-go/week4/webook/internal/service/sms"
	"github.com/gevinzone/basic-go/week4/webook/internal/service/sms/failover"
	"github.com/gevinzone/basic-go/week4/webook/pkg/ratelimit"
)

type RateLimitFailOverService struct {
	svc     failover.FailureRateFailOverService
	limiter ratelimit.Limiter
}

func NewRateLimitFailOverService(svc failover.FailureRateFailOverService, limiter ratelimit.Limiter) sms.Service {
	return &RateLimitFailOverService{
		svc:     svc,
		limiter: limiter,
	}
}

func (r *RateLimitFailOverService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	if limited, err := r.limiter.Limit(ctx, "sms:tencent"); err != nil || limited {
		go r.asyncHandleSend(ctx, tpl, args, numbers...)
		return errors.New("限流了，稍后重试")
	}
	err := r.svc.Send(ctx, tpl, args, numbers...)
	if err != nil {
		go r.asyncHandleSend(ctx, tpl, args, numbers...)
		return fmt.Errorf("%w, 稍后重试", err)
	}
	return nil
}

func (r *RateLimitFailOverService) asyncHandleSend(ctx context.Context, tpl string, args []string, numbers ...string) error {
	panic("implement later")
}
