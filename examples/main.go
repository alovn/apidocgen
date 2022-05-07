package main

import "fmt"

//@title svc-greeter
//@desc greeter接口文档
//@baseurl http://bytego.dev/admin/
func main() {
	//@group greeter
	//@title greeter分组
	//@desc greeter分组说明
	resigter(greet)

	resigter(greet)
}

func resigter(f func()) {

}

type Response struct {
	Code int         `json:"code" example:"0"`             //返回状态码
	Msg  string      `json:"msg,omitempty" example:"返回消息"` //返回文本消息
	Data interface{} `json:"data,omitempty"`               //返回的具体数据
}
type TestData struct {
	MyTitle string `json:"my_title,omitempty"`
}

type Request struct {
	ID    int    `query:"id"`
	TID   int    `param:"tid"`
	Name  string `json:"name,omitempty"`
	Token string `header:"token"`
}

//@title 测试greeter
//@api GET /greeter
//@group greeter
//@accept json
//@request Request
//@response1 200 Response{data=TestData} 输出对象 dd
//@response 200 Response 输出对象 dd
func greet() {
	var msg = "Hello World!"
	fmt.Println(msg)
}

//@title 测试greeter2
//@api GET /greeter2
//@group greeter
//@response 200 TestData 输出对象 dd
func greet2() {
	var msg = "Hello World!"
	fmt.Println(msg)
}

//@title 测试hello
//@api GET /hello
//@group hello
func hello() {
	var msg = "Hello World!"
	fmt.Println(msg)
}

//@title 测试hello2
//@api GET /hello2
//@group hello
func hello2() {
	var msg = "Hello World!"
	fmt.Println(msg)
}

//@title 测试other
//@api GET /other
func other() {
	var msg = "Hello World!"
	fmt.Println(msg)
}
