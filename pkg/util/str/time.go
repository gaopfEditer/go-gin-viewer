package str

import (
	"cambridge-hit.com/gin-base/activateserver/pkg/util/validate"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"sync"
	"time"
	_ "time/tzdata"
)

var (
	OrderFileInstanceLayout = "20060102-15_04_05.000"

	//服务器时间转浏览器时间
	standardTimeLayout = []string{"2006-01-02T15:04:05-07:00", "2006-01-02 15:04:05 -0700 MST", "20060102150405", "2006-01-02 15:04:05 MST", time.DateTime}

	//浏览器时间转服务器时间
	browserTimeLayouts = []string{
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"20060102150405",
		time.RFC3339,
		time.RFC3339Nano,
		time.DateTime,
	}
	respTimeLayout = "2006-01-02 15:04:05"
	timeFormatPool sync.Pool
	locationCache  sync.Map
)

func init() {
	timeFormatPool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
}

type TimeConverter struct {
	location *time.Location
	fields   []string
}

// ConvertSingleFiled 尝试解析时间字符串并转换为指定时区
func ConvertSingleFiled(ctx *gin.Context, timeStr, layout string) string {
	// 获取时区
	location, err := getLocation(ctx.GetHeader("TimeZone"))
	if err != nil {
		return ""
	}

	// 常见的时间布局
	layouts := append(browserTimeLayouts, standardTimeLayout...)

	// 遍历所有时间布局并尝试解析
	var t time.Time
	for _, layout := range layouts {
		t, err = time.Parse(layout, timeStr)
		if err == nil {
			break
		}
	}

	// 如果解析失败，返回空字符串
	if err != nil {
		return ""
	}

	// 将时间转换为指定时区
	t = t.In(location)

	// 返回格式化后的时间字符串
	return t.Format(layout)
}

// ServerTimeToBrowserTimeWS 转换接口返回的时间值为当地时间
func ServerTimeToBrowserTimeWS(tz string, fields []string, param interface{}) interface{} {
	if !validate.TimeZone(tz) || len(fields) == 0 {
		return param
	}

	location, err := getLocation(tz)
	if err != nil {
		return param
	}

	converter := &TimeConverter{
		location: location,
		fields:   fields,
	}

	return converter.Convert(param)
}

// ServerTimeToBrowserTime 转换接口返回的时间值为当地时间
func ServerTimeToBrowserTime(ctx *gin.Context, fields []string, param interface{}) interface{} {
	tz := ctx.GetHeader("TimeZone")
	if !validate.TimeZone(tz) || len(fields) == 0 {
		return param
	}

	location, err := getLocation(tz)
	if err != nil {
		return param
	}

	converter := &TimeConverter{
		location: location,
		fields:   fields,
	}

	return converter.Convert(param)
}

// BrowserTimeToServerTime 转换浏览器时间为服务器时间
func BrowserTimeToServerTime(ctx *gin.Context, browserTimeStr string) (string, error) {
	tz := ctx.GetHeader("TimeZone")
	if !validate.TimeZone(tz) {
		return "", fmt.Errorf("invalid timezone")
	}

	location, err := getLocation(tz)
	if err != nil {
		return "", err
	}

	var browserTime time.Time
	var parseErr error

	for _, layout := range browserTimeLayouts {
		browserTime, parseErr = time.Parse(layout, browserTimeStr)
		if parseErr == nil {
			break
		}
	}

	if parseErr != nil {
		return "", fmt.Errorf("unable to parse browser time: %s", parseErr)
	}

	// 假设服务器使用 UTC 时间
	serverLocation, _ := time.LoadLocation("UTC")
	serverTime := browserTime.In(location).In(serverLocation)

	return serverTime.Format("2006-01-02 15:04:05"), nil
}

func getLocation(tz string) (*time.Location, error) {
	if loc, ok := locationCache.Load(tz); ok {
		return loc.(*time.Location), nil
	}

	location, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}

	locationCache.Store(tz, location)
	return location, nil
}

func (tc *TimeConverter) Convert(param interface{}) interface{} {
	switch v := param.(type) {
	case map[string]interface{}:
		return tc.convertMap(v)
	case []map[string]interface{}:
		return tc.convertSlice(v)
	default:
		return tc.convertStruct(v)
	}
}

func (tc *TimeConverter) convertMap(m map[string]interface{}) map[string]interface{} {
	for _, field := range tc.fields {
		if value, ok := m[field]; ok {
			if strValue, ok := value.(string); ok {
				m[field] = tc.convertTime(strValue)
			}
		}
	}
	return m
}

func (tc *TimeConverter) convertSlice(slice []map[string]interface{}) []map[string]interface{} {
	for i := range slice {
		slice[i] = tc.convertMap(slice[i])
	}
	return slice
}

func (tc *TimeConverter) convertStruct(s interface{}) interface{} {
	jsonData, err := jsoniter.Marshal(s)
	if err != nil {
		return s
	}

	var data map[string]interface{}
	if err := jsoniter.Unmarshal(jsonData, &data); err != nil {
		return s
	}

	return tc.convertMap(data)
}

func (tc *TimeConverter) convertTime(strValue string) string {
	t, err := parseTime(strValue)
	if err != nil {
		return strValue
	}
	t = t.In(tc.location)

	sb := timeFormatPool.Get().(*strings.Builder)
	defer timeFormatPool.Put(sb)
	sb.Reset()

	sb.WriteString(t.Format(respTimeLayout))
	return sb.String()
}

func parseTime(strValue string) (time.Time, error) {
	for _, layout := range standardTimeLayout {
		if t, err := time.Parse(layout, strValue); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse time: %s", strValue)
}
