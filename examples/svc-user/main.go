package main

import (
	"net/http"

	"github.com/alovn/apidoc/examples/svc-user/handler"
)

//@title 用户服务
//@service svc-user
//@desc 用户相关的服务接口
//@baseurl /user
func main() {
	mux := http.DefaultServeMux

	//@group account
	//@title 账户相关
	//@desc 账户相关的接口，含用户注册、登录等
	{
		account := handler.NewAccountHandler()
		mux.HandleFunc("/user/account/register", account.Register)
		mux.HandleFunc("/user/account/login", account.Login)
	}

	//@group address
	//@title 地址管理
	//@desc 收货地址管理接口
	{
		address := handler.NewAddressHandler()
		mux.HandleFunc("/user/address/create", address.Create)
		mux.HandleFunc("/user/address/update", address.Update)
		mux.HandleFunc("/user/address/delete", address.Delete)
		mux.HandleFunc("/user/address/list", address.List)
	}

	//@group profile
	//@title 资料管理
	//@desc 用户资料管理接口
	{
		profile := handler.NewProfileHandler()
		mux.HandleFunc("/user/profile/get", profile.Get)
	}

	//@group menu
	//@title 菜单管理
	//@desc 菜单管理接口
	{
		menu := handler.NewMenuHandler()
		mux.HandleFunc("/user/menu/nodes", menu.Nodes)
	}

	_ = http.ListenAndServe(":8000", mux)
}
