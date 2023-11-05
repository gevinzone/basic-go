package main

import (
	"github.com/gevinzone/basic-go/week7/webook/internal/events"
	"github.com/gin-gonic/gin"
)

type App struct {
	web       *gin.Engine
	consumers []events.Consumer
}
