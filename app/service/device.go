package service

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"strings"
	"time"

	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/device"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/licensetype"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/product"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/productmanager"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DeviceService 设备管理服务
type DeviceService struct{}

// NewDeviceService 创建设备管理服务实例
func NewDeviceService() *DeviceService {
	return &DeviceService{}
}

// ListProducts 获取产品列表（带设备数量统计）
func (s *DeviceService) ListProducts(c *gin.Context, userID int, page, pageSize int) (*dto.PageResult, resource.RspCode) {
	// 构建查询
	q := dto.Client().ProductManager.Query()
	if userID != 1 {
		// 用户ID不是1，列出用户作为管理员或协作者管理的产品
		q = q.Where(productmanager.UserIDEQ(userID))
	}

	// 获取搜索关键字
	searchQuery := c.Query("search")

	// 构建产品查询
	prodQuery := q.QueryProduct()

	// 如果有搜索关键字，添加搜索条件
	if searchQuery != "" {
		prodQuery = prodQuery.Where(
			product.Or(
				product.CodeContainsFold(searchQuery),
				product.ProductNameContainsFold(searchQuery),
				product.ProductTypeContainsFold(searchQuery),
			),
		)
	}

	// 计算总数
	total, err := prodQuery.Count(c)
	if err != nil {
		logger.Error("count products failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 执行分页查询
	offset := (page - 1) * pageSize
	products, err := prodQuery.
		WithManagers(func(pmq *ent.ProductManagerQuery) { pmq.WithUser() }).
		Limit(pageSize).
		Offset(offset).
		All(c)

	if err != nil {
		logger.Error("query user managed products failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 计算设备数量
	var deviceSummaries []dto.DeviceSummary
	for _, p := range products {
		count, err := dto.Client().Device.Query().Where(device.ProductIDEQ(p.ID)).Count(c)
		if err != nil {
			logger.Error("count devices failed", zap.Error(err), zap.Int("product_id", p.ID))
			continue
		}

		deviceSummaries = append(deviceSummaries, dto.DeviceSummary{
			ProductID:   p.ID,
			ProductName: p.ProductName,
			Count:       count,
		})
	}

	// 构建分页结果
	result := &dto.PageResult{
		Total:    int64(total),
		Page:     page,
		PageSize: pageSize,
		List:     products,
		Extra: map[string]interface{}{
			"device_count": deviceSummaries,
		},
	}

	return result, resource.CODE_SUCCESS
}

// ListDevices 获取产品下设备列表
func (s *DeviceService) ListDevices(c *gin.Context, userID int, filter dto.DeviceFilter) (*dto.PageResult, resource.RspCode) {
	// 权限检查
	//if userID != 1 {
	//	exist, err := dto.Client().ProductManager.Query().
	//		Where(
	//			productmanager.ProductIDEQ(filter.ProductID),
	//			productmanager.UserIDEQ(userID),
	//		).Exist(c)
	//	if err != nil || !exist {
	//		return nil, resource.ERR_NO_PERMISSION
	//	}
	//}

	// 构建查询
	q := dto.Client().Device.Query()
	if filter.ProductID != 0 {
		q = q.Where(device.ProductIDEQ(filter.ProductID))
	}

	// 应用过滤条件
	if filter.LicenseTypeID > 0 {
		q = q.Where(device.LicenseTypeIDEQ(filter.LicenseTypeID))
	}

	if filter.SN != "" {
		q = q.Where(device.SnContainsFold(filter.SN))
	}

	if filter.OEMTag != "" {
		q = q.Where(device.OemTagContainsFold(filter.OEMTag))
	}

	// 计算总数
	total, err := q.Count(c)
	if err != nil {
		logger.Error("count devices failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 执行分页查询
	offset := (filter.Page - 1) * filter.PageSize
	devices, err := q.
		WithProduct().
		WithLicenseType().
		WithCreator().
		WithUpdater().
		Order(ent.Desc(device.FieldCreatedAt)).
		Limit(filter.PageSize).
		Offset(offset).
		All(c)

	if err != nil {
		logger.Error("query devices failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 转换为DTO
	deviceInfos := make([]dto.DeviceInfo, 0, len(devices))
	for _, d := range devices {
		deviceInfo := dto.DeviceInfo{
			ID:            d.ID,
			SN:            d.Sn,
			ProductID:     d.ProductID,
			LicenseTypeID: d.LicenseTypeID,
			OEMTag:        d.OemTag,
			Remark:        d.Remark,
			CreatedAt:     d.CreatedAt,
			CreatedBy:     d.CreatedBy,
			UpdatedAt:     d.UpdatedAt,
			UpdatedBy:     d.UpdatedBy,
		}

		// 添加关联信息
		if d.Edges.Product != nil {
			deviceInfo.ProductName = d.Edges.Product.ProductName
			deviceInfo.ProductCode = d.Edges.Product.Code
		}

		if d.Edges.LicenseType != nil {
			deviceInfo.LicenseTypeName = d.Edges.LicenseType.TypeName
			deviceInfo.LicenseTypeCode = d.Edges.LicenseType.LicenseType
		}

		if d.Edges.Creator != nil {
			deviceInfo.CreatedByEmail = d.Edges.Creator.Email
		}

		if d.Edges.Updater != nil {
			deviceInfo.UpdatedByEmail = d.Edges.Updater.Email
		}

		deviceInfos = append(deviceInfos, deviceInfo)
	}

	// 构建分页结果
	result := &dto.PageResult{
		Total:    int64(total),
		Page:     filter.Page,
		PageSize: filter.PageSize,
		List:     deviceInfos,
	}

	return result, resource.CODE_SUCCESS
}

// GetDeviceBySN 通过SN获取设备
func (s *DeviceService) GetDeviceBySN(c *gin.Context, userID int, sn string) (*dto.DeviceInfo, resource.RspCode) {
	d, err := dto.Client().Device.Query().
		Where(device.SnEQ(sn)).
		WithProduct().
		WithLicenseType().
		WithCreator().
		WithUpdater().
		Only(c)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, resource.ERR_DEVICE_NOT_EXIST
		}
		logger.Error("query device failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 权限检查
	if userID != 1 {
		exist, err := dto.Client().ProductManager.Query().
			Where(
				productmanager.ProductIDEQ(d.ProductID),
				productmanager.UserIDEQ(userID),
			).Exist(c)
		if err != nil || !exist {
			return nil, resource.ERR_NO_PERMISSION
		}
	}

	// 转换为DTO
	deviceInfo := dto.DeviceInfo{
		ID:            d.ID,
		SN:            d.Sn,
		ProductID:     d.ProductID,
		LicenseTypeID: d.LicenseTypeID,
		OEMTag:        d.OemTag,
		Remark:        d.Remark,
		CreatedAt:     d.CreatedAt,
		CreatedBy:     d.CreatedBy,
		UpdatedAt:     d.UpdatedAt,
		UpdatedBy:     d.UpdatedBy,
	}

	// 添加关联信息
	if d.Edges.Product != nil {
		deviceInfo.ProductName = d.Edges.Product.ProductName
		deviceInfo.ProductCode = d.Edges.Product.Code
	}

	if d.Edges.LicenseType != nil {
		deviceInfo.LicenseTypeName = d.Edges.LicenseType.TypeName
		deviceInfo.LicenseTypeCode = d.Edges.LicenseType.LicenseType
	}

	if d.Edges.Creator != nil {
		deviceInfo.CreatedByEmail = d.Edges.Creator.Email
	}

	if d.Edges.Updater != nil {
		deviceInfo.UpdatedByEmail = d.Edges.Updater.Email
	}

	return &deviceInfo, resource.CODE_SUCCESS
}

// GetLicenseTypesByProductID 获取产品下的许可证类型
func (s *DeviceService) GetLicenseTypesByProductID(c *gin.Context, userID int, productID int) ([]*ent.LicenseType, resource.RspCode) {
	// 权限检查
	if userID != 1 {
		exist, err := dto.Client().ProductManager.Query().
			Where(
				productmanager.ProductIDEQ(productID),
				productmanager.UserIDEQ(userID),
			).Exist(c)
		if err != nil || !exist {
			return nil, resource.ERR_NO_PERMISSION
		}
	}

	// 查询产品下的许可证类型
	licenseTypes, err := dto.Client().LicenseType.Query().
		Where(licensetype.ProductIDEQ(productID)).
		Order(ent.Asc(licensetype.FieldTypeName)).
		All(c)

	if err != nil {
		logger.Error("query license types failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	return licenseTypes, resource.CODE_SUCCESS
}

// AddDevice 添加单个设备
func (s *DeviceService) AddDevice(c *gin.Context, userID int, param dto.DeviceAdd) resource.RspCode {
	// 权限检查
	if userID != 1 {
		exist, err := dto.Client().ProductManager.Query().
			Where(
				productmanager.ProductIDEQ(param.ProductID),
				productmanager.UserIDEQ(userID),
			).Exist(c)
		if err != nil || !exist {
			return resource.ERR_NO_PERMISSION
		}
	}

	// 检查产品是否存在
	productExist, err := dto.Client().Product.Query().
		Where(product.IDEQ(param.ProductID)).
		Exist(c)
	if err != nil {
		logger.Error("check product failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if !productExist {
		return resource.ERR_PRODUCT_NOT_EXIST
	}

	// 检查许可证类型是否存在
	licenseTypeExist, err := dto.Client().LicenseType.Query().
		Where(
			licensetype.IDEQ(param.LicenseTypeID),
			licensetype.ProductIDEQ(param.ProductID),
		).Exist(c)
	if err != nil {
		logger.Error("check license type failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if !licenseTypeExist {
		return resource.ERR_LICENSE_TYPE_NOT_EXIST
	}

	// 检查SN是否重复
	exist, err := dto.Client().Device.Query().
		Where(device.SnEQ(param.SN)).
		Exist(c)
	if err != nil {
		logger.Error("check device sn failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if exist {
		return resource.ERR_DEVICE_SN_EXIST
	}

	// 开启事务
	tx, err := dto.Client().Tx(c)
	if err != nil {
		logger.Error("begin transaction failed", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	now := time.Now()

	// 创建设备
	newDevice, err := tx.Device.Create().
		SetSn(param.SN).
		SetProductID(param.ProductID).
		SetLicenseTypeID(param.LicenseTypeID).
		SetOemTag(param.OEMTag).
		SetRemark(param.Remark).
		SetCreatedAt(now).
		SetCreatedBy(userID).
		SetUpdatedAt(now).
		SetUpdatedBy(userID).
		Save(c)

	if err != nil {
		logger.Error("create device failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_ADD_FAILED
	}

	// 创建审计日志
	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:    userID,
		Action:    dto.ActionCreate,
		Module:    dto.ModuleDevice,
		ProductID: param.ProductID,
		DetailInfo: map[string]interface{}{
			"device": newDevice,
		},
	})

	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_ADD_LOG_FAILED
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("commit transaction failed", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}

	return resource.CODE_SUCCESS
}

// BatchAddDevices 批量添加设备
func (s *DeviceService) BatchAddDevices(c *gin.Context, userID int, param dto.DeviceBatchAdd) resource.RspCode {
	// 权限检查
	if userID != 1 {
		exist, err := dto.Client().ProductManager.Query().
			Where(
				productmanager.ProductIDEQ(param.ProductID),
				productmanager.UserIDEQ(userID),
			).Exist(c)
		if err != nil || !exist {
			return resource.ERR_NO_PERMISSION
		}
	}

	// 检查产品是否存在
	productExist, err := dto.Client().Product.Query().
		Where(product.IDEQ(param.ProductID)).
		Exist(c)
	if err != nil {
		logger.Error("check product failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if !productExist {
		return resource.ERR_PRODUCT_NOT_EXIST
	}

	// 检查许可证类型是否存在
	licenseTypeExist, err := dto.Client().LicenseType.Query().
		Where(
			licensetype.IDEQ(param.LicenseTypeID),
			licensetype.ProductIDEQ(param.ProductID),
		).Exist(c)
	if err != nil {
		logger.Error("check license type failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if !licenseTypeExist {
		return resource.ERR_LICENSE_TYPE_NOT_EXIST
	}

	// 过滤空的SN
	var validSNs []string
	for _, sn := range param.SNs {
		sn = strings.TrimSpace(sn)
		if sn != "" {
			validSNs = append(validSNs, sn)
		}
	}

	if len(validSNs) == 0 {
		return resource.ERR_INVALID_PARAMETER
	}

	// 检查SN是否重复
	existingSNs, err := dto.Client().Device.Query().
		Where(device.SnIn(validSNs...)).
		Select(device.FieldSn).
		Strings(c)
	if err != nil {
		logger.Error("check device sn failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	if len(existingSNs) > 0 {
		// 有重复的SN，创建详细的错误信息
		return resource.ERR_DEVICE_SN_EXIST
	}

	// 开启事务
	tx, err := dto.Client().Tx(c)
	if err != nil {
		logger.Error("begin transaction failed", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	now := time.Now()

	// 批量创建设备
	bulk := make([]*ent.DeviceCreate, len(validSNs))
	for i, sn := range validSNs {
		bulk[i] = tx.Device.Create().
			SetSn(sn).
			SetProductID(param.ProductID).
			SetLicenseTypeID(param.LicenseTypeID).
			SetOemTag(param.OEMTag).
			SetRemark(param.Remark).
			SetCreatedAt(now).
			SetCreatedBy(userID).
			SetUpdatedAt(now).
			SetUpdatedBy(userID)
	}

	devices, err := tx.Device.CreateBulk(bulk...).Save(c)
	if err != nil {
		logger.Error("create devices in bulk failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_ADD_FAILED
	}

	// 创建审计日志
	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:    userID,
		Action:    dto.ActionCreate,
		Module:    dto.ModuleDevice,
		ProductID: param.ProductID,
		DetailInfo: map[string]interface{}{
			"count": len(devices),
			"sns":   validSNs,
		},
	})

	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_ADD_LOG_FAILED
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("commit transaction failed", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}

	return resource.CODE_SUCCESS
}

// UpdateDevice 更新设备
func (s *DeviceService) UpdateDevice(c *gin.Context, userID int, param dto.DeviceUpdate) resource.RspCode {
	// 获取设备信息
	d, err := dto.Client().Device.Query().
		Where(device.IDEQ(param.ID)).
		WithProduct().
		Only(c)

	if err != nil {
		if ent.IsNotFound(err) {
			return resource.ERR_DEVICE_NOT_EXIST
		}
		logger.Error("query device failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	// 权限检查
	if userID != 1 {
		exist, err := dto.Client().ProductManager.Query().
			Where(
				productmanager.ProductIDEQ(d.ProductID),
				productmanager.UserIDEQ(userID),
			).Exist(c)
		if err != nil || !exist {
			return resource.ERR_NO_PERMISSION
		}
	}

	// 检查许可证类型是否存在
	licenseTypeExist, err := dto.Client().LicenseType.Query().
		Where(
			licensetype.IDEQ(param.LicenseTypeID),
			licensetype.ProductIDEQ(d.ProductID),
		).Exist(c)
	if err != nil {
		logger.Error("check license type failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if !licenseTypeExist {
		return resource.ERR_LICENSE_TYPE_NOT_EXIST
	}

	// 开启事务
	tx, err := dto.Client().Tx(c)
	if err != nil {
		logger.Error("begin transaction failed", zap.Error(err))
		return resource.ERR_MOD_FAILED
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	oldDevice := *d

	// 更新设备
	updatedDevice, err := tx.Device.UpdateOne(d).
		SetLicenseTypeID(param.LicenseTypeID).
		SetOemTag(param.OEMTag).
		SetRemark(param.Remark).
		SetUpdatedAt(time.Now()).
		SetUpdatedBy(userID).
		Save(c)

	if err != nil {
		logger.Error("update device failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_MOD_FAILED
	}

	// 创建审计日志
	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:    userID,
		Action:    dto.ActionUpdate,
		Module:    dto.ModuleDevice,
		ProductID: d.ProductID,
		DetailInfo: map[string]interface{}{
			"old_device": oldDevice,
			"new_device": updatedDevice,
		},
	})

	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_ADD_LOG_FAILED
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("commit transaction failed", zap.Error(err))
		return resource.ERR_MOD_FAILED
	}

	return resource.CODE_SUCCESS
}

// DeleteDevice 删除设备
func (s *DeviceService) DeleteDevice(c *gin.Context, userID int, deviceID int) resource.RspCode {
	// 获取设备信息
	d, err := dto.Client().Device.Query().
		Where(device.IDEQ(deviceID)).
		Only(c)

	if err != nil {
		if ent.IsNotFound(err) {
			return resource.ERR_DEVICE_NOT_EXIST
		}
		logger.Error("query device failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	// 权限检查
	if userID != 1 {
		exist, err := dto.Client().ProductManager.Query().
			Where(
				productmanager.ProductIDEQ(d.ProductID),
				productmanager.UserIDEQ(userID),
				productmanager.RoleEQ(productmanager.RoleMain),
			).Exist(c)
		if err != nil || !exist {
			return resource.ERR_NO_PERMISSION
		}
	}

	// 开启事务
	tx, err := dto.Client().Tx(c)
	if err != nil {
		logger.Error("begin transaction failed", zap.Error(err))
		return resource.ERR_DEL_FAILED
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	// 保存设备信息用于审计日志
	deviceInfo := *d

	// 删除设备
	err = tx.Device.DeleteOne(d).Exec(c)
	if err != nil {
		logger.Error("delete device failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_DEL_FAILED
	}

	// 创建审计日志
	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:    userID,
		Action:    dto.ActionDelete,
		Module:    dto.ModuleDevice,
		ProductID: deviceInfo.ProductID,
		DetailInfo: map[string]interface{}{
			"device": deviceInfo,
		},
	})

	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		_ = tx.Rollback()
		return resource.ERR_ADD_LOG_FAILED
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("commit transaction failed", zap.Error(err))
		return resource.ERR_DEL_FAILED
	}

	return resource.CODE_SUCCESS
}

// BatchUpdateLicenseType 批量更新设备许可证类型
func (s *DeviceService) BatchUpdateLicenseType(c *gin.Context, userID int, param dto.DeviceBatchUpdateLicense) resource.RspCode {
	if len(param.DeviceIDs) == 0 {
		return resource.ERR_INVALID_PARAMETER
	}

	// 获取设备列表
	devices, err := dto.Client().Device.Query().
		Where(device.IDIn(param.DeviceIDs...)).
		WithProduct().
		All(c)

	if err != nil {
		logger.Error("query devices failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	if len(devices) == 0 {
		return resource.ERR_DEVICE_NOT_EXIST
	}

	// 分组设备（按产品ID）
	devicesByProduct := make(map[int][]*ent.Device)
	for _, d := range devices {
		devicesByProduct[d.ProductID] = append(devicesByProduct[d.ProductID], d)
	}

	// 检查权限和许可证类型
	for productID := range devicesByProduct {
		// 权限检查
		if userID != 1 {
			exist, err := dto.Client().ProductManager.Query().
				Where(
					productmanager.ProductIDEQ(productID),
					productmanager.UserIDEQ(userID),
				).Exist(c)
			if err != nil || !exist {
				return resource.ERR_NO_PERMISSION
			}
		}

		// 检查许可证类型是否存在
		licenseTypeExist, err := dto.Client().LicenseType.Query().
			Where(
				licensetype.IDEQ(param.LicenseTypeID),
				licensetype.ProductIDEQ(productID),
			).Exist(c)
		if err != nil {
			logger.Error("check license type failed", zap.Error(err))
			return resource.ERR_QUERY_FAILED
		}
		if !licenseTypeExist {
			return resource.ERR_LICENSE_TYPE_NOT_EXIST
		}
	}

	// 开启事务
	tx, err := dto.Client().Tx(c)
	if err != nil {
		logger.Error("begin transaction failed", zap.Error(err))
		return resource.ERR_MOD_FAILED
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	now := time.Now()

	// 批量更新设备许可证类型
	for _, d := range devices {
		oldDevice := *d

		updatedDevice, err := tx.Device.UpdateOne(d).
			SetLicenseTypeID(param.LicenseTypeID).
			SetRemark(param.Remark).
			SetUpdatedAt(now).
			SetUpdatedBy(userID).
			Save(c)

		if err != nil {
			logger.Error("update device failed", zap.Error(err))
			_ = tx.Rollback()
			return resource.ERR_MOD_FAILED
		}

		// 创建审计日志
		err = CreateAuditLog(c, tx, dto.AuditLogData{
			UserID:    userID,
			Action:    dto.ActionUpdate,
			Module:    dto.ModuleDevice,
			ProductID: d.ProductID,
			DetailInfo: map[string]interface{}{
				"old_device": oldDevice,
				"new_device": updatedDevice,
				"operation":  "batch_update_license",
			},
		})

		if err != nil {
			logger.Error("create audit log failed", zap.Error(err))
			_ = tx.Rollback()
			return resource.ERR_ADD_LOG_FAILED
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("commit transaction failed", zap.Error(err))
		return resource.ERR_MOD_FAILED
	}

	return resource.CODE_SUCCESS
}

// GetActivationFile 获取设备激活文件 公钥加密私钥解密
func (s *DeviceService) GetActivationFile(c *gin.Context, sn string) ([]byte, resource.RspCode) {
	// 获取设备信息
	device, err := dto.Client().Device.Query().
		//Where(device.SnEQ(sn), device.ProductIDEQ(productID)).
		Where(device.SnEQ(sn)).
		WithProduct().
		WithLicenseType().
		Only(c)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, resource.ERR_DEVICE_NOT_EXIST
		}
		logger.Error("query device failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 获取许可证类型对应的功能编码
	features, err := dto.Client().LicenseType.Query().
		Where(licensetype.IDEQ(device.LicenseTypeID)).
		QueryFeatures().
		All(c)
	if err != nil {
		logger.Error("query features failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 提取功能编码
	featureCodes := make([]string, 0, len(features))
	for _, f := range features {
		featureCodes = append(featureCodes, f.FeatureCode)
	}

	// 生成激活文件内容
	activationData := dto.ActivationData{
		SN:           device.Sn,
		ProductID:    device.ProductID,
		LicenseType:  device.LicenseTypeID,
		OEMTag:       device.OemTag,
		CreatedAt:    time.Now().Unix(),
		FeatureCodes: featureCodes,
	}

	// 将数据转换为JSON，用于签名
	jsonData, err := jsoniter.Marshal(activationData)
	if err != nil {
		logger.Error("marshal activation data failed", zap.Error(err))
		return nil, resource.ERR_OPERATION_FAILED
	}

	// 生成签名
	signature, err := signData(jsonData)
	if err != nil {
		logger.Error("sign activation data failed", zap.Error(err))
		return nil, resource.ERR_OPERATION_FAILED
	}

	// 构建激活文件
	activationFile := &dto.ActivationFile{
		Data:      activationData,
		Signature: signature,
	}
	jsonEnc, _ := jsoniter.Marshal(activationFile)
	enc, err := encrypt(jsonEnc)
	if err != nil {
		return nil, resource.ERR_OPERATION_FAILED
	}
	return enc, resource.CODE_SUCCESS
}

// signData 签名数据
func signData(data []byte) ([]byte, error) {
	// 从配置中获取私钥
	privateKey := resource.Conf.App.GetPrivateKey()
	if privateKey == nil {
		return nil, fmt.Errorf("private key not initialized")
	}

	// 计算哈希
	hash := sha256.Sum256(data)

	// 签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %v", err)
	}

	return signature, nil
}

// 加密（使用AES-GCM模式）
func encrypt(plaintext []byte) ([]byte, error) {
	block, _ := aes.NewCipher([]byte("lqFrzHIimXT66RgpglhASciWerqFMEjJ"))
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize()) // 生成随机Nonce
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil // 将nonce拼接到密文前
}

// 解密
func decrypt(key, ciphertext []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:] // 分离nonce和密文
	return gcm.Open(nil, nonce, ciphertext, nil)
}
