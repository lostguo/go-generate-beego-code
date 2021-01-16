package controllers

import beego "github.com/beego/beego/v2/server/web"

const (
	SuccessWebCode = 200
	ErrorWebCode   = 500
)

type CommonRsp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type BaseController struct {
	beego.Controller
}
