package controllers

import (
	"go-generate-code/models"

	"github.com/beego/beego/v2/core/logs"

	beego "github.com/beego/beego/v2/server/web"
)

type StructController struct {
	beego.Controller
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

func (c *StructController) StructConfig() {
	var rsp CommonRsp
	rsp.Data = map[string]interface{}{
		"orm_types":    models.OrmStructFieldType,
		"normal_types": models.NormalStructFieldType,
	}
	c.Data["json"] = rsp
	_ = c.ServeJSON()
}