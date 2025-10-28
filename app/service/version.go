package service

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/firmwareversion"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/product"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/productfeature"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/productmanager"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/softwareversion"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/user"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/mytime"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// VersionService 版本管理服务
type VersionService struct{}

// NewVersionService 创建版本管理服务
func NewVersionService() *VersionService {
	return &VersionService{}
}

// 检查韧件版本是否已存在
func (s *VersionService) checkFirmwareVersionExists(ctx context.Context, productID int, version string, excludeID int) (bool, error) {
	client := dto.Client()
	query := client.FirmwareVersion.Query().
		Where(
			firmwareversion.ProductID(productID),
			firmwareversion.VersionEQ(version),
		)

	if excludeID > 0 {
		query = query.Where(firmwareversion.IDNEQ(excludeID))
	}

	count, err := query.Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 检查软件版本是否已存在
func (s *VersionService) checkSoftwareVersionExists(ctx context.Context, productID int, version string, excludeID int) (bool, error) {
	client := dto.Client()
	query := client.SoftwareVersion.Query().
		Where(
			softwareversion.ProductID(productID),
			softwareversion.VersionEQ(version),
		)

	if excludeID > 0 {
		query = query.Where(softwareversion.IDNEQ(excludeID))
	}

	count, err := query.Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 检查产品权限
func (s *VersionService) checkProductPermission(ctx context.Context, userID int, productID int) (bool, error) {
	client := dto.Client()

	// 超级管理员有全部权限
	if userID == dto.SuperAdminID {
		return true, nil
	}

	// 检查是否是产品管理员
	count, err := client.ProductManager.Query().
		Where(
			productmanager.ProductIDEQ(productID),
		).
		QueryUser().
		Where(
			user.ID(userID),
		).
		Count(ctx)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ListFirmwareVersions 获取韧件版本列表
func (s *VersionService) ListFirmwareVersions(c *gin.Context, userID int, productID int, page, pageSize int) (*dto.PageResult, resource.RspCode) {
	client := dto.Client()
	ctx := c.Request.Context()

	// 检查产品是否存在
	productExists, err := client.Product.Query().Where(product.ID(productID)).Exist(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}
	if !productExists {
		return nil, resource.ERR_INVALID_PARAMETER
	}

	// 查询总数
	total, err := client.FirmwareVersion.Query().
		Where(firmwareversion.ProductID(productID)).
		Count(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 分页查询
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 获取排序参数
	orderField := c.Query("sort")
	order := c.Query("order")
	if orderField == "" {
		orderField = "release_date" // 默认按发布日期排序
	}

	orderFunc := ent.Desc
	if order == "asc" {
		orderFunc = ent.Asc
	}

	// 根据排序字段选择排序方式
	var orderByField string
	switch orderField {
	case "release_date":
		orderByField = firmwareversion.FieldReleaseDate
	case "created_at":
		orderByField = firmwareversion.FieldCreatedAt
	case "version":
		orderByField = firmwareversion.FieldVersion
	default:
		orderByField = firmwareversion.FieldReleaseDate
	}

	firmwares, err := client.FirmwareVersion.Query().
		Where(firmwareversion.ProductID(productID)).
		Order(orderFunc(orderByField)).
		Limit(pageSize).
		Offset(offset).
		WithCreator(). // 加载创建人信息
		All(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 构造返回结果
	result := make([]dto.FirmwareVersionResponse, 0, len(firmwares))
	for _, v := range firmwares {
		// 获取创建者邮箱
		creatorEmail := ""
		if v.Edges.Creator != nil {
			creatorEmail = v.Edges.Creator.Email
		}

		result = append(result, dto.FirmwareVersionResponse{
			ID:             v.ID,
			ProductID:      v.ProductID,
			Version:        v.Version,
			ReleaseDate:    v.ReleaseDate.Format("2006-01-02 15:04"),
			CreatedBy:      v.CreatedBy,
			CreatedByEmail: creatorEmail,
			CreatedAt:      v.CreatedAt.Format("2006-01-02 15:04:05"),
			Remark:         v.Remark,
		})
	}

	return &dto.PageResult{
		Total:    int64(total),
		Page:     page,
		PageSize: pageSize,
		List:     result,
	}, resource.CODE_SUCCESS
}

// AddFirmwareVersion 添加韧件版本
func (s *VersionService) AddFirmwareVersion(c *gin.Context, userID int, param dto.AddFirmwareVersion) resource.RspCode {
	client := dto.Client()
	ctx := c.Request.Context()

	// 检查产品是否存在
	productExists, err := client.Product.Query().Where(product.ID(param.ProductID)).Exist(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}
	if !productExists {
		return resource.ERR_INVALID_PARAMETER
	}

	// 检查是否有权限
	hasPermission, err := s.checkProductPermission(ctx, userID, param.ProductID)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}
	if !hasPermission {
		return resource.ERR_NO_PERMISSION
	}

	// 检查版本是否已存在
	exists, err := s.checkFirmwareVersionExists(ctx, param.ProductID, param.Version, 0)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}
	if exists {
		return resource.ERR_FIRMWARE_VERSION_EXIST
	}

	// 解析发布日期
	releaseDate, err := mytime.ParseTime("2006-01-02 15:04", param.ReleaseDate)
	if err != nil {
		logger.Error("Failed to parse release date:", zap.Error(err))
		return resource.ERR_INVALID_PARAMETER
	}

	tx, _ := dto.Client().Tx(c)

	// 创建韧件版本
	_, err = tx.FirmwareVersion.Create().
		SetProductID(param.ProductID).
		SetVersion(param.Version).
		SetReleaseDate(releaseDate).
		SetRemark(param.Remark).
		SetCreatedBy(userID).
		Save(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_ADD_FAILED
	}

	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:     userID,
		Action:     dto.ActionCreate,
		Module:     dto.ModuleFirmwareVersion,
		ProductID:  param.ProductID,
		DetailInfo: param,
	})
	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		tx.Rollback()
		return resource.ERR_ADD_LOG_FAILED
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("commit transaction failed", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}

	return resource.CODE_SUCCESS
}

func (s *VersionService) ModifyFirmwareVersion(c *gin.Context, userID int, param dto.ModifyFirmwareVersion) resource.RspCode {
	client := dto.Client()
	ctx := c.Request.Context()

	// 获取韧件版本信息
	fw, err := client.FirmwareVersion.Get(ctx, param.ID)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_FIRMWARE_NOT_EXIST
	}

	// 检查是否有权限
	hasPermission, err := s.checkProductPermission(ctx, userID, fw.ProductID)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_MOD_FAILED
	}
	if !hasPermission {
		return resource.ERR_NO_PERMISSION
	}

	tx, _ := dto.Client().Tx(c)

	update := tx.FirmwareVersion.UpdateOneID(param.ID)
	changed := false

	oldVersion := fw.Version

	// 版本号更新：只有提供了新版本号且与原版本不同时才更新
	if param.Version != "" && param.Version != oldVersion {
		// 检查新版本号是否已经存在
		exists, err := s.checkFirmwareVersionExists(ctx, fw.ProductID, param.Version, param.ID)
		if err != nil {
			logger.Error("err:", zap.Error(err))
			return resource.ERR_MOD_FAILED
		}
		if exists {
			return resource.ERR_FIRMWARE_VERSION_EXIST
		}
		update = update.SetVersion(param.Version)
		changed = true
	}

	// 发布日期更新：只有提供了新日期且与原日期不同才更新
	if param.ReleaseDate != "" {
		releaseDate, err := mytime.ParseTime("2006-01-02 15:04", param.ReleaseDate)
		if err != nil {
			logger.Error("Failed to parse release date:", zap.Error(err))
			return resource.ERR_INVALID_PARAMETER
		}
		if fw.ReleaseDate.Format("2006-01-02 15:04") != param.ReleaseDate {
			update = update.SetReleaseDate(releaseDate)
			changed = true
		}
	}

	// 备注更新：只有提供了备注且与原备注不同才更新
	if fw.Remark != param.Remark {
		update = update.SetRemark(param.Remark)
		changed = true
	}

	// 如果没有任何改动，则直接返回成功，不写日志
	if !changed {
		return resource.CODE_SUCCESS
	}

	// 执行更新
	new_fw, err := update.Save(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_MOD_FAILED
	}

	// 查询产品
	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:    userID,
		Action:    dto.ActionUpdate,
		Module:    dto.ModuleFirmwareVersion,
		ProductID: fw.ProductID,
		DetailInfo: map[string]interface{}{
			"old_version_info": fw,
			"new_version_info": new_fw,
		},
	})
	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		tx.Rollback()
		return resource.ERR_ADD_LOG_FAILED
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("commit transaction failed", zap.Error(err))
		return resource.ERR_MOD_FAILED
	}

	return resource.CODE_SUCCESS
}

