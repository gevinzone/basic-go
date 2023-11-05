package startup

import (
	"github.com/gevinzone/basic-go/week7/webook/pkg/logger"
)

func InitLog() logger.LoggerV1 {
	return logger.NewNoOpLogger()
}
