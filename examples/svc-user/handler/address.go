package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alovn/apidoc/examples/common"
)

type AddressHandler struct {
}

func NewAddressHandler() *AddressHandler {
	return &AddressHandler{}
}

type CreateAddressRequest struct {
	CityID  int64  `form:"city_id" validate:"required" json:"city_id,omitempty"` //城市ID
	Address string `form:"address" validate:"required" json:"address,omitempty"` //地址
} //添加地址请求参数

//@api POST /address/create
//@title 添加地址接口
//@group address
//@request CreateAddressRequest
//@response 200 common.Response{code=0,msg="success"}
//@author alovn
func (h *AddressHandler) Create(w http.ResponseWriter, r *http.Request) {
	res := common.NewResponse(200, "添加成功", nil)
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)
}

type UpdateAddressRequest struct {
	ID      int64  `form:"id" validate:"required"`      //地址ID
	Address string `form:"address" validate:"required"` //地址
} //添加地址请求参数

//@api POST /address/update
//@title 更新地址接口
//@group address
//@request UpdateAddressRequest
//@response 200 common.Response{code=0,msg="success"}
//@author alovn
func (h *AddressHandler) Update(w http.ResponseWriter, r *http.Request) {
	res := common.NewResponse(200, "更新成功", nil)
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)
}

type DeleteAddressRequest struct {
	ID int64 `form:"id" validate:"required"` //地址ID
} //添加地址请求参数

//@api POST /address/delete
//@title 删除地址接口
//@group address
//@request DeleteAddressRequest
//@response 200 common.Response{code=0,msg="success"}
//@author alovn
func (h *AddressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	res := common.NewResponse(200, "删除成功", nil)
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)
}

type GetAddressRequest struct {
	ID int64 `param:"id"` //地址ID
} //获取地址请求参数

type AddressResponse struct {
	ID      int64  `json:"id,omitempty"`      //地址ID
	CityID  int64  `json:"city_id,omitempty"` //城市ID
	Address string `json:"address,omitempty"` //地址信息
} //返回地址信息

//@api GET /address/get/:id
//@title 获取地址信息
//@group address
//@request GetAddressRequest
//@response 200 common.Response{code=0,msg="success",data=AddressResponse}
//@author alovn
func (h *AddressHandler) Get(w http.ResponseWriter, r *http.Request) {
	address := AddressResponse{}
	res := common.NewResponse(200, "获取成功", address)
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)
}

//@api GET /address/list
//@title 获取地址列表
//@group address
//@response 200 common.Response{code=0,msg="success",data=[]AddressResponse}
//@author alovn
//@desc 获取收货地址列表
// Deprecated: or use @deprecated
func (h *AddressHandler) List(w http.ResponseWriter, r *http.Request) {
	addresses := []AddressResponse{}
	res := common.NewResponse(200, "删除成功", addresses)
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)
}
