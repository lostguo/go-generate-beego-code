package logic

import (
	"fmt"
	"go-generate-code/core"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
)

var OrmStructFieldType []string
var NormalStructFieldType []string

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

type StructLogic struct {
}

type GenerateStructReq struct {
	Id          int                       `json:"id"`
	Description string                    `json:"description"`
	Name        string                    `json:"name"`
	Type        string                    `json:"type" description:"枚举值：orm-数据库类型结构体、normal-常规结构体"`
	Content     []BeforeConvertStructItem `json:"content"`
}

/*
仿rap接口文档的参数结构
*/
type BeforeConvertStructItem struct {
	Id          string `json:"id" description:"id"`
	Name        string `json:"name" description:"字段名-遵循大驼峰"`
	Description string `json:"description" description:"字段简介"`
	ParentId    string `json:"parent_id" description:"父节点，基于id，若无父节点-1"`
	Type        string `json:"type" description:"类型：int、int8、int64、string、text、array、object、bool、child，array、object、bool、child不可用于orm"`
	FieldLen    int    `json:"field_len" description:"字符串长度(1-255)"`
	InitValue   string `json:"init_value" description:"初始值"`
}

/*
用户提交过来的原始数据
*/
func ConvertOriginStruct(db GenerateStructReq) (afterStruct string, err error) {

	for {
		list := db.Content

		afterStruct += "type " + core.ToUpperCamel(db.Name) + " struct { \n"

		logs.Debug(list)

		for _, item := range list {
			if item.ParentId == "-1" {
				if item.Type == StructFieldTypeArray || item.Type == StructFieldTypeObject {
					afterStruct += convertStructItemChildToString(item, list, db.Type)
				} else {
					afterStruct += convertStructItemToString(item, db.Type) + "\n"
				}
			}
		}

		afterStruct += " } \n"

		break
	}
	return afterStruct, err
}

/*
结构体单行记录生成 - 选择器
*/
func convertStructItemToString(item BeforeConvertStructItem, structType string) (str string) {

	switch structType {
	case StructTypeOrm:
		str = convertToOrmString(item)
	case StructTypeNormal:
		str = convertToNormalString(item)
	}
	return str
}

/*
结构体单行记录生成 - orm生成
*/
func convertToOrmString(item BeforeConvertStructItem) string {
	underScoreName := item.Name
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

	return fmt.Sprintf("%s %s `orm:\"%s\" json:\"%s\"` \n", upperCamelName, fieldType, s, underScoreName)
}

/*
结构体单行记录生成 - 常规结构体生成
*/
func convertToNormalString(item BeforeConvertStructItem) string {
	s := ""

	underScoreName := item.Name
	upperCamelName := core.ToUpperCamel(item.Name)
	fieldType := item.Type

	if item.Type == StructFieldTypeInt || item.Type == StructFieldTypeInt64 || item.Type == StructFieldTypeString || item.Type == StructFieldTypeBool {
		s += `json:"` + underScoreName + `"`
	}

	if item.Description != "" {
		s += " description:\"" + item.Description + "\""
	}

	return fmt.Sprintf("%s %s `%s` \n", upperCamelName, fieldType, s)
}

/*
单条记录转化为结构体字符串，采用递归方法

例：MyInfo struct {
   Id string `json:"id"`
   Name string `json:"name"`
} `json:"my_info"`
*/
func convertStructItemChildToString(actionItem BeforeConvertStructItem, list []BeforeConvertStructItem, structType string) string {

	upperCamelName := core.ToUpperCamel(actionItem.Name)
	underScoreName := actionItem.Name
	var str string
	if actionItem.Type == StructFieldTypeObject {
		str += "%s struct { \n"
	} else {
		str += "%s []struct { \n"
	}
	for _, i := range list {
		if i.ParentId == actionItem.Id {
			if i.Type == StructFieldTypeObject || i.Type == StructFieldTypeArray {
				str += convertStructItemChildToString(i, list, structType)
			} else {
				str += convertStructItemToString(i, structType) + "\n"
			}
		}
	}
	str += "} `json:\"%s\"` \n"
	return fmt.Sprintf(str, upperCamelName, underScoreName)
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
