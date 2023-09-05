package web

import (
	"bytes"
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
