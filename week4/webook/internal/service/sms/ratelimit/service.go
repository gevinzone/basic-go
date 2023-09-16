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

package ratelimit

import (
	"context"
	"fmt"
	"github.com/gevinzone/basic-go/week4/webook/internal/service/sms"
	"github.com/gevinzone/basic-go/week4/webook/pkg/ratelimit"
)

var ErrLimited = fmt.Errorf("触发了限流")

type SmsService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewSmsService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &SmsService{
		svc:     svc,
		limiter: limiter,
	}
}

func (s *SmsService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil || limited {
		return ErrLimited
	}
	return s.svc.Send(ctx, tpl, args, numbers...)
}
