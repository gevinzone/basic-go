//gobuild e2e

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
	"github.com/gevinzone/basic-go/week4/webook/internal/domain"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

type SmsSuite struct {
	suite.Suite
	db  *gorm.DB
	dsn string
}

func (s *SmsSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(s.dsn))
	s.Require().NoError(err)
	s.db = db
	err = db.AutoMigrate(&Sms{})
	s.Require().NoError(err)
}

func (s *SmsSuite) SetupTest() {

}

func (s *SmsSuite) TearDownTest() {
	s.db.Exec("truncate table `sms`")
}

func (s *SmsSuite) TearDownSuite() {
	//s.db.Raw("truncate sms")
}

func TestSmsSuite(t *testing.T) {
	suite.Run(t, &SmsSuite{
		dsn: "root:root@tcp(localhost:3306)/webook2",
	})
}

func (s *SmsSuite) prepareDataNotProcessingData() *Sms {
	now := time.Now().UnixMilli()
	sms := &Sms{
		Id:         10000,
		Tpl:        "default",
		Args:       "default",
		Numbers:    "1",
		Processing: domain.SmsUnprocessed,
		Retry:      0,
		Ctime:      now,
		Utime:      now,
	}
	s.db.WithContext(context.Background()).Create(sms)
	return sms
}

func (s *SmsSuite) prepareDataProcessingData() *Sms {
	now := time.Now().UnixMilli()
	sms := &Sms{
		Id:         10000,
		Tpl:        "default",
		Args:       "default",
		Numbers:    "1",
		Processing: domain.SmsProcessing,
		Retry:      0,
		Ctime:      now,
		Utime:      now,
	}
	s.db.WithContext(context.Background()).Create(sms)
	return sms
}

func (s *SmsSuite) TestGetOne() {
	s.prepareDataNotProcessingData()
	dao := NewSmsGormDao(s.db)
	sms, err := dao.GetFirst(context.Background())
	s.Require().NoError(err)
	err = s.db.WithContext(context.Background()).First(&sms).Error
	s.Require().NoError(err)
	s.Assert().Equal(domain.SmsProcessing, sms.Processing)

	sms, err = dao.GetFirst(context.Background())
	s.Assert().Equal(gorm.ErrRecordNotFound, err)
	s.Assert().Equal(Sms{}, sms)
}

func (s *SmsSuite) TestSaveSms() {
	bo := domain.Sms{
		Tpl:        "tmp",
		Args:       []string{"tmp"},
		Numbers:    []string{"123"},
		Processing: domain.SmsUnprocessed,
		Ctime:      time.Now(),
		Utime:      time.Now(),
	}
	dao := NewSmsGormDao(s.db)
	sms, err := dao.SaveSms(context.Background(), bo)
	s.Require().NoError(err)
	s.Assert().Equal(int64(1), sms.Id)
	wantSms := fromDomain(bo)
	wantSms.Id = sms.Id
	s.Assert().Equal(wantSms, sms)
}

func (s *SmsSuite) TestUpdateStatusAsProcessed() {
	sms := s.prepareDataProcessingData()
	dao := NewSmsGormDao(s.db)
	affected, err := dao.UpdateStatusAsProcessed(context.Background(), sms.Id)
	s.Assert().NoError(err)
	s.Assert().Equal(int64(1), affected)

	err = s.db.WithContext(context.Background()).First(sms).Error
	s.Require().NoError(err)
	s.Assert().Equal(domain.SmsProcessed, sms.Processing)
}

func (s *SmsSuite) TestUpdateStatusAsProcessFailed() {
	sms := s.prepareDataProcessingData()
	dao := NewSmsGormDao(s.db)
	affected, err := dao.UpdateStatusAsProcessFailed(context.Background(), sms.Id)
	s.Assert().NoError(err)
	s.Assert().Equal(int64(1), affected)

	err = s.db.WithContext(context.Background()).First(sms).Error
	s.Require().NoError(err)
	s.Assert().Equal(domain.SmsProcessFailed, sms.Processing)
}
