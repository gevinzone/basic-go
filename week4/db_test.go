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

package week4

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	db, err := initDB()
	if err != nil {
		return
	}
	user := createUser(time.Now().Unix())
	res := db.Create(user)
	assert.NoError(t, res.Error)
	assert.Equal(t, int64(1), res.RowsAffected)
}

func TestInsertBatch(t *testing.T) {
	db, err := initDB()
	if err != nil {
		return
	}
	// batch更多可能会报错
	batch := 10000
	idx := 2000
	users := make([]User, 0, batch)
	for i := 0; i < batch; i++ {
		u := createUser(int64(idx + i))
		users = append(users, *u)
	}
	res := db.Create(&users)
	assert.NoError(t, res.Error)
	assert.Equal(t, int64(batch), res.RowsAffected)
}

func TestInsertBatchMore(t *testing.T) {
	db, err := initDB()
	if err != nil {
		return
	}
	// batch更多可能会报错
	batch := 10000
	idx := 0
	loops := 100
	for loop := 0; loop < loops; loop++ {
		idx = idx + batch
		users := make([]User, 0, batch)
		for i := 0; i < batch; i++ {
			u := createUser(int64(idx + i))
			users = append(users, *u)
		}
		res := db.Create(&users)
		assert.NoError(t, res.Error)
		assert.Equal(t, int64(batch), res.RowsAffected)
	}

}

func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/webook"))
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createUser(id int64) *User {
	now := time.Now().UnixMilli()
	user := &User{
		Id:       id,
		Email:    fmt.Sprintf("gevin%d@gmail.com", id),
		Password: "123",
		Ctime:    now,
		Utime:    now,
	}
	return user
}

type User struct {
	Id       int64  `gorm:"primaryKey, autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	Ctime    int64
	Utime    int64
}
