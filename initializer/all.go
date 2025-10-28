package initializer

import (
	"cambridge-hit.com/gin-base/activateserver/pkg/util/str"
	"cambridge-hit.com/gin-base/activateserver/resource"
)

// InitAll 不直接用init函数，初始化顺序不确定，采用函数初始化更稳妥
func InitAll() {
	//根据具体项目需求调整初始化顺序
	resource.ConfigInit()
	str.SnowflakeInit(resource.Conf.App.MachineID)
	loggerInit()
	cacheInit()
	dbInit()
	//ossInit("aliyun")
}
