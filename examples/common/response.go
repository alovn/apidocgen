package common

type Response struct {
	Code int         `json:"code"`           //返回状态码
	Msg  string      `json:"msg"`            //返回消息
	Data interface{} `json:"data,omitempty"` //返回具体数据
} //通用返回结果

func NewResponse(code int, msg string, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
