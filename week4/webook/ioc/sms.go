package ioc

import (
	"github.com/gevinzone/basic-go/week4/webook/internal/service/sms"
	"github.com/gevinzone/basic-go/week4/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	// 换内存，还是换别的
	return memory.NewService()
}
