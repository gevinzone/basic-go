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

package ginx

import (
	"errors"
	"github.com/gevinzone/basic-go/week5/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func HandleReq[T any](fn func(ctx *gin.Context, req T, uc *jwt.UserClaims) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		err := ctx.Bind(&req)
		if err != nil {
			return
		}

		c, _ := ctx.Get("claims")
		uc, ok := c.(*jwt.UserClaims)
		if !ok {
			zap.L().Error("claims 错误", zap.Error(errors.New("claims 错误")))
			ctx.String(http.StatusOK, "系统错误")
			return
		}

		res, err := fn(ctx, req, uc)
		if err != nil {
			zap.L().Error("系统异常", zap.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}

type Result struct {
	// 这个叫做业务错误码
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