// DeleteFirmwareVersion 删除韧件版本
func (s *VersionService) DeleteFirmwareVersion(c *gin.Context, userID int, firmwareID int) resource.RspCode {
	client := dto.Client()
	ctx := c.Request.Context()

	// 查询韧件版本
	firmwareVersion, err := client.FirmwareVersion.Get(ctx, firmwareID)
	if err != nil {
		if ent.IsNotFound(err) {
			return resource.ERR_FIRMWARE_NOT_EXIST
		}
		logger.Error("err:", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	// 检查权限
	hasPermission, err := s.checkProductPermission(ctx, userID, firmwareVersion.ProductID)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if !hasPermission {
		return resource.ERR_NO_PERMISSION
	}

	// 开始事务
	tx, err := client.Tx(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_DEL_FAILED
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	// 删除韧件版本
	err = tx.FirmwareVersion.DeleteOne(firmwareVersion).Exec(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_DEL_FAILED
	}

	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:    userID,
		Action:    dto.ActionDelete,
		Module:    dto.ModuleFirmwareVersion,
		ProductID: firmwareVersion.ProductID,
		DetailInfo: map[string]interface{}{
			"old_version_info": firmwareVersion,
		},
	})
	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		tx.Rollback()
		return resource.ERR_ADD_LOG_FAILED
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_DEL_FAILED
	}

	return resource.CODE_SUCCESS
}

