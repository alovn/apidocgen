package handler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/alovn/apidocgen/examples/common"
)

type DemoHandler struct{}

func NewDemoHandler() *DemoHandler {
	return &DemoHandler{}
}

type DemoData struct {
	Title       string         `json:"title,omitempty"` //标题
	Description string         `json:"description,omitempty"`
	Count       int            `json:"count,omitempty"`
	IntArray    []int          `json:"int_array,omitempty"`
	FloatArray  []float64      `json:"float_array,omitempty"`
	IntPointer  *int           `json:"int_pointer,omitempty"`
	Map         map[string]int `json:"map,omitempty"`
	Object1     DemoObject     `json:"object_1,omitempty"`
	Object2     *DemoObject    `json:"object_2,omitempty"`
}

type DemoObject struct {
	Name string `json:"name,omitempty"`
}

type DemoMap map[string]DemoData

//@api GET /demo/struct_array
//@title struct数组
//@group demo
//@response 200 []DemoData "demo struct array"
//@version 1.0.2.1
func (h *DemoHandler) StructArray(w http.ResponseWriter, r *http.Request) {

}

type Struct1 struct {
	Name string
	Struct2
}
type Struct2 struct {
	Name2 string
}

//@api GET /demo/struct_nested
//@title struct嵌套
//@group demo
//@response 200 Struct1 "nested struct"
func (h *DemoHandler) StructNested(w http.ResponseWriter, r *http.Request) {

}

//@api GET /demo/int_array
//@title int数组
//@group demo
//@response 200 []int "demo int array"
func (h *DemoHandler) IntArray(w http.ResponseWriter, r *http.Request) {

}

//@api GET /demo/int
//@title int
//@group demo
//@response 200 int "demo int"
func (h *DemoHandler) Int(w http.ResponseWriter, r *http.Request) {

}

//@api GET /demo/map
//@title map
//@group demo
//@response 200 DemoMap "demo map"
func (h *DemoHandler) Map(w http.ResponseWriter, r *http.Request) {

}

type DemoXMLRequest struct {
	XMLName xml.Name `xml:"request"`
	ID      int64    `param:"id" xml:"id"` //DemoID
} //XML测试请求对象

type DemoXMLResponse struct {
	XMLName xml.Name `xml:"demo"`
	ID      int64    `xml:"id"`      //地址ID
	CityID  int64    `xml:"city_id"` //城市ID
	Address string   `xml:"address"` //地址信息
} //XML测试返回对象
type DemoXMLResponse2 struct {
	ID      int64  `xml:"id"`      //地址ID
	CityID  int64  `xml:"city_id"` //城市ID
	Address string `xml:"address"` //地址信息
} //XML测试返回对象2

//@api GET /demo/xml
//@title xml
//@group demo
//@accept xml
//@request DemoXMLRequest
//@format xml
//@response 200 common.Response{code=0,msg="success",data=DemoXMLResponse}
//@response 200 common.Response{code=0,msg="success",data=DemoXMLResponse2}
//@response 200 common.Response{code=10010,msg="sme error"}
//@author alovn
func (h *AddressHandler) XML(w http.ResponseWriter, r *http.Request) {
	address := DemoXMLRequest{}
	res := common.NewResponse(200, "获取成功", address)
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)
}

type DemoTime struct {
	// Title string    //测试
	Time1 time.Time `xml:"time_1" json:"time_1"`                               //example1
	Time2 time.Time `xml:"time_2" json:"time_2" example:"2022-05-14 15:04:05"` //example2
	Time3 MyTime    `xml:"time_3" json:"time_3"`
}

type MyTime time.Time

//@api GET /demo/time
//@title time
//@group demo
//@accept xml
//@format json
//@response 200 common.Response{code=0,msg="success",data=DemoTime}
//@author alovn
func (h *AddressHandler) Time(w http.ResponseWriter, r *http.Request) {
	address := DemoXMLRequest{}
	res := common.NewResponse(200, "获取成功", address)
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)
}

//@api GET /demo/jsonp
//@title jsonp
//@group demo
//@format jsonp
//@response 200 common.Response{code=0,msg="success"}
//@author alovn
func (h *AddressHandler) Jsonp(w http.ResponseWriter, r *http.Request) {
	res := common.NewResponse(200, "获取成功", nil)
	b, _ := json.Marshal(res)
	callback := r.URL.Query().Get("callback")
	if callback == "" {
		callback = "callback"
	}
	result := fmt.Sprintf("%s(%s)", callback, string(b))
	_, _ = w.Write([]byte(result))
}
