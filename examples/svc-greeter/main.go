package main

import (
	"fmt"

	"github.com/alovn/apidoc/examples/common"
)

//@title svc-greeter
//@desc greeter接口文档
//@baseurl http://bytego.dev/admin/
func main() {
	//@group greeter
	//@title greeter分组
	//@desc greeter分组说明
	group(greet)
}

func group(f func()) {

}

type MyInt int

type Response struct {
	Code int         `json:"code" example:"0"`             //返回状态码
	Msg  string      `json:"msg,omitempty" example:"返回消息"` //返回文本消息
	Data interface{} `json:"data,omitempty"`               //返回的具体数据
} //通用返回结果

type TestData2 struct {
	MyTitle2 string //标题2
	MyAge2   int
}
type Map map[string]interface{}
type Map2 map[string]TestData2
type Node struct {
	Name  string
	Nodes map[string]Node
}
type TestData struct {
	MyTitle          string     `json:"my_title,omitempty"` //标题
	Data2            *TestData2 `json:"data2,omitempty"`
	MyIntData        int
	MyFloat64        float64
	MyFloat32        float32
	MyIntArray       []int
	MyTestData2Array []TestData2
	Int              *int
	MyInt            MyInt
	MyInts           []MyInt
	Map              Map `json:"amap"`
	Map2             Map2
	Map3             map[string]TestData2
	Nodes            map[string]Node
	Map4             map[int]Node
}

type Request struct {
	ID    int    `query:"id" header:"id" required:"true" example:"12357"` //this id
	TID   int    `param:"tid" validate:"required"`
	Name  string `json:"name,omitempty"`
	Token string `header:"token"`
} //请求对象

//@title 测试greeter
//@api GET /greeter
//@group greeter
//@accept json
//@request1 Request
//@response 200 Response{data=TestData} 输出对象 dd
//@response1 200 common.Response{data=TestData} 输出对象 dd
//@response1 500 Response{code=10010,msg="异常"} 出错了
//@response1 500 int 错误
func greet() {
	var msg = "Hello World!"
	fmt.Println(msg)
	fmt.Println(&common.Response{})
}

//@title 测试greeter2
//@api GET /greeter2
//@group greeter
//@response1 200 TestData 输出对象 dd
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
