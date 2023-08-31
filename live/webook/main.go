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

package main

import (
	"github.com/gevinzone/basic-go/live/webook/config"
	"github.com/gevinzone/basic-go/live/webook/internal/repository"
	"github.com/gevinzone/basic-go/live/webook/internal/repository/cache"
	"github.com/gevinzone/basic-go/live/webook/internal/repository/dao"
	"github.com/gevinzone/basic-go/live/webook/internal/service"
	"github.com/gevinzone/basic-go/live/webook/internal/web"
	"github.com/gevinzone/basic-go/live/webook/internal/web/middleware"
	"github.com/gevinzone/basic-go/live/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func main() {
	db := initDB()
	redisClient := initRedis()
	server := initWebServer(redisClient)
	u := initUserHandler(db, redisClient)
	u.RegisterRoutes(server)
	//err := server.Run(":8000")
	//if err != nil {
	//	panic(err)
	//}
	//server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，你来了")
	})

	server.Run(":8080")
}

func initWebServer(redisClient redis.Cmdable) *gin.Engine {
	server := gin.Default()

	server.Use(func(ctx *gin.Context) {
		println("this is the first middleware")
	})
	server.Use(func(ctx *gin.Context) {
		println("this is the second middleware")
	})

	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 你不加这个，前端是拿不到的
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your_domain.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	server.Use(sessions.Sessions("gevin_session", memstore.NewStore([]byte("this is secret"))))
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/signup", "/users/login").
	//	Build())
	server.Use(middleware.NewLoginJwtMiddlewareBuilder().
		IgnorePaths("/users/signup", "/users/login", "/hello").
		Build())
	return server
}

func initUserHandler(db *gorm.DB, client redis.Cmdable) *web.UserHandler {
	userDao := dao.NewUserDAO(db)
	profileDao := dao.NewProfileDAO(db)
	userProfileDao := dao.NewUserWithProfileDAO(db, userDao, profileDao)
	c := cache.NewUserCache(client)
	return web.NewUserHandler(
		service.NewUserService(
			repository.NewUserRepository(userDao, profileDao, userProfileDao, c)))
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
}