// ListSoftwareVersions 获取软件版本列表
func (s *VersionService) ListSoftwareVersions(c *gin.Context, userID int, productID int, page, pageSize int) (*dto.PageResult, resource.RspCode) {
	client := dto.Client()
	ctx := c.Request.Context()

	// 检查产品是否存在
	productExists, err := client.Product.Query().Where(product.ID(productID)).Exist(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}
	if !productExists {
		return nil, resource.ERR_INVALID_PARAMETER
	}

	// 查询总数
	total, err := client.SoftwareVersion.Query().
		Where(softwareversion.ProductID(productID)).
		Count(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 分页查询
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	versions, err := client.SoftwareVersion.Query().
		Where(softwareversion.ProductID(productID)).
		Order(ent.Desc(softwareversion.FieldCreatedAt)).
		Limit(pageSize).
		Offset(offset).
		WithFeatures().
		WithFirmwareVersions().
		WithCreator().
		All(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 构造返回结果
	result := make([]dto.SoftwareVersionResponse, 0, len(versions))
	for _, v := range versions {
		// 处理功能列表
		featureInfos := make([]dto.FeatureInfo, 0)
		if v.Edges.Features != nil {
			for _, f := range v.Edges.Features {
				featureInfos = append(featureInfos, dto.FeatureInfo{
					ID:          f.ID,
					FeatureName: f.FeatureName,
					FeatureCode: f.FeatureCode,
				})
			}
		}

		// 处理韧件版本列表
		firmwareInfos := make([]dto.FirmwareInfo, 0)
		if v.Edges.FirmwareVersions != nil {
			for _, f := range v.Edges.FirmwareVersions {
				firmwareInfos = append(firmwareInfos, dto.FirmwareInfo{
					ID:          f.ID,
					Version:     f.Version,
					ReleaseDate: f.ReleaseDate,
				})
			}
		}

		// 添加软件版本信息
		creatorEmail := ""
		if v.Edges.Creator != nil {
			creatorEmail = v.Edges.Creator.Email
		}

		result = append(result, dto.SoftwareVersionResponse{
			ID:             v.ID,
			ProductID:      v.ProductID,
			Version:        v.Version,
			ReleaseDate:    v.ReleaseDate.Format("2006-01-02 15:04"),
			CreatedBy:      v.CreatedBy,
			CreatedByEmail: creatorEmail,
			CreatedAt:      v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdateLog:      v.UpdateLog,
			Remark:         v.Remark,
			Features:       featureInfos,
			Firmwares:      firmwareInfos,
		})
	}

	return &dto.PageResult{
		Total:    int64(total),
		Page:     page,
		PageSize: pageSize,
		List:     result,
	}, resource.CODE_SUCCESS
}

// AddSoftwareVersion 添加软件版本
func (s *VersionService) AddSoftwareVersion(c *gin.Context, userID int, param dto.AddSoftwareVersion) resource.RspCode {
	client := dto.Client()
	ctx := c.Request.Context()

	// 检查产品是否存在
	productExists, err := client.Product.Query().Where(product.ID(param.ProductID)).Exist(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if !productExists {
		return resource.ERR_INVALID_PARAMETER
	}

	// 检查权限
	hasPermission, err := s.checkProductPermission(ctx, userID, param.ProductID)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if !hasPermission {
		return resource.ERR_NO_PERMISSION
	}

	// 检查版本是否已存在
	exists, err := s.checkSoftwareVersionExists(ctx, param.ProductID, param.Version, 0)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if exists {
		return resource.ERR_SOFTWARE_VERSION_EXIST
	}
	// 解析发布日期
	releaseDate, err := mytime.ParseTime("2006-01-02 15:04", param.ReleaseDate)
	if err != nil {
		logger.Error("Failed to parse release date:", zap.Error(err))
		return resource.ERR_INVALID_PARAMETER
	}

	// 开始事务
	tx, err := client.Tx(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	// 创建软件版本
	softwareCreate := tx.SoftwareVersion.Create().
		SetProductID(param.ProductID).
		SetVersion(param.Version).
		SetReleaseDate(releaseDate).
		SetUpdateLog(param.UpdateLog).
		SetRemark(param.Remark).
		SetCreatedBy(userID)

	// 添加关联的功能
	if len(param.FeatureIDs) > 0 {
		features, err := tx.ProductFeature.Query().
			Where(
				productfeature.IDIn(param.FeatureIDs...),
				productfeature.ProductID(param.ProductID),
			).All(ctx)
		if err != nil {
			logger.Error("err:", zap.Error(err))
			_ = tx.Rollback()
			return resource.ERR_ADD_FAILED
		}
		softwareCreate.AddFeatures(features...)
	}

	// 添加关联的韧件版本
	if len(param.FirmwareIDs) > 0 {
		firmwares, err := tx.FirmwareVersion.Query().
			Where(
				firmwareversion.IDIn(param.FirmwareIDs...),
				firmwareversion.ProductID(param.ProductID),
			).All(ctx)
		if err != nil {
			logger.Error("err:", zap.Error(err))
			_ = tx.Rollback()
			return resource.ERR_ADD_FAILED
		}
		softwareCreate.AddFirmwareVersions(firmwares...)
	}

	// 保存软件版本
	sv, err := softwareCreate.Save(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_ADD_FAILED
	}

	// 查询关联的功能
	featuresList, err := sv.QueryFeatures().All(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_QUERY_FAILED
	}

	features := make(map[int]string)
	for _, f := range featuresList {
		features[f.ID] = f.FeatureName
	}
	// 查询关联的韧件版本
	firmwareList, err := sv.QueryFirmwareVersions().All(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_QUERY_FAILED
	}

	firmwares := make(map[int]string)
	for _, fw := range firmwareList {
		firmwares[fw.ID] = fw.Version
	}

	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:    userID,
		Action:    dto.ActionCreate,
		Module:    dto.ModuleSoftwareVersion,
		ProductID: param.ProductID,
		DetailInfo: map[string]interface{}{
			"new_version_info": sv,
		},
	})
	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		tx.Rollback()
		return resource.ERR_ADD_LOG_FAILED
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}

	return resource.CODE_SUCCESS
}

// ModifySoftwareVersion 修改软件版本
func (s *VersionService) ModifySoftwareVersion(c *gin.Context, userID int, param dto.ModifySoftwareVersion) resource.RspCode {
	client := dto.Client()
	ctx := c.Request.Context()

	// 查询软件版本（包含关联关系）
	softwareVersion, err := client.SoftwareVersion.Query().
		Where(softwareversion.ID(param.ID)).
		WithFeatures().
		WithFirmwareVersions().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return resource.ERR_SOFTWARE_NOT_EXIST
		}
		logger.Error("err:", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	// 检查权限
	hasPermission, err := s.checkProductPermission(ctx, userID, softwareVersion.ProductID)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if !hasPermission {
		return resource.ERR_NO_PERMISSION
	}
	// 检查版本是否已存在
	versionChanged := param.Version != "" && param.Version != softwareVersion.Version
	if versionChanged {
		exists, err := s.checkSoftwareVersionExists(ctx, softwareVersion.ProductID, param.Version, param.ID)
		if err != nil {
			logger.Error("err:", zap.Error(err))
			return resource.ERR_QUERY_FAILED
		}
		if exists {
			return resource.ERR_SOFTWARE_VERSION_EXIST
		}
	}

	// 开始事务
	tx, err := client.Tx(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_MOD_FAILED
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	// 更新软件版本基础信息
	update := tx.SoftwareVersion.UpdateOneID(param.ID)
	hasUpdate := false

	if versionChanged {
		update = update.SetVersion(param.Version)
		hasUpdate = true
	}

	if param.ReleaseDate != "" {
		releaseDate, err := mytime.ParseTime("2006-01-02 15:04", param.ReleaseDate)
		if err != nil {
			logger.Error("Failed to parse release date:", zap.Error(err))
			_ = tx.Rollback()
			return resource.ERR_INVALID_PARAMETER
		}
		if !releaseDate.Equal(softwareVersion.ReleaseDate) {
			update = update.SetReleaseDate(releaseDate)
			hasUpdate = true
		}
	}

	if param.UpdateLog != softwareVersion.UpdateLog {
		update = update.SetUpdateLog(param.UpdateLog)
		hasUpdate = true
	}

	if param.Remark != softwareVersion.Remark {
		update = update.SetRemark(param.Remark)
		hasUpdate = true
	}

	// 仅当有修改时执行基础更新
	if hasUpdate {
		_, err = update.Save(ctx)
		if err != nil {
			logger.Error("err:", zap.Error(err))
			_ = tx.Rollback()
			return resource.ERR_MOD_FAILED
		}
	}

	// 处理功能关联变更
	featuresChanged := false
	currentFeatureIDs := getFeatureIDs(softwareVersion.Edges.Features)
	if param.FeatureIDs != nil && !equalIntSlices(param.FeatureIDs, currentFeatureIDs) {
		// 清除并重建关联
		_, err = tx.SoftwareVersion.UpdateOneID(param.ID).ClearFeatures().Save(ctx)
		if err != nil {
			logger.Error("err:", zap.Error(err))
			_ = tx.Rollback()
			return resource.ERR_MOD_FAILED
		}

		if len(param.FeatureIDs) > 0 {
			features, err := tx.ProductFeature.Query().
				Where(
					productfeature.IDIn(param.FeatureIDs...),
					productfeature.ProductID(softwareVersion.ProductID),
				).All(ctx)
			if err != nil {
				logger.Error("err:", zap.Error(err))
				_ = tx.Rollback()
				return resource.ERR_MOD_FAILED
			}

			_, err = tx.SoftwareVersion.UpdateOneID(param.ID).AddFeatures(features...).Save(ctx)
			if err != nil {
				logger.Error("err:", zap.Error(err))
				_ = tx.Rollback()
				return resource.ERR_MOD_FAILED
			}
		}
		featuresChanged = true
	}

	// 处理韧件关联变更
	firmwaresChanged := false
	currentFirmwareIDs := getFirmwareIDs(softwareVersion.Edges.FirmwareVersions)
	if param.FirmwareIDs != nil && !equalIntSlices(param.FirmwareIDs, currentFirmwareIDs) {
		// 清除并重建关联
		_, err = tx.SoftwareVersion.UpdateOneID(param.ID).ClearFirmwareVersions().Save(ctx)
		if err != nil {
			logger.Error("err:", zap.Error(err))
			_ = tx.Rollback()
			return resource.ERR_MOD_FAILED
		}

		if len(param.FirmwareIDs) > 0 {
			firmwares, err := tx.FirmwareVersion.Query().
				Where(
					firmwareversion.IDIn(param.FirmwareIDs...),
					firmwareversion.ProductID(softwareVersion.ProductID),
				).All(ctx)
			if err != nil {
				logger.Error("err:", zap.Error(err))
				_ = tx.Rollback()
				return resource.ERR_MOD_FAILED
			}

			_, err = tx.SoftwareVersion.UpdateOneID(param.ID).AddFirmwareVersions(firmwares...).Save(ctx)
			if err != nil {
				logger.Error("err:", zap.Error(err))
				_ = tx.Rollback()
				return resource.ERR_MOD_FAILED
			}
		}
		firmwaresChanged = true
	}

	// 无任何修改直接返回
	if !hasUpdate && !featuresChanged && !firmwaresChanged {
		_ = tx.Rollback()
		return resource.CODE_SUCCESS
	}

	// 获取更新后的数据
	updatedVersion, err := client.SoftwareVersion.Query().
		Where(softwareversion.ID(param.ID)).
		WithFeatures().
		WithFirmwareVersions().
		Only(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:    userID,
		Action:    dto.ActionUpdate,
		Module:    dto.ModuleSoftwareVersion,
		ProductID: updatedVersion.ProductID,
		DetailInfo: map[string]interface{}{
			"old_version_info": softwareVersion,
			"new_version_info": updatedVersion,
		},
	})
	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		tx.Rollback()
		return resource.ERR_ADD_LOG_FAILED
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_MOD_FAILED
	}

	return resource.CODE_SUCCESS
}

