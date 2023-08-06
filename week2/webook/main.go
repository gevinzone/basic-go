// Copyright 2023 igevin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	server := initWebServer()
	err := server.Run(":8000")
	if err != nil {
		panic(err)
	}
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(func(ctx *gin.Context) {
		println("this is the first middleware")
	})
	server.Use(func(ctx *gin.Context) {
		println("this is the second middleware")
	})
	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your_domain.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	return server
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("gevin:gevin@tcp(localhost:3306)/webook"))
	if err != nil {
		panic(err)
	}

	return db
}
