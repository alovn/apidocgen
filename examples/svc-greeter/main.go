package main

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/alovn/apidoc/examples/common"
)

//@title greeter服务
//@service svc-greeter
//@desc greeter接口文档
//@baseurl /
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
	Code int         `json:"code"`           //返回状态码
	Msg  string      `json:"msg"`            //返回消息
	Data interface{} `json:"data,omitempty"` //返回具体数据
} //通用返回结果

type TestData2 struct {
	MyTitle2 string //标题2
	MyAge2   int
}
type Map map[string]interface{}
type Map2 map[string]TestData2
type Map3 map[string]string
type Node struct {
	Name  string
	Nodes map[string]Node
}
type TestData struct {
	MyTitle string `json:"my_title,omitempty"` //标题
	// Data2            *TestData2 `json:"data2,omitempty"`
	// MyIntData int
	// MyFloat64        float64
	// MyFloat32        float32
	// MyIntArray []int
	MyTestData2Array []TestData2 `json:"my_test_data_2_array,omitempty"`
	// Int              *int
	// MyInt            MyInt
	MyInts []MyInt `json:"my_ints,omitempty"`
	// Map    Map `json:"amap"`
	// Map map[string]string
	Map2 Map2 `json:"map_2,omitempty"`
	// Map3             map[string]TestData2
	// Nodes map[string]Node
	// Map4 map[int]Node
	Time1 time.Time `xml:"time_1,omitempty"`
}

type Request struct {
	ID    int    `query:"id" header:"id" required:"true" example:"12357"` //this id
	TID   int    `param:"tid" validate:"required"`
	Name  string `json:"name,omitempty"`
	Token string `header:"token"`
} //请求对象

type Struct1 struct {
	Name  string
	Time1 time.Time
	// Struct2
}
type Struct2 struct {
	Name2 string
}

type XMLResponse struct {
	XMLName xml.Name    `xml:"response"`
	Attr    int         `json:"attr" xml:"attr,attr"` //返回状态码
	Code    int         `json:"code" xml:"code"`      //返回状态码
	Msg     string      `json:"msg" xml:"msg"`        //返回消息
	Data    interface{} //`json:"data,omitempty" xml:"data,omitempty"` //返回具体数据
	// InnerText string      `xml:",innerxml"`
	// Arr       []string    `xml:"arr"`
} //通用返回结果

type XMLData struct {
	XMLName xml.Name `xml:"mydata"`
	Title   string
	Desc    string
}
type XMLData2 struct {
	Title string
	Desc  string
}

type DemoTime struct {
	// Title string    //测试
	// Map     map[string]string //map测试
	MyTime1 MyTime
	// Time1 time.Time `xml:"time_1" json:"time_1"`                               //example1
	// Time2 time.Time `xml:"time_2" json:"time_2" example:"2022-05-14 15:04:05"` //example2
}
type MyTime time.Time

//@title 测试greeter
//@api GET /greeter
//@group greeter
//@accept xml
//@format json
//@request1 Request
//@response1 200 TestData "输出对象"
//@response1 200 Struct1 "struct1"
//@response1 200 XMLResponse{data=[]XMLData} "输出xml"
//@response1 200 XMLResponse{data=[]XMLData2} "输出xml2"
//@response 200 common.Response{data=DemoTime} "输出对象 dd"
//@response1 500 Response{code=10010,msg="异常"} "出错了"
//@response1 500 int 错误
func greet() {
	var msg = "Hello World!"
	fmt.Println(msg)
	fmt.Println(&common.Response{})
}
