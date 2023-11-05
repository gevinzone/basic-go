package startup

import (
	"github.com/gevinzone/basic-go/week7/webook/internal/service/oauth2/wechat"
	"github.com/gevinzone/basic-go/week7/webook/pkg/logger"
)

// InitPhantomWechatService 没啥用的虚拟的 wechatService
func InitPhantomWechatService(l logger.LoggerV1) wechat.Service {
	return wechat.NewService("", "", l)
}
