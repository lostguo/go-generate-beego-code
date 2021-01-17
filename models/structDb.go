package models

import (
	"time"
)

/*
结构体生成类目，用于按需生成结构体，并记录结构体的位置
*/
type StructDb struct {
	Id                   int       `orm:"default(0);pk" json:"id"`
	Description          string    `orm:"size(100);description(结构体简介，例：数据库结构体)" json:"description"`
	Name                 string    `orm:"size(100);description(结构体名称，例：AbStruct)" json:"name"`
	Type                 string    `orm:"size(100);description(枚举值：orm-数据库类型结构体、normal-常规结构体)" json:"type"`
	BeforeConvertContent string    `orm:"type(text);description(结构体原始内容：用户提交来的字段内容，便于过程分析处理)" json:"before_convert_content"`
	AfterConvertContent  string    `orm:"type(text);description(结构体生成内容：用户提交来的字段内容进行处理后)" json:"after_convert_content"`
	Uid                  int64     `orm:"default(0);description(用户id 后面考虑放到平台上，让感兴趣的小伙伴都用用)" json:"uid"`
	CreateTime           time.Time `orm:"auto_now_add;type(datetime);description(创建日期)" json:"create_time"`
	UpdateTime           time.Time `orm:"auto_now;type(datetime);description(更新日期)" json:"update_time"`
}
