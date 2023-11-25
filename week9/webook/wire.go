//go:build wireinject

package main

import (
	"github.com/gevinzone/basic-go/week9/webook/internal/events/article"
	"github.com/gevinzone/basic-go/week9/webook/internal/repository"
	article2 "github.com/gevinzone/basic-go/week9/webook/internal/repository/article"
	"github.com/gevinzone/basic-go/week9/webook/internal/repository/cache"
	"github.com/gevinzone/basic-go/week9/webook/internal/repository/dao"
	article3 "github.com/gevinzone/basic-go/week9/webook/internal/repository/dao/article"
	"github.com/gevinzone/basic-go/week9/webook/internal/service"
	"github.com/gevinzone/basic-go/week9/webook/internal/web"
	ijwt "github.com/gevinzone/basic-go/week9/webook/internal/web/jwt"
	"github.com/gevinzone/basic-go/week9/webook/ioc"
	"github.com/google/wire"
)

var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCachedInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

var rankingServiceSet = wire.NewSet(
	repository.NewCachedRankingRepository,
	cache.NewRankingRedisCache,
	service.NewBatchRankingService,
)

func InitWebServer() *App {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.NewConsumers,
		ioc.NewSyncProducer,

		interactiveSvcProvider,
		rankingServiceSet,
		ioc.InitJobs,
		ioc.InitRankingJob,

		// consumer
		article.NewInteractiveReadEventBatchConsumer,
		article.NewKafkaProducer,

		// 初始化 DAO
		dao.NewUserDAO,
		article3.NewGORMArticleDAO,

		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,
		article2.NewArticleRepository,

		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,

		// 直接基于内存实现
		ioc.InitSMSService,
		ioc.InitWechatService,

		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,
		//ioc.NewWechatHandlerConfig,
		ijwt.NewRedisJWTHandler,
		// 你中间件呢？
		// 你注册路由呢？
		// 你这个地方没有用到前面的任何东西
		//gin.Default,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
		// 组装我这个结构体的所有字段
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
