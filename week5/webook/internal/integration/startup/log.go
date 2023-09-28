package startup

import (
	"github.com/gevinzone/basic-go/week5/webook/pkg/logger"
)

func InitLog() logger.LoggerV1 {
	return &logger.NopLogger{}
}
