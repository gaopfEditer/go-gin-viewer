package initializer

import (
	_ "cambridge-hit.com/gin-base/activateserver/pkg/util" //导入i18n等模块
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 自定义日志初始化配置
func loggerInit() {
	var err error
	env := resource.Conf.Env
	if env == "dev" {
		err = InitDev()
	} else {
		// 生产环境
		err = InitProd()
	}
	if err != nil {
		panic(err)
	}
}

// 开发日志初始化配置
func InitDev() error {
	return logger.InitLoggers(
		// 日志等级
		logger.WithDebugLevel(),
		// 时间格式化
		logger.WithTimeLayout(time.DateTime),
	)
}

// 生产日志初始化配置
func InitProd() error {
	// 获取当前可执行文件的完整路径
	execPath, _ := os.Executable()

	// 提取可执行文件的文件名
	fileName := filepath.Base(execPath)

	// 将文件名拼接到指定的目录前缀
	file := filepath.Join("log", "web", strings.TrimSuffix(filepath.Base(execPath), filepath.Ext(fileName))+".log")

	return logger.InitLoggers(
		// 日志等级
		logger.WithInfoLevel(),
		// 写出的文件
		logger.WithFileRotationP(file),
		// 不在控制台打印
		logger.WithDisableConsole(),
		// 时间格式化
		logger.WithTimeLayout(time.DateTime),
	)
}
