package str

import (
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

// PrettyJson 将 JSON 字符串格式化为易读的 JSON 格式
func PrettyJson(is *string) {
	// 去掉外层引号并取消转义符号
	unquoted, err := strconv.Unquote(*is)
	if err != nil {
		return
	}

	// 解析 JSON 并去掉转义符号
	var jsonData map[string]interface{}
	err = jsoniter.Unmarshal([]byte(unquoted), &jsonData)
	if err != nil {
		return
	}
	// 将结果转换为纯 JSON 格式
	prettyJSON, err := jsoniter.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return
	}
	*is = string(prettyJSON)
}

func CleanString(input string) string {
	result := ""
	for _, r := range input {
		if strconv.IsPrint(r) { // 只保留可打印字符
			result += string(r)
		}
	}
	return result
}
