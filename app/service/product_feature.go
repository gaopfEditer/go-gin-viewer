package service

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/licensetypefeatures"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/productfeature"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/productmanager"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProductFeatureService struct{}

func NewProductFeatureService() *ProductFeatureService {
	return &ProductFeatureService{}
}

// ListProductFeatures 获取产品功能列表
func (s *ProductFeatureService) ListProductFeatures(c *gin.Context, userID, productID int, page, pageSize int) (*dto.PageResult, resource.RspCode) {
	// 检查用户是否有权限查看该产品的功能
	_, err := dto.Client().ProductManager.Query().
		Where(
			productmanager.ProductIDEQ(productID),
			productmanager.UserIDEQ(userID),
		).Only(c)

	// 超级管理员跳过权限检查
	if err != nil && userID != 1 {
		logger.Error("check permission failed", zap.Error(err))
		return nil, resource.ERR_NO_PERMISSION
	}

	// 构建查询
	q := dto.Client().ProductFeature.Query().
		Where(productfeature.ProductIDEQ(productID))

	// 计算总数
	total, err := q.Count(c)
	if err != nil {
		logger.Error("count product features failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 初始化分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // 默认每页10条
	}

	// 执行分页查询
	offset := (page - 1) * pageSize
	features, err := q.
		Limit(pageSize).
		Offset(offset).
		All(c)

	if err != nil {
		logger.Error("query product features failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 构建分页结果
	result := &dto.PageResult{
		Total:    int64(total),
		Page:     page,
		PageSize: pageSize,
		List:     features,
	}

	return result, resource.CODE_SUCCESS
}

// AddProductFeature 添加产品功能
func (s *ProductFeatureService) AddProductFeature(c *gin.Context, userID int, param dto.AddProductFeature) resource.RspCode {
	// 1. 检查用户权限
	pm, err := dto.Client().ProductManager.Query().
		Where(
			productmanager.ProductIDEQ(param.ProductID),
			productmanager.UserIDEQ(userID),
		).Only(c)
	if err != nil || (userID != 1 && pm.Permissions == productmanager.PermissionsRead) {
		return resource.ERR_NO_PERMISSION
	}

	// 2.1 检查功能编码是否已存在
	exist, err := dto.Client().ProductFeature.Query().
		Where(
			productfeature.ProductIDEQ(param.ProductID),
			productfeature.FeatureCodeEQ(param.FeatureCode),
		).Exist(c)
	if err != nil {
		logger.Error("check feature code failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if exist {
		return resource.ERR_FEATURE_CODE_EXIST
	}

	// 2.2 检查功能名称是否已存在
	exist, err = dto.Client().ProductFeature.Query().
		Where(
			productfeature.ProductIDEQ(param.ProductID),
			productfeature.FeatureNameEQ(param.FeatureName),
		).Exist(c)
	if err != nil {
		logger.Error("check feature name failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if exist {
		return resource.ERR_FEATURE_NAME_EXIST
	}
	tx, _ := dto.Client().Tx(c)
	// 3. 创建产品功能
	_, err = tx.ProductFeature.Create().
		SetProductID(param.ProductID).
		SetFeatureName(param.FeatureName).
		SetFeatureCode(param.FeatureCode).
		Save(c)
	if err != nil {
		logger.Error("create product feature failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_ADD_FAILED
	}

	// 4. 创建审计日志
	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:     userID,
		Action:     dto.ActionCreate,
		Module:     dto.ModuleFeature,
		ProductID:  param.ProductID,
		DetailInfo: param,
	})
	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_ADD_FAILED
	}

	// 5. 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("failed to commit transaction", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}
	return resource.CODE_SUCCESS
}

// DeleteProductFeature 删除产品功能
func (s *ProductFeatureService) DeleteProductFeature(c *gin.Context, userID, featureID int) resource.RspCode {
	// 1. 获取功能信息
	feature, err := dto.Client().ProductFeature.Get(c, featureID)
	if err != nil {
		logger.Error("get product feature failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	// 2. 检查用户权限
	pm, err := dto.Client().ProductManager.Query().
		Where(
			productmanager.ProductIDEQ(feature.ProductID),
			productmanager.UserIDEQ(userID),
		).Only(c)
	if err != nil || (userID != 1 && pm.Permissions == productmanager.PermissionsRead) {
		return resource.ERR_NO_PERMISSION
	}

	// 3. Start transaction
	tx, err := dto.Client().Tx(c)
	if err != nil {
		logger.Error("failed to start transaction", zap.Error(err))
		return resource.ERR_DEL_FAILED
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	// 4. Clear license type features if exists
	_, err = tx.LicenseTypeFeatures.Delete().Where(licensetypefeatures.FeatureIDEQ(featureID)).Exec(c)
	if err != nil {
		logger.Error("failed to clear license type features", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_DEL_FAILED
	}

	// 5. Delete feature
	err = tx.ProductFeature.DeleteOne(feature).Exec(c)
	if err != nil {
		logger.Error("failed to delete product feature", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_DEL_FAILED
	}

	// 6. 创建审计日志
	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:     userID,
		Action:     dto.ActionDelete,
		Module:     dto.ModuleFeature,
		ProductID:  feature.ProductID,
		DetailInfo: feature,
	})
	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_ADD_FAILED
	}

	// 6. Commit transaction
	if err := tx.Commit(); err != nil {
		logger.Error("failed to commit transaction", zap.Error(err))
		return resource.ERR_DEL_FAILED
	}

	return resource.CODE_SUCCESS
}
