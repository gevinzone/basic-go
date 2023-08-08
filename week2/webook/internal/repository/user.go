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
	"github.com/gevinzone/basic-go/week2/webook/internal/domain"
	"github.com/gevinzone/basic-go/week2/webook/internal/repository/dao"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	db         *gorm.DB
	userDAO    *dao.UserDAO
	profileDAO *dao.ProfileDAO
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db:         db,
		userDAO:    dao.NewUserDAO(db),
		profileDAO: dao.NewProfileDAO(db),
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.userDAO.FindByEmail(ctx, email)
	if err != nil {
		var user domain.User
		return user, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var (
			user dao.User
			err  error
		)
		userDAO := dao.NewUserDAO(tx)
		if user, err = userDAO.Insert(ctx, dao.User{
			Email:    u.Email,
			Password: u.Password,
		}); err != nil {
			return err
		}

		profileDAO := dao.NewProfileDAO(tx)
		err = profileDAO.Insert(ctx, dao.Profile{UserId: user.Id, Birthday: time.Now()})
		return err
	})
}

func (r *UserRepository) FindById(int64) {
	// 先从 cache 里面找
	// 再从 dao 里面找
	// 找到了回写 cache
}

func (r *UserRepository) FindProfileByEmail(ctx context.Context, email string) (domain.Profile, error) {
	var (
		profile domain.Profile
		err     error
	)
	err = r.db.Transaction(func(tx *gorm.DB) error {
		var (
			u  dao.User
			p  dao.Profile
			er error
		)
		if u, er = r.userDAO.FindByEmail(ctx, email); er != nil {
			return er
		}
		if p, er = r.profileDAO.FindByUserId(ctx, u.Id); er != nil {
			return er
		}
		profile = domain.Profile{
			UserId:   u.Id,
			Email:    u.Email,
			Nickname: p.Nickname,
			Biology:  p.Biology,
			Birthday: p.Birthday,
			Ctime:    time.UnixMilli(p.Ctime),
			Utime:    time.UnixMilli(p.Utime),
		}
		return nil
	})
	return profile, err
}
