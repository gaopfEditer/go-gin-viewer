package controller

import (
	"strconv"

	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"cambridge-hit.com/gin-base/activateserver/app/service"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/auth"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/req-resp/resp"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
)

type MetricsController struct {
	s *service.MetricsService
}

func NewMetricsController() *MetricsController {
	return &MetricsController{s: service.NewMetricsService()}
}

// ListMetrics 获取指标事件列表
// @Tags     metrics
// @Summary  获取指标事件列表
// @Produce  application/json
// @Param    page      query int    false "页码" default(1)
// @Param    page_size query int    false "每页数量" default(10)
// @Param    type      query string false "事件类型"
// @Param    pageId    query string false "页面ID"
// @Param    route     query string false "页面"
// @Param    startTs   query int64  false "开始时间戳(ms)"
// @Param    endTs     query int64  false "结束时间戳(ms)"
// @Param    Authorization header string true "Authorization"
// @Success  200 {object} resp.Response "列表"
// @Router   /activate/metrics/list [get]
func (cl *MetricsController) ListMetrics(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	var startTS, endTS int64
	if v := c.Query("startTs"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			startTS = n
		}
	}
	if v := c.Query("endTs"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			endTS = n
		}
	}

	q := dto.MetricEventQuery{
		PageParams: dto.PageParams{Page: page, PageSize: pageSize},
		Type:       c.Query("type"),
		Route:      c.Query("route"),
		PageID:     c.Query("pageId"),
		StartTS:    startTS,
		EndTS:      endTS,
	}

	result, code := cl.s.ListMetrics(c, q)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}
	resp.Success(c, result)
}

// AddMetric 添加指标事件
// @Tags     metrics
// @Summary  添加指标事件
// @Produce  application/json
// @Param    Authorization header string true "Authorization"
// @Param    data body dto.AddMetricEvent true "参数：添加指标事件"
// @Success  200 {object} resp.Response{message=string} "添加"
// @Router   /activate/metrics/add [post]
func (cl *MetricsController) AddMetric(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	var param dto.AddMetricEvent
	if err := c.ShouldBindJSON(&param); err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}
	code := cl.s.AddMetric(c, param)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}
	resp.Success(c)
}

// UpdateMetric 更新指标事件
// @Tags     metrics
// @Summary  更新指标事件
// @Produce  application/json
// @Param    Authorization header string true "Authorization"
// @Param    data body dto.ModifyMetricEvent true "参数：更新指标事件"
// @Success  200 {object} resp.Response{message=string} "更新"
// @Router   /activate/metrics/update [post]
func (cl *MetricsController) UpdateMetric(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	var param dto.ModifyMetricEvent
	if err := c.ShouldBindJSON(&param); err != nil || param.ID == 0 {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}
	code := cl.s.UpdateMetric(c, param)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}
	resp.Success(c)
}

// DeleteMetric 删除指标事件
// @Tags     metrics
// @Summary  删除指标事件
// @Produce  application/json
// @Param    Authorization header string true "Authorization"
// @Param    id query string true "事件ID"
// @Success  200 {object} resp.Response{message=string} "删除"
// @Router   /activate/metrics/delete [delete]
func (cl *MetricsController) DeleteMetric(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id == 0 {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}
	code := cl.s.DeleteMetric(c, id)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}
	resp.Success(c)
}

// GetMetric 获取指标事件详情
// @Tags     metrics
// @Summary  获取指标事件详情
// @Produce  application/json
// @Param    Authorization header string true "Authorization"
// @Param    id query string true "事件ID"
// @Success  200 {object} resp.Response "详情"
// @Router   /activate/metrics/detail [get]
func (cl *MetricsController) GetMetric(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id == 0 {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}
	result, code := cl.s.GetMetric(c, id)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}
	resp.Success(c, result)
}
