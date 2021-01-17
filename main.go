package main

import (
	_ "go-generate-code/routers"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	/*
		开启swagger
	*/
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	/*
		开启静态目录
	*/
	beego.BConfig.WebConfig.StaticDir["/static"] = "static"

	/*
		开启数据库配置
	*/
	if err := orm.RegisterDataBase("sqliteDB", "sqlite3", "./sqlite3/data.db"); err != nil {
		logs.Error("[main] sqlite3 db init err is :", err.Error())
	}

	/*
		在开发环境下，使用以下指令来开启查询调试模式
	*/
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}

	if err := logs.SetLogger(logs.AdapterFile, `{"filename":"./logs/runtime/runtime.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true,"perm":"755"}`); err != nil {
	}

	beego.Run()
}
