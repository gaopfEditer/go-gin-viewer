package str

import "strings"

func Unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// 去除文件名的后缀
func RemoveSuffix(filename string) string {
	if lastDotIndex := strings.LastIndex(filename, "."); lastDotIndex != -1 {
		return filename[:lastDotIndex]
	}
	return filename
}
