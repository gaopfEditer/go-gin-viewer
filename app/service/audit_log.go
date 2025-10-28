package service

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/auditlog"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"context"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type AuditLogService struct{}

func NewAuditLogService() *AuditLogService {
	return &AuditLogService{}
}

// // CreateLog 创建审计日志
// func CreateLog(ctx context.Context, operatorID int, module, actionType string, details interface{}) {
// 	// 将details转换为JSON字符串
// 	detailsJSON, err := json.Marshal(details)
// 	if err != nil {
// 		logger.Error("marshal audit log details failed", zap.Error(err))
// 		return
// 	}

// 	// 获取IP地址
// 	var ipAddress string
// 	if gc, ok := ctx.(*gin.Context); ok {
// 		ipAddress = gc.ClientIP()
// 	}

// 	// 创建审计日志
// 	_, err = dto.Client().AuditLog.Create().
// 		SetOperatorID(operatorID).
// 		SetModule(module).
// 		SetActionType(actionType).
// 		SetDetails(string(detailsJSON)).
// 		SetIPAddress(ipAddress).
// 		Save(ctx)

//		if err != nil {
//			logger.Error("create audit log failed", zap.Error(err))
//		}
//	}
func CreateAuditLog(ctx context.Context, tx *ent.Tx, data dto.AuditLogData) error {
	details, err := jsoniter.Marshal(data.DetailInfo)
	if err != nil {
		logger.Error("marshal audit log details failed", zap.Error(err))
		return err
	}
	// 获取IP地址
	var ipAddress string
	if gc, ok := ctx.(*gin.Context); ok {
		ipAddress = gc.ClientIP()
	}

	var builder *ent.AuditLogCreate
	// 判断是否传入了事务对象
	if tx != nil {
		builder = tx.AuditLog.Create() // 使用事务对象的 client
	} else {
		builder = dto.Client().AuditLog.Create() // 使用默认 client
	}

	builder = builder.
		SetOperatorID(data.UserID).
		SetModule(string(data.Module)).
		SetActionType(string(data.Action)).
		SetDetails(string(details)).
		SetIPAddress(ipAddress)

	if data.ProductID > 0 {
		builder.SetProductID(data.ProductID)
	}

	var saveErr error
	if tx != nil {
		_, saveErr = builder.Save(ctx) // 使用事务执行 Save
	} else {
		_, saveErr = builder.Save(ctx) // 使用默认 client 执行 Save
	}

	if saveErr != nil {
		logger.Error("create audit log failed", zap.Error(saveErr))
		return saveErr
	}
	return nil
}

//func CreateAuditLog(ctx context.Context, data dto.AuditLogData) error {
//	details, err := jsoniter.Marshal(data.DetailInfo)
//	if err != nil {
//		logger.Error("marshal audit log details failed", zap.Error(err))
//		return err
//	}
//	// 获取IP地址
//	var ipAddress string
//	if gc, ok := ctx.(*gin.Context); ok {
//		ipAddress = gc.ClientIP()
//	}
//
//	builder := dto.Client().AuditLog.Create().
//		SetOperatorID(data.UserID).
//		SetModule(string(data.Module)).
//		SetActionType(string(data.Action)).
//		SetDetails(string(details)).
//		SetIPAddress(ipAddress)
//
//	if data.ProductID > 0 {
//		builder.SetProductID(data.ProductID)
//	}
//
//	_, err = builder.Save(ctx)
//	if err != nil {
//		logger.Error("create audit log failed", zap.Error(err))
//		return err
//	}
//	return nil
//}

// ListLogs 查询审计日志列表
func (s *AuditLogService) ListLogs(c *gin.Context, query dto.OperationLogQuery) (*dto.PageResult, resource.RspCode) {
	// 构建查询
	q := dto.Client().AuditLog.Query().
		WithOperator().
		WithProduct()

	// 添加时间范围过滤
	if !query.StartTime.IsZero() {
		q = q.Where(auditlog.CreatedAtGTE(query.StartTime))
	}
	if !query.EndTime.IsZero() {
		q = q.Where(auditlog.CreatedAtLTE(query.EndTime))
	}

	// 添加模块过滤
	if query.Module != "" {
		q = q.Where(auditlog.ModuleEQ(query.Module))
	}

	// 添加操作类型过滤
	if query.Operation != "" {
		q = q.Where(auditlog.ActionTypeEQ(query.Operation))
	}

	// 添加用户ID过滤
	if query.UserID > 0 {
		q = q.Where(auditlog.OperatorIDEQ(query.UserID))
	}

	// 添加产品ID过滤
	if query.ProductID > 0 {
		q = q.Where(auditlog.ProductIDEQ(query.ProductID))
	}

	// 计算总数
	total, err := q.Count(c)
	if err != nil {
		logger.Error("count audit logs failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 分页参数处理
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}

	// 查询数据
	logs, err := q.
		Order(ent.Desc(auditlog.FieldCreatedAt)).
		Limit(query.PageSize).
		Offset((query.Page - 1) * query.PageSize).
		All(c)

	if err != nil {
		logger.Error("query audit logs failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 构造返回数据
	var responses []dto.OperationLogResponse
	for _, log := range logs {
		response := dto.OperationLogResponse{
			ID:        log.ID,
			UserID:    log.OperatorID,
			Module:    log.Module,
			Operation: log.ActionType,
			Detail:    log.Details,
			IPAddress: log.IPAddress,
			CreatedAt: log.CreatedAt,
		}

		// 获取操作者邮箱
		if operator := log.Edges.Operator; operator != nil {
			response.UserEmail = operator.Email
		}

		// 获取产品信息
		if product := log.Edges.Product; product != nil {
			response.ProductID = product.ID
			response.ProductName = product.ProductName
		}

		responses = append(responses, response)
	}

	result := &dto.PageResult{
		Total:    int64(total),
		Page:     query.Page,
		PageSize: query.PageSize,
		List:     responses,
	}

	return result, resource.CODE_SUCCESS
}
