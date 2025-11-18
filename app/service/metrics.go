package service

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/metricevent"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MetricsService struct{}

func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

// ListMetrics 分页查询
func (s *MetricsService) ListMetrics(c *gin.Context, query dto.MetricEventQuery) (*dto.PageResult, resource.RspCode) {
	q := dto.Client().MetricEvent.Query()

	if query.Type != "" {
		q = q.Where(metricevent.TypeEQ(query.Type))
	}
	if query.Route != "" {
		q = q.Where(metricevent.PageEQ(query.Route))
	}
	if query.PageID != "" {
		q = q.Where(metricevent.PageIDEQ(query.PageID))
	}
	if query.StartTS > 0 {
		q = q.Where(metricevent.TsGTE(query.StartTS))
	}
	if query.EndTS > 0 {
		q = q.Where(metricevent.TsLTE(query.EndTS))
	}

	total, err := q.Count(c)
	if err != nil {
		logger.Error("count metrics failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	pageNum := query.Page
	pageSize := query.PageSize
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNum - 1) * pageSize

	list, err := q.
		Limit(pageSize).
		Offset(offset).
		Order(ent.Desc(metricevent.FieldCreatedAt)).
		All(c)
	if err != nil {
		logger.Error("query metrics failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	return &dto.PageResult{
		Total:    int64(total),
		Page:     pageNum,
		PageSize: pageSize,
		List:     list,
	}, resource.CODE_SUCCESS
}

// AddMetric 新增
func (s *MetricsService) AddMetric(c *gin.Context, param dto.AddMetricEvent) resource.RspCode {
	_, err := dto.Client().MetricEvent.Create().
		SetType(param.Type).
		SetTs(param.TS).
		SetPage(param.Page).
		SetPageID(param.PageID).
		SetReferrer(param.Referrer).
		SetUserAgent(param.UserAgent).
		SetPayload(param.Payload).
		Save(c)
	if err != nil {
		logger.Error("create metric failed", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}
	return resource.CODE_SUCCESS
}

// UpdateMetric 修改
func (s *MetricsService) UpdateMetric(c *gin.Context, param dto.ModifyMetricEvent) resource.RspCode {
	u := dto.Client().MetricEvent.UpdateOneID(param.ID)

	if param.Type != "" {
		u = u.SetType(param.Type)
	}
	if param.TS > 0 {
		u = u.SetTs(param.TS)
	}
	if param.Page != "" {
		u = u.SetPage(param.Page)
	}
	if param.PageID != "" {
		u = u.SetPageID(param.PageID)
	}
	if param.Referrer != "" {
		u = u.SetReferrer(param.Referrer)
	}
	if param.UserAgent != "" {
		u = u.SetUserAgent(param.UserAgent)
	}
	if param.Payload != "" {
		u = u.SetPayload(param.Payload)
	}

	if _, err := u.Save(c); err != nil {
		logger.Error("update metric failed", zap.Error(err))
		return resource.ERR_MOD_FAILED
	}
	return resource.CODE_SUCCESS
}

// DeleteMetric 删除
func (s *MetricsService) DeleteMetric(c *gin.Context, id int) resource.RspCode {
	if err := dto.Client().MetricEvent.DeleteOneID(id).Exec(c); err != nil {
		logger.Error("delete metric failed", zap.Error(err))
		return resource.ERR_DEL_FAILED
	}
	return resource.CODE_SUCCESS
}

// GetMetric 详情
func (s *MetricsService) GetMetric(c *gin.Context, id int) (*dto.MetricEventResponse, resource.RspCode) {
	m, err := dto.Client().MetricEvent.Get(c, id)
	if err != nil {
		logger.Error("get metric failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}
	resp := &dto.MetricEventResponse{
		ID:        m.ID,
		Type:      m.Type,
		TS:        m.Ts,
		Page:      m.Page,
		PageID:    m.PageID,
		Referrer:  m.Referrer,
		UserAgent: m.UserAgent,
		Payload:   m.Payload,
	}
	return resp, resource.CODE_SUCCESS
}
