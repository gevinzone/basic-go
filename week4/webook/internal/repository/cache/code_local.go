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

package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type LocalCodeCache struct {
	c     *cache.Cache
	mutex *sync.RWMutex
}

func NewLocalCodeCache() CodeCache {
	return &LocalCodeCache{
		c:     cache.New(5*time.Minute, 10*time.Minute),
		mutex: &sync.RWMutex{},
	}
}

func (l *LocalCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	key := l.key(biz, phone)
	l.mutex.Lock()
	defer l.mutex.Unlock()
	_, expiration, found := l.c.GetWithExpiration(key)
	zero := time.Time{}
	if found {
		if expiration == zero {
			return errors.New("系统错误")
		}
		if expiration.Unix()-time.Now().Unix() > 540 {
			return ErrCodeSendTooMany
		}
	}
	l.c.Set(key, code, time.Minute*10)
	l.c.Set(l.cntKey(key), 3, time.Minute*10)
	return nil
}

func (l *LocalCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	key := l.key(biz, phone)
	cntKey := l.cntKey(key)
	l.mutex.Lock()
	defer l.mutex.Unlock()
	cntV, expiration, found := l.c.GetWithExpiration(cntKey)
	if !found {
		return false, ErrUnknownForCode
	}
	cnt, ok := cntV.(int)
	if !ok || cnt <= 0 {
		return false, ErrCodeVerifyTooManyTimes
	}
	codeV, _ := l.c.Get(key)
	code := codeV.(string)
	if code != inputCode {
		l.c.Set(cntKey, cnt-1, expiration.Sub(time.Now()))
		return false, nil
	}
	l.c.Set(cntKey, -1, expiration.Sub(time.Now()))
	return true, nil

}

func (l *LocalCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (l *LocalCodeCache) cntKey(key string) string {
	return fmt.Sprintf("%s:cnt", key)
}
