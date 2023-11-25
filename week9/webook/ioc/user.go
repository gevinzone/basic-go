package ioc

import (
	"github.com/gevinzone/basic-go/week9/webook/internal/repository/cache"
	"github.com/gevinzone/basic-go/week9/webook/pkg/redisx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

//func InitUserHandler(repo repository.UserRepository) service.UserService {
//	l, err := zap.NewDevelopment()
//	if err != nil {
//		panic(err)
//	}
//	return service.NewUserService(repo, )
//}

// InitUserCache 配合 PrometheusHook 使用
func InitUserCache(client *redis.ClusterClient) cache.UserCache {
	client.AddHook(redisx.NewPrometheusHook(
		prometheus.SummaryOpts{
			Namespace: "geekbang_daming",
			Subsystem: "webook",
			Name:      "gin_http",
			Help:      "统计 GIN 的 HTTP 接口",
			ConstLabels: map[string]string{
				"biz": "user",
			},
		}))
	panic("你别调用")
}
