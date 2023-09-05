package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gevinzone/basic-go/week4/webook/internal/domain"
	svcmocks "github.com/gevinzone/basic-go/week4/webook/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/gevinzone/basic-go/week4/webook/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestEncrypt(t *testing.T) {
	password := "hello#world123"
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
	assert.NoError(t, err)
}

func TestNil(t *testing.T) {
	testTypeAssert(nil)
}

func testTypeAssert(c any) {
	claims := c.(*UserClaims)
	println(claims.Uid)
}

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				svc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "abc123@email.com",
					Password: "abc#123@TAG",
				}).Return(nil)
				return svc
			},
			reqBody: `
{
	"email": "abc123@email.com",
	"confirmPassword": "abc#123@TAG",
	"password": "abc#123@TAG"
}
`,
			wantBody: "注册成功",
			wantCode: 200,
		},
		{
			name: "json bind 错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				return svc
			},
			reqBody: `
{
	"email": "abc123@email.com",
	"confirmPassword": "abc#123@TAG",
	"passworda": "abc#123@TAG",
}
`,
			wantCode: 400,
		}, {
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				svc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "abc123@email.com",
					Password: "abc#123@TAG",
				}).Return(nil)
				return svc
			},
			reqBody: `
{
	"email": "abc123@email.com",
	"confirmPassword": "abc#123@TAG",
	"password": "abc#123@TAG"
}
`,
			wantBody: "注册成功",
			wantCode: 200,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				return svc
			},
			reqBody: `
{
	"email": "abc123@@email.com",
	"confirmPassword": "abc#123@TAG",
	"password": "abc#123@TAG"
}
`,
			wantBody: "你的邮箱格式不对",
			wantCode: 200,
		},
		{
			name: "两次输入的密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				return svc
			},
			reqBody: `
{
	"email": "abc123@email.com",
	"confirmPassword": "abc#123@TAG",
	"password": "abc#123@aTAG"
}
`,
			wantBody: "两次输入的密码不一致",
			wantCode: 200,
		},
		{
			name: "密码必须大于8位，包含数字、特殊字符",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				return svc
			},
			reqBody: `
{
	"email": "abc123@email.com",
	"confirmPassword": "abcG",
	"password": "abcG"
}
`,
			wantBody: "密码必须大于8位，包含数字、特殊字符",
			wantCode: 200,
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				svc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "abc123@email.com",
					Password: "abc#123@TAG",
				}).Return(service.ErrUserDuplicateEmail)
				return svc
			},
			reqBody: `
{
	"email": "abc123@email.com",
	"confirmPassword": "abc#123@TAG",
	"password": "abc#123@TAG"
}
`,
			wantBody: "邮箱冲突",
			wantCode: 200,
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				svc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "abc123@email.com",
					Password: "abc#123@TAG",
				}).Return(errors.New("any error"))
				return svc
			},
			reqBody: `
{
	"email": "abc123@email.com",
	"confirmPassword": "abc#123@TAG",
	"password": "abc#123@TAG"
}
`,
			wantBody: "系统异常",
			wantCode: 200,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBufferString(tc.reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}

func TestUserHandler_LoginJWT(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				svc.EXPECT().
					Login(gomock.Any(), "abc123@email.com", "!ABC@123").
					Return(domain.User{
						Id:    1,
						Email: "abc123@email.com",
					}, nil)
				return svc
			},
			reqBody: `
{
	"email": "abc123@email.com",
	"password": "!ABC@123"
}
`,
			wantCode: 200,
			wantBody: "登录成功",
		},
		{
			name: "json 解析不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				return svc
			},
			reqBody: `
{
	"email": "abc123@email.com",
	"password": 123
}
`,
			wantCode: 400,
		},
		{
			name: "用户名或密码不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				svc.EXPECT().
					Login(gomock.Any(), "abc123@email.com", "!ABC@123").
					Return(domain.User{}, service.ErrInvalidUserOrPassword)
				return svc
			},
			reqBody: `
{
	"email": "abc123@email.com",
	"password": "!ABC@123"
}
`,
			wantCode: 200,
			wantBody: "用户名或密码不对",
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				svc.EXPECT().
					Login(gomock.Any(), "abc123@email.com", "!ABC@123").
					Return(domain.User{}, errors.New("any error"))
				return svc
			},
			reqBody: `
{
	"email": "abc123@email.com",
	"password": "!ABC@123"
}
`,
			wantCode: 200,
			wantBody: "系统错误",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			userService := tc.mock(ctrl)
			userHandler := NewUserHandler(userService, nil)
			userHandler.RegisterRoutes(server)
			req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBufferString(tc.reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}

func TestUserHandler_LoginSMS(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBody  string
		wantCode int
		wantRes  Result
	}{
		{
			name: "验证码校验通过",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), biz, "18834506987", "567135").
					Return(true, nil)
				userSvc.EXPECT().FindOrCreate(gomock.Any(), "18834506987").
					Return(domain.User{
						Id: 1,
					}, nil)
				return userSvc, codeSvc
			},
			reqBody: `
{
	"phone": "18834506987",
	"code": "567135"
}
`,
			wantCode: 200,
			wantRes: Result{
				Msg: "验证码校验通过",
			},
		},
		{
			name: "json bind 错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBody: `
{
	"phone": "18834506987",
	"code": 567135
}
`,
			wantCode: 400,
		},

		{
			name: "验证码校验系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), biz, "18834506987", "567135").
					Return(true, errors.New("any"))
				return userSvc, codeSvc
			},
			reqBody: `
{
	"phone": "18834506987",
	"code": "567135"
}
`,
			wantCode: 200,
			wantRes: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},

		{
			name: "验证码有误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), biz, "18834506987", "567135").
					Return(false, nil)
				return userSvc, codeSvc
			},
			reqBody: `
{
	"phone": "18834506987",
	"code": "567135"
}
`,
			wantCode: 200,
			wantRes: Result{
				Code: 4,
				Msg:  "验证码有误",
			},
		},

		{
			name: "查询或创建用户有误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), biz, "18834506987", "567135").
					Return(true, nil)
				userSvc.EXPECT().FindOrCreate(gomock.Any(), "18834506987").
					Return(domain.User{}, errors.New("any error"))
				return userSvc, codeSvc
			},
			reqBody: `
{
	"phone": "18834506987",
	"code": "567135"
}
`,
			wantCode: 200,
			wantRes: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userSvc, codeSvc := tc.mock(ctrl)
			userHandler := NewUserHandler(userSvc, codeSvc)
			server := gin.Default()
			userHandler.RegisterRoutes(server)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms", bytes.NewBufferString(tc.reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			if len(resp.Body.Bytes()) == 0 {
				return
			}
			var res Result
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
