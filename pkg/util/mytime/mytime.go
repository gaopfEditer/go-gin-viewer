package mytime

import "time"

var AppLocation *time.Location

func init() {
	// 初始化时区（以中国时区为例，可按需修改）
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err) // 或记录日志
	}
	AppLocation = loc
}

// ParseTime 全局时间解析函数
func ParseTime(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, AppLocation)
}
