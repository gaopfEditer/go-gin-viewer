package req

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"reflect"
	"strconv"
	"strings"
)

var (
	typeString  = reflect.TypeOf("")
	typeInt64   = reflect.TypeOf(int64(0))
	typeFloat64 = reflect.TypeOf(float64(0))
	typeBool    = reflect.TypeOf(true)
	typeSlice   = reflect.TypeOf([]interface{}{})
	typeMap     = reflect.TypeOf(map[string]interface{}{})
)

//公开的只读访问函数

func TypeString() reflect.Type {
	return typeString
}

func TypeInt64() reflect.Type {
	return typeInt64
}

func TypeFloat64() reflect.Type {
	return typeFloat64
}

func TypeBool() reflect.Type {
	return typeBool
}

func TypeSlice() reflect.Type {
	return typeSlice
}

func TypeMap() reflect.Type {
	return typeMap
}

type PageSearch struct {
	PageNum  int64  `json:"pageNum,omitempty" form:"pageNum"`   // 页码
	PageSize int64  `json:"pageSize,omitempty" form:"pageSize"` // 每页显示数量
	Keyword  string `json:"keyword,omitempty" form:"keyword"`   // 关键词
}

type Ids struct {
	Ids []int `json:"ids,omitempty" form:"ids"` // id切片
}

// GetJsonParam 从请求的JSON体中获取指定字段的值
func GetJsonParam(context *gin.Context) map[string]interface{} {
	jsonData := map[string]interface{}{}
	err := context.ShouldBindBodyWith(&jsonData, binding.JSON)
	if err != nil {
		return nil
	}
	return jsonData
}

// GetParam 从请求的JSON体中获取指定字段的值，并尝试将其转换为指定类型
func GetParam(context *gin.Context, fieldType reflect.Type, fieldName string) (interface{}, error) {
	value := GetJsonParam(context)
	if value == nil || value[fieldName] == nil {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	fieldValue := value[fieldName]
	val := reflect.ValueOf(fieldValue)

	// 使用类型断言来处理常见类型
	switch fieldType.Kind() {
	case reflect.String:
		return fmt.Sprintf("%v", fieldValue), nil
	case reflect.Int:
		if val.Kind() == reflect.Float64 { // JSON解析数字时默认为float64
			return int(val.Float()), nil
		}
	case reflect.Int64:
		if val.Kind() == reflect.Float64 {
			return int64(val.Float()), nil
		}
	case reflect.Float64:
		if val.Kind() == reflect.Float64 {
			return val.Float(), nil
		}
	case reflect.Bool:
		if val.Kind() == reflect.Bool {
			return val.Bool(), nil
		}
	case reflect.Slice:
		if val.Kind() == reflect.Slice {
			return fieldValue, nil
		}
	case reflect.Map:
		if val.Kind() == reflect.Map {
			return fieldValue, nil
		}
	}
	// 如果类型不匹配，返回错误
	return nil, fmt.Errorf("field %s cannot be converted to %s", fieldName, fieldType.Kind())
}

// GetJsonParamStr 从请求的JSON体中获取指定字段的值，并尝试将其转换为字符串
func GetJsonParamStr(context *gin.Context, fieldName string) string {
	value := GetJsonParam(context)
	if value != nil && value[fieldName] != nil {
		return fmt.Sprintf("%v", value[fieldName])
	}
	return ""
}

// GetJsonParamInt 从请求的JSON体中获取指定字段的值，并尝试将其转换为int
func GetJsonParamInt(context *gin.Context, fieldName string) (int, bool) {
	strValue := GetJsonParamStr(context, fieldName) // 首先尝试将值转换为字符串
	if strValue == "" {
		return 0, false // 如果没有值，返回 0 和 false
	}
	intValue, err := strconv.Atoi(strValue) // 尝试将字符串转换为 int64
	if err != nil {
		return 0, false // 如果转换失败，返回 0 和 false
	}
	return intValue, true // 如果成功，返回转换后的整数值和 true
}

// GetJsonParamInt64 从请求的JSON体中获取指定字段的值，并尝试将其转换为int64
func GetJsonParamInt64(context *gin.Context, fieldName string) (int64, bool) {
	strValue := GetJsonParamStr(context, fieldName) // 首先尝试将值转换为字符串
	if strValue == "" {
		return 0, false // 如果没有值，返回 0 和 false
	}
	intValue, err := strconv.ParseInt(strValue, 10, 64) // 尝试将字符串转换为 int64
	if err != nil {
		return 0, false // 如果转换失败，返回 0 和 false
	}
	return intValue, true // 如果成功，返回转换后的整数值和 true
}

//// GetJsonParamInt64 从请求的JSON体中获取指定字段的值，并返回int64
//func GetJsonParamInt64(context *gin.Context, fieldName string) (int64, bool) {
//	if value, exists := GetJsonParam(context,  fieldName); exists {
//		if floatValue, ok := value.(float64); ok {
//			return int64(floatValue), true // 直接转换以避免精度丢失
//		}
//	}
//	return 0, false
//}

func GetJsonParamInt64Convert(context *gin.Context, fieldName string) (int64, bool) {
	value := GetJsonParam(context)
	if value != nil {
		if floatValue, ok := value[fieldName].(string); ok {
			num, err := strconv.ParseInt(floatValue, 10, 64)
			if err != nil {
				return 0, false
			}
			return num, true
		}
	}
	return 0, false
}

func ParamGet(context *gin.Context) map[string]string {
	param := make(map[string]string)
	if err := context.ShouldBindQuery(&param); err != nil {
		return nil
	}
	return param
}

func ParamPost(context *gin.Context) map[string]interface{} {
	param := make(map[string]interface{})
	if err := context.ShouldBindJSON(&param); err != nil {
		return nil
	}
	return param
}

func ParamPostString(context *gin.Context) map[string]string {
	param := make(map[string]string)
	if err := context.ShouldBindJSON(&param); err != nil {
		return nil
	}
	return param
}

func CheckParamEmpty(param interface{}, check ...string) bool {
	// 辅助函数：检查接口是否为map[string]string类型
	isMapStringString := func(param interface{}) bool {
		_, ok := param.(map[string]string)
		return ok
	}

	// 辅助函数：检查接口是否为map[string]interface{}类型
	isMapStringInterface := func(param interface{}) bool {
		_, ok := param.(map[string]interface{})
		return ok
	}

	// 辅助函数：检查接口是否为string类型
	isString := func(value interface{}) bool {
		_, ok := value.(string)
		return ok
	}

	// 检查param是否为nil或check的长度
	if param == nil || (isMapStringString(param) && len(param.(map[string]string)) < len(check)) || (isMapStringInterface(param) && len(param.(map[string]interface{})) < len(check)) {
		return true
	}

	switch p := param.(type) {
	case map[string]string:
		for _, v := range check {
			if strings.TrimSpace(p[v]) == "" {
				return true
			}
		}
	case map[string]interface{}:
		for _, v := range check {
			value, ok := p[v]
			if !ok || value == nil || (isString(value) && strings.TrimSpace(value.(string)) == "") {
				return true
			}
		}
	default:
		// 如果param不是map[string]string或map[string]interface{}
		return true
	}
	return false
}
