package ioc

import (
	"github.com/gevinzone/basic-go/week7/webook/internal/service/sms"
	"github.com/gevinzone/basic-go/week7/webook/internal/service/sms/memory"
	"github.com/redis/go-redis/v9"
)

func InitSMSService(cmd redis.Cmdable) sms.Service {
	// 换内存，还是换别的
	//svc := ratelimit.NewRatelimitSMSService(memory.NewService(),
	//	limiter.NewRedisSlidingWindowLimiter(cmd, time.Second, 100))
	//return retryable.NewService(svc, 3)
	return memory.NewService()
}
