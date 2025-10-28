package validate

import (
	"reflect"
	"regexp"
)

func IsValidValue(value interface{}, rule interface{}) bool {
	ruleValue := reflect.ValueOf(rule)
	if ruleValue.Kind() == reflect.String {
		//正则
		reg, err := regexp.Compile(rule.(string))
		if err != nil {
			return false
		}
		strValue, ok := value.(string)
		if !ok {
			return false // If value is not a string, consider it invalid
		}
		return reg.MatchString(strValue)
	}

	// 枚举值
	for i := 0; i < ruleValue.Len(); i++ {
		if reflect.DeepEqual(value, ruleValue.Index(i).Interface()) {
			return true
		}
	}
	return false
}

func UserType(t string) bool {
	if t == "0" || t == "1" {
		return true
	}
	return false
}

func Email(s string) bool {
	// 定义邮箱的正则表达式
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(s)
}

//func UploadFileType(fileType string) bool {
//	fileTypes := domain.GetFileTypes()
//	for _, ft := range fileTypes {
//		if ft == fileType {
//			return true
//		}
//	}
//	return false
//}
