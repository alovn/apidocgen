package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alovn/apidoc/examples/common"
)

type AccountHandler struct {
}

func NewAccountHandler() *AccountHandler {
	return &AccountHandler{}
}

type RegisterRequest struct {
	Username string `form:"username" validate:"required"` //用户名
	Password string `form:"password" validate:"required"` //密码
} //注册请求参数

type RegisterResponse struct {
	Username   string `json:"username,omitempty"`    //注册的用户名
	UserID     int64  `json:"user_id,omitempty"`     //注册的用户ID
	WelcomeMsg string `json:"welcome_msg,omitempty"` //注册后的欢迎语
} //注册返回数据

//@api POST /account/register
//@title 用户注册接口
//@group account
//@request RegisterRequest
//@response 200 common.Response{code=0,msg="success",data=RegisterResponse} "注册成功返回数据"
//@response 200 common.Response{code=10011,msg="password format error"} "密码格式错误"
//@author alovn
func (h *AccountHandler) Register(w http.ResponseWriter, r *http.Request) {
	res := common.NewResponse(200, "注册成功", &RegisterResponse{
		Username:   "abc",
		UserID:     123,
		WelcomeMsg: "welcome",
	})
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)
}

type LoginRequest struct {
	Username     string `form:"username" validate:"required"` //登录用户名
	Password     string `form:"password" validate:"required"` //登录密码
	ValidateCode string `form:"validate_code"`                //验证码
} //登录请求参数

type LoginResponse struct {
	WelcomeMsg string `json:"welcome_msg,omitempty"` //登录成功欢迎语
} //登录返回数据

//@title 用户登录接口
//@api POST /account/login
//@group account
//@request LoginRequest
//@response 200 common.Response{code=0,msg="success",data=LoginResponse} "登录成功返回数据"
//@response 200 common.Response{code=10020,msg="password_error"} "密码错误"
//@author alovn
func (h *AccountHandler) Login(w http.ResponseWriter, r *http.Request) {
	//bind LoginRequest
	res := common.NewResponse(0, "success", &LoginResponse{
		WelcomeMsg: "welcome",
	})
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)
}