// 辅助函数
func getFeatureIDs(features []*ent.ProductFeature) []int {
	ids := make([]int, len(features))
	for i, f := range features {
		ids[i] = f.ID
	}
	return ids
}

// 辅助函数
func getFirmwareIDs(firmwares []*ent.FirmwareVersion) []int {
	ids := make([]int, len(firmwares))
	for i, fw := range firmwares {
		ids[i] = fw.ID
	}
	return ids
}

// 辅助函数
func equalIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[int]bool, len(a))
	for _, v := range a {
		m[v] = true
	}
	for _, v := range b {
		if !m[v] {
			return false
		}
	}
	return true
}

// DeleteSoftwareVersion 删除软件版本
func (s *VersionService) DeleteSoftwareVersion(c *gin.Context, userID int, softwareID int) resource.RspCode {
	client := dto.Client()
	ctx := c.Request.Context()

	// 查询软件版本
	softwareVersion, err := client.SoftwareVersion.Get(ctx, softwareID)
	if err != nil {
		if ent.IsNotFound(err) {
			return resource.ERR_SOFTWARE_NOT_EXIST
		}
		logger.Error("err:", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	// 检查权限
	hasPermission, err := s.checkProductPermission(ctx, userID, softwareVersion.ProductID)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if !hasPermission {
		return resource.ERR_NO_PERMISSION
	}

	// 开始事务
	tx, err := client.Tx(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_DEL_FAILED
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	// 删除软件版本
	err = tx.SoftwareVersion.DeleteOne(softwareVersion).Exec(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_DEL_FAILED
	}
	// 记录审计日志
	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:    userID,
		Action:    dto.ActionDelete,
		Module:    dto.ModuleSoftwareVersion,
		ProductID: softwareVersion.ProductID,
		DetailInfo: map[string]interface{}{
			"deleted_version_info": softwareVersion,
		},
	})
	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		tx.Rollback()
		return resource.ERR_ADD_LOG_FAILED
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("err:", zap.Error(err))
		return resource.ERR_DEL_FAILED
	}

	return resource.CODE_SUCCESS
}

