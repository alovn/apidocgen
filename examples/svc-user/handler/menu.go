package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alovn/apidocgen/examples/common"
)

type MenuHandler struct{}

func NewMenuHandler() *MenuHandler {
	return &MenuHandler{}
}

type Node struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Nodes []Node `json:"nodes"`
}

//@api GET /menu/nodes
//@title 获取菜单节点
//@group menu
//@response 200 common.Response{code=0,msg="success",data=[]Node}
//@author alovn
//@desc 测试数组、递归结构体
func (h *MenuHandler) Nodes(w http.ResponseWriter, r *http.Request) {
	var nodes []Node
	res := common.NewResponse(200, "success", nodes)
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)
}
