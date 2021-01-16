package models

import (
	"encoding/json"
	"fmt"
	"go-generate-code/core"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

const (
	StructTypeOrm    = "orm"    /* orm 类型结构体 */
	StructTypeNormal = "normal" /* normal 常规结构体 */

	StructFieldTypeInt      = "int"
	StructFieldTypeInt64    = "int64"
	StructFieldTypeString   = "string"
	StructFieldTypeText     = "text"
	StructFieldTypeBool     = "bool"
	StructFieldTypeArray    = "array"
	StructFieldTypeObject   = "object"
	StructFieldTypeChild    = "child"
	StructFieldTypeDateTime = "datetime"
)

var OrmStructFieldType []string
var NormalStructFieldType []string

/*
结构体生成类目，用于按需生成结构体，并记录结构体的位置
*/
type StructDb struct {
	Id                   int       `orm:"default(0);pk" json:"id"`
	Description          string    `orm:"size(100);description(结构体简介，例：数据库结构体)" json:"description"`
	Name                 string    `orm:"size(100);description(结构体名称，例：AbStruct)" json:"name"`
	Type                 string    `orm:"size(100);description(枚举值：orm-数据库类型结构体、normal-常规结构体)" json:"type"`
	BeforeConvertContent string    `orm:"type(text);description(结构体原始内容：用户提交来的服务内容，便于过程分析处理)" json:"before_convert_content"`
	AfterConvertContent  string    `orm:"type(text);description(结构体原始内容：用户提交来的服务内容，便于过程分析处理)" json:"after_convert_content"`
	Uid                  int64     `orm:"default(0);description(用户id 后面考虑放到平台上，让感兴趣的小伙伴都用用)" json:"uid"`
	CreateTime           time.Time `orm:"auto_now_add;type(datetime);description(创建日期)" json:"create_time"`
	UpdateTime           time.Time `orm:"auto_now;type(datetime);description(更新日期)" json:"update_time"`
}

/*
前端提交结构生成要素
*/
type BeforeConvertStructItem struct {
	Name        string `json:"name" description:"字段名-遵循大驼峰"`
	Description string `json:"description" description:"字段简介"`
	MemoryId    string `json:"memory_id" description:"字段本地id,根节点为-1，其他为memory-:id"`
	ParentId    string `json:"parent_id" description:"父节点，基于memory_id"`
	Type        string `json:"type" description:"类型：int、int8、int64、string、text、array、object、bool、child，array、object、bool、child不可用于orm"`
	FieldLen    int    `json:"field_len" description:"字符串长度(1-255)"`
	InitValue   string `json:"init_value" description:"初始值"`
}

/*

 */
func convertOriginStruct(db StructDb) (afterStruct string, err error) {

	for {
		var list []BeforeConvertStructItem

		strByte := []byte(db.BeforeConvertContent)
		if err = json.Unmarshal(strByte, &list); err != nil {
			logs.Error("结构体入参Json字符串解析失败", err.Error())
			break
		}

		afterStruct += "type " + core.ToUpperCamel(db.Name) + " struct { \n"

		for _, item := range list {
			afterStruct += convertStructItemToString(item) + "\n"
		}

		afterStruct += " } \n"

		break
	}
	return afterStruct, err
}

// 结构体 - 单行生成方法
// todo 需要补充 xml、json、form 灵活组合生成方式
func convertStructItemToString(item BeforeConvertStructItem) (str string) {

	if item.Type == StructTypeOrm {

		name := item.Name
		upperCamelName := core.ToUpperCamel(item.Name)
		fieldType := item.Type

		s := ""
		if item.Type == StructFieldTypeInt || item.Type == StructFieldTypeInt64 {
			s = "default(0)"
		}
		if item.Type == StructFieldTypeString {
			s = "size(" + strconv.Itoa(item.FieldLen) + ")"
		}
		if item.Type == StructFieldTypeText {
			s = "type(text)"
			fieldType = "string"
		}
		if item.Type == StructFieldTypeDateTime {
			s = "type(datetime)"
			fieldType = "time.Time"
		}

		if item.Description != "" {
			s += ";description(" + item.Description + ")"
		}

		str = fmt.Sprintf("%s %s `orm:\"%s\" json:\"%s\"` \n", upperCamelName, fieldType, s, name)
	}

	if item.Type == StructTypeNormal {
		s := ""

		name := item.Name
		upperCamelName := core.ToUpperCamel(item.Name)
		fieldType := item.Type

		if item.Type == StructFieldTypeInt || item.Type == StructFieldTypeInt64 || item.Type == StructFieldTypeString || item.Type == StructFieldTypeBool {
			s += `json:"` + name + `"`
		}

		if item.Description != "" {
			s += " description:\"" + item.Description + "\""
		}

		str = fmt.Sprintf("%s %s `%s` \n", upperCamelName, fieldType, s)
	}

	return str
}

func init() {
	/*
		orm 结构体支持的类型
	*/
	OrmStructFieldType = []string{
		StructFieldTypeInt,
		StructFieldTypeInt64,
		StructFieldTypeString,
		StructFieldTypeText,
		StructFieldTypeDateTime,
	}

	/*
		常规结构体支持的类型
	*/
	NormalStructFieldType = []string{
		StructFieldTypeInt,
		StructFieldTypeInt64,
		StructFieldTypeString,
		StructFieldTypeBool,
		StructFieldTypeArray,
		StructFieldTypeObject,
		StructFieldTypeChild,
	}
}