// GetProductFirmwareVersions 获取产品的所有韧件版本
func (s *VersionService) GetProductFirmwareVersions(c *gin.Context, productID int) ([]dto.FirmwareInfo, resource.RspCode) {
	client := dto.Client()
	ctx := c.Request.Context()

	// 检查产品是否存在
	productExists, err := client.Product.Query().Where(product.ID(productID)).Exist(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}
	if !productExists {
		return nil, resource.ERR_INVALID_PARAMETER
	}

	// 获取韧件版本
	firmwares, err := client.FirmwareVersion.Query().
		Where(firmwareversion.ProductID(productID)).
		Order(ent.Desc(firmwareversion.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 构造返回结果
	result := make([]dto.FirmwareInfo, 0, len(firmwares))
	for _, v := range firmwares {
		result = append(result, dto.FirmwareInfo{
			ID:          v.ID,
			Version:     v.Version,
			ReleaseDate: v.ReleaseDate,
		})
	}

	return result, resource.CODE_SUCCESS
}

// GetProductFeatures 获取产品的所有功能
func (s *VersionService) GetProductFeatures(c *gin.Context, productID int) ([]dto.FeatureInfo, resource.RspCode) {
	client := dto.Client()
	ctx := c.Request.Context()

	// 检查产品是否存在
	productExists, err := client.Product.Query().Where(product.ID(productID)).Exist(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}
	if !productExists {
		return nil, resource.ERR_INVALID_PARAMETER
	}

	// 获取功能列表
	features, err := client.ProductFeature.Query().
		Where(productfeature.ProductID(productID)).
		Order(ent.Asc(productfeature.FieldFeatureName)).
		All(ctx)
	if err != nil {
		logger.Error("err:", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 构造返回结果
	result := make([]dto.FeatureInfo, 0, len(features))
	for _, v := range features {
		result = append(result, dto.FeatureInfo{
			ID:          v.ID,
			FeatureName: v.FeatureName,
			FeatureCode: v.FeatureCode,
		})
	}

	return result, resource.CODE_SUCCESS
}
