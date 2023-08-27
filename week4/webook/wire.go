//go:build wireinject

package main

import (
	"github.com/gevinzone/basic-go/week4/webook/internal/repository"
	"github.com/gevinzone/basic-go/week4/webook/internal/repository/cache"
	"github.com/gevinzone/basic-go/week4/webook/internal/repository/dao"
	"github.com/gevinzone/basic-go/week4/webook/internal/service"
	"github.com/gevinzone/basic-go/week4/webook/internal/web"
	"github.com/gevinzone/basic-go/week4/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,

		// 初始化 DAO
		dao.NewUserDAO,

		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewUserService,
		service.NewCodeService,
		// 直接基于内存实现
		ioc.InitSMSService,
		web.NewUserHandler,
		// 你中间件呢？
		// 你注册路由呢？
		// 你这个地方没有用到前面的任何东西
		//gin.Default,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
