package startup

import (
	"github.com/gevinzone/basic-go/week9/webook/pkg/logger"
)

func InitLog() logger.LoggerV1 {
	return logger.NewNoOpLogger()
}
