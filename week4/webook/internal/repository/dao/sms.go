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

package dao

import (
	"context"
	"errors"
	"github.com/gevinzone/basic-go/week4/webook/internal/domain"
	"github.com/gevinzone/basic-go/week4/webook/internal/repository"
	"gorm.io/gorm"
	"strings"
	"time"
)

var ErrCompetitionFailed = errors.New("未抢到记录")

type SmsDao interface {
	GetFirst(ctx context.Context) (Sms, error)
	SaveSms(ctx context.Context, sms domain.Sms) (Sms, error)
	UpdateStatusAsProcessed(ctx context.Context, id int64) (int64, error)
	UpdateStatusAsProcessFailed(ctx context.Context, id int64) (int64, error)
}

type SmsGormDao struct {
	db *gorm.DB
}

var _ SmsDao = (*SmsGormDao)(nil)

func NewSmsGormDao(db *gorm.DB) SmsDao {
	return &SmsGormDao{
		db: db,
	}
}

func (dao *SmsGormDao) GetFirst(ctx context.Context) (Sms, error) {
	var s Sms
	err := dao.db.WithContext(ctx).Where("processing=?", domain.SmsUnprocessed).First(&s).Error
	if err != nil {
		return Sms{}, err
	}
	s.Processing = domain.SmsProcessing
	s.Utime = time.Now().UnixMilli()
	affected := dao.db.WithContext(ctx).Updates(&s).
		Where("id=? AND processing=?", s.Id, domain.SmsUnprocessed).RowsAffected
	if affected != 1 {
		return Sms{}, ErrCompetitionFailed
	}
	return s, nil
}

func (dao *SmsGormDao) SaveSms(ctx context.Context, sms domain.Sms) (Sms, error) {
	s := fromDomain(sms)
	err := dao.db.WithContext(ctx).Create(&s).Error
	return s, err
}

func (dao *SmsGormDao) UpdateStatusAsProcessed(ctx context.Context, id int64) (int64, error) {
	res := dao.db.WithContext(ctx).Update("processing", domain.SmsProcessed).Where("`id`=? AND `processing`=?", id, domain.SmsProcessing)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (dao *SmsGormDao) UpdateStatusAsProcessFailed(ctx context.Context, id int64) (int64, error) {
	res := dao.db.WithContext(ctx).Update("processing", domain.SmsProcessFailed).Where("`id`=? AND `processing`=?", id, domain.SmsProcessing)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

type Sms struct {
	Id         int64
	Tpl        string
	Args       string
	Numbers    string
	Processing int
	Retry      int
	Ctime      int64
	Utime      int64
}

func fromDomain(bo domain.Sms) Sms {
	return Sms{
		Id:         bo.Id,
		Tpl:        bo.Tpl,
		Args:       strings.Join(bo.Args, repository.Delimiter),
		Numbers:    strings.Join(bo.Numbers, repository.Delimiter),
		Processing: bo.Processing,
		Retry:      bo.Retry,
		Ctime:      bo.Ctime.UnixMilli(),
		Utime:      bo.Utime.UnixMilli(),
	}
}
