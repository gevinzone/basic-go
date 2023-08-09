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
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) Insert(ctx context.Context, u User) (User, error) {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			// 邮箱冲突
			return u, ErrUserDuplicateEmail
		}
	}
	return u, err
}

type ProfileDAO struct {
	db *gorm.DB
}

func NewProfileDAO(db *gorm.DB) *ProfileDAO {
	return &ProfileDAO{db: db}
}

func (dao *ProfileDAO) FindByUserId(ctx context.Context, id int64) (Profile, error) {
	var p Profile
	err := dao.db.WithContext(ctx).Where("user_id=?", id).First(&p).Error
	return p, err
}

func (dao *ProfileDAO) Insert(ctx context.Context, p Profile) error {
	now := time.Now().UnixMilli()
	p.Ctime = now
	p.Utime = now
	return dao.db.WithContext(ctx).Create(&p).Error
}

func (dao *ProfileDAO) Update(ctx context.Context, p Profile) error {
	p.Utime = time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&p).Where("user_id=?", p.UserId).Updates(p).Error
}

// User 直接对应数据库表结构
// 有些人叫做 entity，有些人叫做 model，有些人叫做 PO(persistent object)
type User struct {
	Id       int64  `gorm:"primaryKey, autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	Ctime    int64
	Utime    int64
}

type Profile struct {
	Id       int64 `gorm:"primaryKey, autoIncrement"`
	UserId   int64
	Nickname string
	Biology  string
	Birthday int64
	Ctime    int64
	Utime    int64
}
