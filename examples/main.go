package main

import "fmt"

//@title svc-greeter
//@desc greeter接口文档
func main() {

}

//@title 测试接口
//@api GET /greeter
func greet() {
	var msg = "Hello World!"
	fmt.Println(msg)
}

type Response struct {
	Code int
	Msg  string
	Data interface{}
}
type TestData struct {
	Name string
}

//@title 测试接口2
//@api GET /greeter2
//@success Response{data=TestData}
func greet2() {
	var msg = "Hello World!"
	fmt.Println(msg)
}
