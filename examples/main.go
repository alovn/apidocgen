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
	Code int
	Msg  string
	Data interface{}
}
type TestData struct {
	Name string
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
func greet() {
	var msg = "Hello World!"
	fmt.Println(msg)
}

//@title 测试greeter2
//@api GET /greeter2
//@group greeter
//@request Request
//@success Response{data=TestData}
func greet2() {
	var msg = "Hello World!"
	fmt.Println(msg)
}

//@title 测试hello
//@api GET /hello
//@group hello
//@request Request
//@success Response{data=TestData}
func hello() {
	var msg = "Hello World!"
	fmt.Println(msg)
}

//@title 测试hello2
//@api GET /hello2
//@group hello
//@request Request
//@success Response{data=TestData}
func hello2() {
	var msg = "Hello World!"
	fmt.Println(msg)
}

//@title 测试other
//@api GET /other
//@request Request
//@success Response{data=TestData}
func other() {
	var msg = "Hello World!"
	fmt.Println(msg)
}
