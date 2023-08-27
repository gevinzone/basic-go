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
	"github.com/gevinzone/basic-go/live/webook/internal/domain"
	"github.com/gevinzone/basic-go/live/webook/internal/repository/cache"
	"github.com/gevinzone/basic-go/live/webook/internal/repository/dao"
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
	UpdateProfile(ctx context.Context, profile domain.Profile) error
	FindProfileByEmail(ctx context.Context, email string) (domain.Profile, error)
}

type CachedUserRepository struct {
	db         *gorm.DB
	userDAO    dao.UserDAO
	profileDAO dao.ProfileDAO
	cache      cache.UserCache
}

func NewUserRepository(db *gorm.DB, userDao dao.UserDAO, profileDao dao.ProfileDAO, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		db:         db,
		userDAO:    userDao,
		profileDAO: profileDao,
		cache:      c,
	}
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.userDAO.FindByEmail(ctx, email)
	if err != nil {
		var user domain.User
		return user, err
	}
	return r.userEntityToDomain(u), nil
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.userDAO.FindByPhone(ctx, phone)
	if err != nil {
		var user domain.User
		return user, err
	}
	return r.userEntityToDomain(u), nil
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
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}
	// 没这个数据
	//if err == cache.ErrKeyNotExist {
	// 去数据库里面加载
	//}

	ue, err := r.userDAO.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = r.userEntityToDomain(ue)

	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 我这里怎么办？
			// 打日志，做监控
			//return domain.User{}, err
		}
	}()
	return u, err
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

func (r *CachedUserRepository) userEntityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}
