package service

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/licensetype"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/productfeature"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/productmanager"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LicenseTypeService struct{}

func NewLicenseTypeService() *LicenseTypeService {
	return &LicenseTypeService{}
}

// ListLicenseTypes 获取许可证类型列表
func (s *LicenseTypeService) ListLicenseTypes(c *gin.Context, userID, productID int, page, pageSize int) (*dto.PageResult, resource.RspCode) {
	// 检查用户是否有权限查看该产品的许可证类型
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
	q := dto.Client().LicenseType.Query().
		Where(licensetype.ProductIDEQ(productID)).
		WithFeatures()

	// 计算总数
	total, err := q.Count(c)
	if err != nil {
		logger.Error("count license types failed", zap.Error(err))
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
	licenseTypes, err := q.
		Limit(pageSize).
		Offset(offset).
		All(c)

	if err != nil {
		logger.Error("query license types failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 构建分页结果
	result := &dto.PageResult{
		Total:    int64(total),
		Page:     page,
		PageSize: pageSize,
		List:     licenseTypes,
	}

	return result, resource.CODE_SUCCESS
}

// AddLicenseType 添加许可证类型
func (s *LicenseTypeService) AddLicenseType(c *gin.Context, userID int, param dto.AddLicenseType) resource.RspCode {
	// 1. 检查用户权限
	pm, err := dto.Client().ProductManager.Query().
		Where(
			productmanager.ProductIDEQ(param.ProductID),
			productmanager.UserIDEQ(userID),
		).Only(c)
	if err != nil || (userID != 1 && pm.Permissions == productmanager.PermissionsRead) {
		return resource.ERR_NO_PERMISSION
	}

	// 2.1. 检查类型名称是否已存在
	exist, err := dto.Client().LicenseType.Query().
		Where(
			licensetype.ProductIDEQ(param.ProductID),
			licensetype.TypeNameEQ(param.TypeName),
		).Exist(c)
	if err != nil {
		logger.Error("check license type name failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if exist {
		return resource.ERR_LICENSE_TYPE_EXIST
	}

	// 2.2. 检查许可证类型代码是否已存在
	exist, err = dto.Client().LicenseType.Query().
		Where(
			licensetype.ProductIDEQ(param.ProductID),
			licensetype.LicenseTypeEQ(param.LicenseType),
		).Exist(c)
	if err != nil {
		logger.Error("check license type code failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if exist {
		return resource.ERR_LICENSE_CODE_EXIST
	}

	// 3. 开始事务
	client := dto.Client()
	tx, err := client.Tx(c.Request.Context())
	if err != nil {
		logger.Error("start transaction failed", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}

	// 定义延迟函数处理事务回滚
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	// 创建许可证类型
	lt, err := tx.LicenseType.Create().
		SetProductID(param.ProductID).
		SetTypeName(param.TypeName).
		SetLicenseType(param.LicenseType).
		Save(c)
	if err != nil {
		logger.Error("create license type failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_ADD_FAILED
	}

	// 添加功能关联
	if len(param.FeatureIDs) > 0 {
		features, err := tx.ProductFeature.Query().
			Where(
				productfeature.IDIn(param.FeatureIDs...),
				productfeature.ProductIDEQ(param.ProductID),
			).All(c)
		if err != nil {
			logger.Error("query features failed", zap.Error(err))
			_ = tx.Rollback()
			return resource.ERR_QUERY_FAILED
		}

		err = lt.Update().AddFeatures(features...).Exec(c)
		if err != nil {
			logger.Error("add features failed", zap.Error(err))
			_ = tx.Rollback()
			return resource.ERR_ADD_FAILED
		}
	}

	// 查询许可证类型详情
	license, err := tx.LicenseType.Query().WithLicenseTypeFeatures().Where(licensetype.IDEQ(lt.ID)).Only(c)
	if err != nil {
		logger.Error("get license type failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_QUERY_FAILED
	}

	// 创建审计日志
	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:     userID,
		Action:     dto.ActionCreate,
		Module:     dto.ModuleLicenseType,
		ProductID:  param.ProductID,
		DetailInfo: license,
	})
	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_ADD_FAILED
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("commit transaction failed", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}

	return resource.CODE_SUCCESS
}

// DeleteLicenseType 删除许可证类型
func (s *LicenseTypeService) DeleteLicenseType(c *gin.Context, userID, typeID int) resource.RspCode {
	// 1. 获取许可证类型信息
	lt, err := dto.Client().LicenseType.Get(c, typeID)
	if err != nil {
		logger.Error("get license type failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	// 2. 检查用户权限
	pm, err := dto.Client().ProductManager.Query().
		Where(
			productmanager.ProductIDEQ(lt.ProductID),
			productmanager.UserIDEQ(userID),
		).Only(c)
	if err != nil || (userID != 1 && pm.Permissions == productmanager.PermissionsRead) {
		return resource.ERR_NO_PERMISSION
	}

	// 3. 开始事务
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

	// 4. 先清除许可证类型与功能的关联关系
	_, err = tx.LicenseType.UpdateOneID(typeID).
		ClearFeatures().
		Save(c)
	if err != nil {
		logger.Error("failed to clear license type features", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_DEL_FAILED
	}

	// 5. 删除许可证类型
	err = tx.LicenseType.DeleteOne(lt).Exec(c)
	if err != nil {
		logger.Error("failed to delete license type", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_DEL_FAILED
	}

	// 创建审计日志
	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:     userID,
		Action:     dto.ActionDelete,
		Module:     dto.ModuleLicenseType,
		ProductID:  lt.ProductID,
		DetailInfo: lt,
	})

	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_ADD_FAILED
	}

	// 6. 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("failed to commit transaction", zap.Error(err))
		return resource.ERR_DEL_FAILED
	}

	return resource.CODE_SUCCESS
}

// UpdateLicenseTypeFeatures 更新许可证类型功能列表
func (s *LicenseTypeService) UpdateLicenseTypeFeatures(c *gin.Context, userID int, param dto.UpdateLicenseTypeFeatures) resource.RspCode {
	// 1. 获取许可证类型信息
	lt, err := dto.Client().LicenseType.Query().
		Where(licensetype.ID(param.TypeID)).
		Only(c)
	if err != nil {
		logger.Error("get license type failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	// 2. 检查用户权限
	pm, err := dto.Client().ProductManager.Query().
		Where(
			productmanager.ProductIDEQ(lt.ProductID),
			productmanager.UserIDEQ(userID),
		).Only(c)
	if err != nil || (userID != 1 && pm.Permissions == productmanager.PermissionsRead) {
		return resource.ERR_NO_PERMISSION
	}
	// 3. 开启事务
	tx, err := dto.Client().Tx(c)
	if err != nil {
		logger.Error("failed to start transaction", zap.Error(err))
		return resource.ERR_OPERATION_FAILED
	}

	// 4. 获取旧的功能列表
	oldLicense, err := dto.Client().LicenseType.Query().
		Where(licensetype.IDEQ(lt.ID)).
		WithFeatures().Only(c)
	if err != nil {
		logger.Error("query old features failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	// 5. 清除现有功能关联
	err = tx.LicenseType.UpdateOne(lt).ClearFeatures().Exec(c)
	if err != nil {
		logger.Error("clear features failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_MOD_FAILED
	}

	// 6. 添加新的功能关联
	var newFeatures []*ent.ProductFeature
	if len(param.FeatureIDs) > 0 {
		newFeatures, err = tx.ProductFeature.Query().
			Where(
				productfeature.IDIn(param.FeatureIDs...),
				productfeature.ProductIDEQ(lt.ProductID),
			).All(c)
		if err != nil {
			logger.Error("query features failed", zap.Error(err))
			_ = tx.Rollback()
			return resource.ERR_QUERY_FAILED
		}

		err = tx.LicenseType.UpdateOne(lt).AddFeatures(newFeatures...).Exec(c)
		if err != nil {
			logger.Error("add features failed", zap.Error(err))
			_ = tx.Rollback()
			return resource.ERR_MOD_FAILED
		}
	}

	// 7. 获取新的功能列表
	newLicense, err := dto.Client().LicenseType.Query().
		Where(licensetype.IDEQ(lt.ID)).
		WithFeatures().Only(c)
	if err != nil {
		logger.Error("get new license type failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:    userID,
		Action:    dto.ActionUpdate,
		Module:    dto.ModuleLicenseType,
		ProductID: lt.ProductID,
		DetailInfo: map[string]interface{}{
			"old_features": oldLicense.Edges.Features,
			"new_features": newLicense.Edges.Features,
		},
	})
	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_ADD_FAILED
	}

	// 7. 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("failed to commit transaction", zap.Error(err))
		return resource.ERR_MOD_FAILED
	}

	return resource.CODE_SUCCESS
}
