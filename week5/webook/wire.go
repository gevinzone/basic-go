//go:build wireinject

package main

import (
	"github.com/gevinzone/basic-go/week5/webook/internal/repository"
	"github.com/gevinzone/basic-go/week5/webook/internal/repository/article"
	"github.com/gevinzone/basic-go/week5/webook/internal/repository/cache"
	"github.com/gevinzone/basic-go/week5/webook/internal/repository/dao"
	"github.com/gevinzone/basic-go/week5/webook/internal/service"
	"github.com/gevinzone/basic-go/week5/webook/internal/web"
	ijwt "github.com/gevinzone/basic-go/week5/webook/internal/web/jwt"
	"github.com/gevinzone/basic-go/week5/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,

		// 初始化 DAO
		dao.NewUserDAO,
		dao.NewGORMArticleDAO,

		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,
		article.NewCachedArticleRepository,

		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
		// 直接基于内存实现
		ioc.InitSMSService,
		ioc.InitWechatService,

		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		ioc.NewWechatHandlerConfig,
		ijwt.NewRedisJWTHandler,
		// 你中间件呢？
		// 你注册路由呢？
		// 你这个地方没有用到前面的任何东西
		//gin.Default,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
