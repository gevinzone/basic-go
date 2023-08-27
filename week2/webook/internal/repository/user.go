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
	"database/sql"
	"github.com/gevinzone/basic-go/week2/webook/internal/domain"
	"github.com/gevinzone/basic-go/week2/webook/internal/repository/cache"
	"github.com/gevinzone/basic-go/week2/webook/internal/repository/dao"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
}

type CachedUserRepository struct {
	db         *gorm.DB
	userDAO    *dao.GormUserDAO
	profileDAO *dao.GORMProfileDAO
	cache      *cache.UserCache
}

func NewUserRepository(db *gorm.DB) *CachedUserRepository {
	return &CachedUserRepository{
		db:         db,
		userDAO:    dao.NewUserDAO(db),
		profileDAO: dao.NewProfileDAO(db),
	}
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.userDAO.FindByEmail(ctx, email)
	if err != nil {
		var user domain.User
		return user, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
	}, nil
}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var (
			user dao.User
			err  error
		)
		userDAO := dao.NewUserDAO(tx)
		if user, err = userDAO.Insert(ctx, r.userDomainToEntity(u)); err != nil {
			return err
		}

		profileDAO := dao.NewProfileDAO(tx)
		err = profileDAO.Insert(ctx, dao.Profile{UserId: user.Id, Birthday: time.Now().UnixMilli()})
		return err
	})
}

func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 先从 cache 里面找
	// 再从 dao 里面找
	// 找到了回写 cache
	panic("implement me")
}

func (r *CachedUserRepository) FindProfileByEmail(ctx context.Context, email string) (domain.Profile, error) {
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
		userDAO := dao.NewUserDAO(tx)
		profileDAO := dao.NewProfileDAO(tx)
		if u, er = userDAO.FindByEmail(ctx, email); er != nil {
			return er
		}
		if p, er = profileDAO.FindByUserId(ctx, u.Id); er != nil {
			return er
		}
		profile = domain.Profile{
			UserId:   u.Id,
			Email:    u.Email.String,
			Nickname: p.Nickname,
			Biology:  p.Biology,
			Birthday: time.UnixMilli(p.Birthday),
			Ctime:    time.UnixMilli(p.Ctime),
			Utime:    time.UnixMilli(p.Utime),
		}
		return nil
	})
	return profile, err
}

func (r *CachedUserRepository) UpdateProfile(ctx context.Context, profile domain.Profile) error {
	return r.profileDAO.Update(ctx, dao.Profile{
		UserId:   profile.UserId,
		Birthday: profile.Birthday.UnixMilli(),
		Biology:  profile.Biology,
		Nickname: profile.Nickname,
	})
}

func (r *CachedUserRepository) userDomainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			// 我确实有手机号
			Valid: u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Ctime:    u.Ctime.UnixMilli(),
	}
}
