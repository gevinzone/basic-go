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

package repository

import (
	"context"
	"github.com/gevinzone/basic-go/week4/webook/internal/domain"
	"github.com/gevinzone/basic-go/week4/webook/internal/repository/dao"
	"strings"
	"time"
)

const Delimiter = "&^"

var ErrCompetitionFailed = dao.ErrCompetitionFailed

type SmsRepository interface {
	GetFirst(ctx context.Context) (domain.Sms, error)
	SaveSms(ctx context.Context, sms domain.Sms) (domain.Sms, error)
}

type SmsDbRepository struct {
	smsDao dao.SmsDao
}

func NewSmsDbRepository(dao dao.SmsDao) SmsRepository {
	return &SmsDbRepository{
		smsDao: dao,
	}
}

var _ SmsRepository = (*SmsDbRepository)(nil)

func (s *SmsDbRepository) GetFirst(ctx context.Context) (domain.Sms, error) {
	entity, err := s.smsDao.GetFirst(ctx)
	if err != nil {
		return domain.Sms{}, err
	}
	return toDomain(entity), nil
}

func (s *SmsDbRepository) SaveSms(ctx context.Context, sms domain.Sms) (domain.Sms, error) {
	entity, err := s.smsDao.SaveSms(ctx, sms)
	if err != nil {
		return domain.Sms{}, err
	}
	return toDomain(entity), nil
}

func toDomain(entity dao.Sms) domain.Sms {
	return domain.Sms{
		Id:         entity.Id,
		Tpl:        entity.Tpl,
		Args:       strings.Split(entity.Args, Delimiter),
		Numbers:    strings.Split(entity.Numbers, Delimiter),
		Processing: entity.Processing,
		Retry:      entity.Retry,
		Ctime:      time.UnixMilli(entity.Ctime),
		Utime:      time.UnixMilli(entity.Utime),
	}
}
