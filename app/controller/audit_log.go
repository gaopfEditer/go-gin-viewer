package controller

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"cambridge-hit.com/gin-base/activateserver/app/service"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/auth"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/req-resp/resp"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
)

type AuditLogController struct {
	s *service.AuditLogService
}

func NewAuditLogController() *AuditLogController {
	return &AuditLogController{s: service.NewAuditLogService()}
}

// ListLogs
// @Tags     audit
// @Summary  获取操作日志列表
// @Produce  application/json
// @Param    page     query    int     false  "页码，从1开始"   default(1)
// @Param    page_size query    int     false  "每页数量"        default(10)
// @Param    module    query    string  false  "模块名称"
// @Param    operation query    string  false  "操作类型"
// @Param    user_id   query    int     false  "操作用户ID"
// @Param    start_time query   string  false  "开始时间"
// @Param    end_time   query   string  false  "结束时间"
// @Param    Authorization  header    string  true  "Authorization"
// @Success  200      {object}  resp.Response  "获取操作日志列表"
// @Router   /activate/audit/list [get]
func (cl *AuditLogController) ListLogs(c *gin.Context) {
	// 权限检查
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	// 绑定查询参数
	var query dto.OperationLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	// 调用服务层方法查询日志
	result, code := cl.s.ListLogs(c, query)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, result)
}
