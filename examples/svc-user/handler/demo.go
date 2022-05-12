package handler

import "net/http"

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
func (h *DemoHandler) StructArray(w http.ResponseWriter, r *http.Request) {

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
//@title DemoMap
//@group demo
//@response 200 DemoMap "demo map"
func (h *DemoHandler) Map(w http.ResponseWriter, r *http.Request) {

}
