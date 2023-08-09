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

package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gevinzone/basic-go/week2/webook/internal/domain"
	"github.com/gevinzone/basic-go/week2/webook/internal/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/logout", u.Logout)
	ug.POST("/edit", u.Edit)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		// 记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符")
		return
	}
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	_ = sess.Save()
	ctx.String(http.StatusOK, "登录成功")
}

func (u *UserHandler) Logout(ctx *gin.Context) {

}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type request struct {
		UserId   int64  `json:"user_id"`
		Nickname string `json:"nickname"`
		Biology  string `json:"biology"`
		Birthday string `json:"birthday"`
	}

	req := request{}
	err := ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusBadRequest, "输入数据格式不对")
		return
	}
	p, err := func(r request) (domain.Profile, error) {
		var p domain.Profile
		if len([]rune(r.Nickname)) > 64 {
			return p, errors.New("nickname 过长")
		}
		if len([]rune(r.Biology)) > 256 {
			return p, errors.New("biology 过长")
		}
		t, er := time.Parse("2006-01-02", r.Birthday)
		if er != nil {
			return p, er
		}
		p = domain.Profile{
			UserId:   r.UserId,
			Nickname: r.Nickname,
			Biology:  r.Biology,
			Birthday: t,
		}

		return p, nil
	}(req)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	err = u.svc.EditProfile(ctx, p)
	if err != nil {
		ctx.String(http.StatusBadRequest, "输入数据不对")
		return
	}
	ctx.JSON(http.StatusOK, p)
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	email := ctx.Query("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "url 缺少email 参数"})
		return
	}
	profile, err := u.svc.GetProfileByEmail(ctx, email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "系统错误")
		return
	}
	type response struct {
		domain.Profile
		Birthday string
	}
	toResponse := func(p domain.Profile) response {
		return response{
			Profile:  p,
			Birthday: p.Birthday.Format("2006-01-02"),
		}
	}
	ctx.JSON(http.StatusOK, toResponse(profile))
}

type profileRes struct {
	Nickname string
	Biology  string
	Birthday time.Time
}
