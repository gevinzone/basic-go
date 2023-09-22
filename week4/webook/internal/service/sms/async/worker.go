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
	"github.com/ecodeclub/ekit/retry"
	"github.com/gevinzone/basic-go/week4/webook/internal/repository"
	"github.com/gevinzone/basic-go/week4/webook/internal/service/sms"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type agent func(ctx context.Context)

type Workshop interface {
	Start(ctx context.Context)
	IsStarted() bool
}

type SimpleWorkShop struct {
	started   bool
	agents    []agent
	agentCnt  int
	smsRepo   repository.SmsRepository
	svc       sms.Service
	awaitTime time.Duration
	retry     retry.Strategy
}

var _ Workshop = (*SimpleWorkShop)(nil)

func NewSimpleWorkShop(agentCnt int, repo repository.SmsRepository, svc sms.Service) Workshop {
	// todo: option 模式
	strategy, _ := retry.NewFixedIntervalRetryStrategy(time.Second*30, 5)
	res := &SimpleWorkShop{
		started:   false,
		agentCnt:  agentCnt,
		smsRepo:   repo,
		svc:       svc,
		awaitTime: time.Minute * 5,
		retry:     strategy,
	}

	agents := make([]agent, 0, agentCnt)
	for i := 0; i < agentCnt; i++ {
		agents = append(agents, res.createAgent())
	}
	res.agents = agents
	return res
}

func (w *SimpleWorkShop) IsStarted() bool {
	return w.started
}

func (w *SimpleWorkShop) Start(ctx context.Context) {
	for _, a := range w.agents {
		go a(ctx)
	}
	w.started = true
}

func (w *SimpleWorkShop) createAgent() agent {
	return func(ctx context.Context) {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := w.consume(ctx)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}

func (w *SimpleWorkShop) consume(ctx context.Context) error {
	s, err := w.smsRepo.GetFirst(ctx)
	if err == repository.ErrCompetitionFailed {
		return nil
	}
	if err != nil {
		return err
	}
	d := time.Now().Sub(s.Utime)
	if d < w.awaitTime {
		time.Sleep(w.awaitTime - d)
	}

	interval, canRetry := time.Duration(0), true
	for canRetry {
		time.Sleep(interval)
		err = w.svc.Send(ctx, s.Tpl, s.Args, s.Numbers...)
		if err == nil {
			affect, er := w.smsRepo.UpdateStatusAsProcessed(ctx, s.Id)
			if er != nil || affect != 1 {
				log.Error(errors.New("发生短信成功，但更新数据库状态失败"))
			}
			return nil
		}

		interval, canRetry = w.retry.Next()
	}

	affect, er := w.smsRepo.UpdateStatusAsProcessFailed(ctx, s.Id)
	if er != nil || affect != 1 {
		log.Error(errors.New("更新数据库状态失败"))
	}
	return errors.New("未成功处理任务")
}
