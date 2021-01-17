package controllers

import (
	"encoding/json"
	"go-generate-code/logic"

	"github.com/beego/beego/v2/core/logs"
)

type StructController struct {
	BaseController
}

func (c *StructController) Prepare() {
	c.EnableRender = true
}

func (c *StructController) List() {
	c.TplName = "bot.html"
	if err := c.Render(); err != nil {
		logs.Error("[controller] struct list render err is : ", err.Error())
	}
}

func (c *StructController) Generate() {
	var req logic.GenerateStructReq
	var rsp CommonRsp
	for true {
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
			rsp.Code = ErrorWebCode
			rsp.Msg = "struct controller Generate unmarshal json err : " + err.Error()
			logs.Warn(rsp.Msg)
			break
		}

		str, err := logic.ConvertOriginStruct(req)
		if err != nil {
			rsp.Code = ErrorWebCode
			rsp.Msg = err.Error()
			break
		}

		rsp.Data = str
		rsp.Code = SuccessWebCode
		rsp.Msg = "success"
		break
	}
	c.Data["json"] = rsp
	_ = c.ServeJSON()
}

func (c *StructController) StructConfig() {
	var rsp CommonRsp
	rsp.Code = SuccessWebCode
	rsp.Data = map[string]interface{}{
		"orm_types":    logic.OrmStructFieldType,
		"normal_types": logic.NormalStructFieldType,
	}
	c.Data["json"] = rsp
	_ = c.ServeJSON()
}
