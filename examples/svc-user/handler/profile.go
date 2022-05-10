package handler

import "net/http"

type ProfileHandler struct {
}

func NewAProfileHandler() *ProfileHandler {
	return &ProfileHandler{}
}

type ProfileResponse struct {
	Username string            `json:"username,omitempty"`
	Gender   uint8             `json:"gender,omitempty" example:"1"`
	Extends  map[string]string `json:"extends,omitempty"` //扩展信息
}

//@api GET /profile/get
//@title 获取用户资料
//@group profile
//@response 200 common.Response{code=0,msg="success",data=ProfileResponse}
//@author alovn
func (p *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {

}
