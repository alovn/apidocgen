package common

import "encoding/xml"

type Response struct {
	XMLName xml.Name    `xml:"response"`
	Code    int         `json:"code" xml:"code"`           //返回状态码
	Msg     string      `json:"msg" xml:"msg"`             //返回消息
	Data    interface{} `json:"data,omitempty" xml:"data"` //返回具体数据
} //通用返回结果

func NewResponse(code int, msg string, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
