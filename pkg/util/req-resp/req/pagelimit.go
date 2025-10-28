package req

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetLimitOffset(context *gin.Context) (int, int) {
	param := ParamGet(context)

	// 默认值
	const defaultLimit = 20
	const defaultPage = 1

	// 解析 Page 参数
	page, err := strconv.Atoi(param["Page"])
	if err != nil || page < 1 {
		page = defaultPage
	}

	// 解析 Limit 参数
	limit, err := strconv.Atoi(param["Limit"])
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}

	// 计算 Offset
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	return limit, offset
}
