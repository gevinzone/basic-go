//go:build wireinject

package startup

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

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)
var userSvcProvider = wire.NewSet(
	dao.NewUserDAO,
	cache.NewUserCache,
	repository.NewUserRepository,
	service.NewUserService)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdProvider,
		userSvcProvider,
		//articlSvcProvider,
		cache.NewCodeCache,
		dao.NewGORMArticleDAO,
		repository.NewCodeRepository,
		article.NewArticleRepository,
		// service 部分
		// 集成测试我们显式指定使用内存实现
		ioc.InitSMSService,

		// 指定啥也不干的 wechat service
		InitPhantomWechatService,
		service.NewCodeService,
		service.NewArticleService,
		// handler 部分
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		InitWechatHandlerConfig,
		ijwt.NewRedisJWTHandler,

		// gin 的中间件
		ioc.InitMiddlewares,

		// Web 服务器
		ioc.InitWebServer,
	)
	// 随便返回一个
	return gin.Default()
}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(thirdProvider,
		dao.NewGORMArticleDAO,
		service.NewArticleService,
		web.NewArticleHandler,
		article.NewArticleRepository,
	)
	return &web.ArticleHandler{}
}

func InitUserSvc() service.UserService {
	wire.Build(thirdProvider, userSvcProvider)
	return service.NewUserService(nil, nil)
}

func InitJwtHdl() ijwt.Handler {
	wire.Build(thirdProvider, ijwt.NewRedisJWTHandler)
	return ijwt.NewRedisJWTHandler(nil)
}
