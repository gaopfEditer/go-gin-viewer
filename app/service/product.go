package service

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/product"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/productmanager"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/user"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 定义接口

type ProductService struct{}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (s *ProductService) ListProduct(c *gin.Context, userID int, page, pageSize int) (*dto.PageResult, resource.RspCode) {
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

	// 构建分页结果
	result := &dto.PageResult{
		Total:    int64(total),
		Page:     page,
		PageSize: pageSize,
		List:     products,
	}

	return result, resource.CODE_SUCCESS
}

func (s *ProductService) AddProduct(c *gin.Context, userID int, param dto.AddProduct) resource.RspCode {
	if param.Code == "" || param.ProductName == "" {
		return resource.ERR_INVALID_PARAMETER
	}
	if param.ProductType == "" {
		param.ProductType = product.DefaultProductType
	}
	// 检查产品代号是否重复
	exist, err := dto.Client().Product.Query().Where(product.CodeEQ(param.Code)).Exist(c)
	if err != nil {
		logger.Error("check product code failed", zap.Any("err:", err))
		return resource.ERR_QUERY_FAILED
	}
	if exist {
		return resource.ERR_PRODUCT_CODE_EXIST
	}

	// 检查产品名称是否重复
	exist, err = dto.Client().Product.Query().Where(product.ProductNameEQ(param.ProductName)).Exist(c)
	if err != nil {
		logger.Error("check product name failed", zap.Any("err:", err))
		return resource.ERR_QUERY_FAILED
	}
	if exist {
		return resource.ERR_PRODUCT_NAME_EXIST
	}

	// 创建产品
	tx, err := dto.Client().Tx(c)
	if err != nil {
		logger.Error("begin transaction failed", zap.Any("err:", err))
		return resource.ERR_ADD_FAILED
	}
	defer func() {
		if r := recover(); r != nil || err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Error("rollback transaction failed", zap.Any("err:", err))
			}
		}
	}()

	p, err := tx.Product.Create().
		SetCode(param.Code).
		SetProductName(param.ProductName).
		SetProductType(param.ProductType).
		Save(c)

	if err != nil {
		logger.Error("create product failed", zap.Any("err:", err))
		if err := tx.Rollback(); err != nil {
			logger.Error("rollback transaction failed", zap.Any("err:", err))
		}
		return resource.ERR_ADD_FAILED
	}

	// 设置创建者为产品管理员
	_, err = tx.ProductManager.Create().
		SetUserID(userID).
		SetProductID(p.ID).
		SetRole(productmanager.RoleMain).
		SetPermissions(productmanager.PermissionsFull).
		Save(c)

	if err != nil {
		logger.Error("create product manager failed", zap.Any("err:", err))
		if err := tx.Rollback(); err != nil {
			logger.Error("rollback transaction failed", zap.Any("err:", err))
		}
		return resource.ERR_ADD_FAILED
	}

	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:     userID,
		Action:     dto.ActionCreate,
		Module:     dto.ModuleProduct,
		ProductID:  p.ID,
		DetailInfo: nil,
	})

	if err := tx.Commit(); err != nil {
		logger.Error("commit transaction failed", zap.Any("err:", err))
		return resource.ERR_ADD_FAILED
	}

	return resource.CODE_SUCCESS
}

func (s *ProductService) ModifyProduct(c *gin.Context, userID int, param dto.ModifyProduct) resource.RspCode {
	// 开启事务
	tx, err := dto.Client().Tx(c)
	if err != nil {
		logger.Error("failed to begin transaction", zap.Error(err))
		return resource.ERR_MOD_FAILED
	}
	defer func() {
		if v := recover(); v != nil || err != nil {
			tx.Rollback()
		}
	}()

	// 1. 权限检查
	pm, err := tx.ProductManager.Query().Where(
		productmanager.ProductIDEQ(param.ID),
		productmanager.UserIDEQ(userID),
	).Only(c)
	if err != nil || (userID != 1 && pm.Role != productmanager.RoleMain) {
		tx.Rollback()
		return resource.ERR_NO_PERMISSION
	}

	// 2. 获取产品信息
	_product, err := tx.Product.Query().WithManagers(func(pmq *ent.ProductManagerQuery) { pmq.WithUser() }).Where(product.IDEQ(param.ID)).Only(c)
	if err != nil {
		logger.Error("query product failed", zap.Error(err))
		tx.Rollback()
		return resource.ERR_QUERY_FAILED
	}

	// 3. 主管理员或超级管理员操作
	if pm.Role == productmanager.RoleMain || userID == dto.SuperAdminID {
		// 查询现有主管理员
		mainManager, err := tx.ProductManager.Query().
			Where(
				productmanager.ProductIDEQ(param.ID),
				productmanager.RoleEQ(productmanager.RoleMain),
			).Only(c)
		if err != nil {
			logger.Error("query main manager failed", zap.Error(err))
			tx.Rollback()
			return resource.ERR_QUERY_FAILED
		}

		newMainManager := mainManager
		if param.ManagerMain != 0 {
			// 查询新主管理员
			newMainManager, err = tx.ProductManager.Query().
				Where(
					productmanager.ProductIDEQ(param.ID),
					productmanager.UserIDEQ(param.ManagerMain),
				).Only(c)
			if err != nil {
				logger.Error("query new main manager failed", zap.Error(err))
				tx.Rollback()
				return resource.ERR_QUERY_FAILED
			}
		}

		// 更新主管理员角色
		if mainManager.UserID != newMainManager.UserID {
			if _, err = tx.ProductManager.UpdateOne(mainManager).
				SetRole(productmanager.RoleAssistant).
				SetPermissions(productmanager.DefaultPermissions).
				Save(c); err != nil {
				logger.Error("update old main failed", zap.Error(err))
				tx.Rollback()
				return resource.ERR_MOD_FAILED
			}

			if _, err = tx.ProductManager.UpdateOne(newMainManager).
				SetRole(productmanager.RoleMain).
				SetPermissions(productmanager.PermissionsFull).
				SetRemark("").
				Save(c); err != nil {
				logger.Error("update new main failed", zap.Error(err))
				tx.Rollback()
				return resource.ERR_MOD_FAILED
			}
		}

		// 更新副管理员权限和备注
		for _, assistant := range param.ManagerAssistant {
			if assistant.UserID == newMainManager.UserID {
				continue // 跳过主管理员
			}

			assistantManager, err := tx.ProductManager.Query().
				Where(
					productmanager.ProductIDEQ(param.ID),
					productmanager.UserIDEQ(assistant.UserID),
				).Only(c)
			if err != nil {
				logger.Error("query assistant manager failed", zap.Error(err))
				continue
			}

			update := tx.ProductManager.UpdateOne(assistantManager)
			if assistant.Permission != "" {
				update.SetPermissions(assistant.Permission)
			}
			if assistant.Remark != "" {
				update.SetRemark(assistant.Remark)
			}

			if _, err = update.Save(c); err != nil {
				logger.Error("update assistant failed", zap.Error(err))
				continue
			}
		}
	}

	// 4. 更新产品信息
	update := tx.Product.UpdateOne(_product)
	if param.ProductName != "" && param.ProductName != _product.ProductName {
		exist, err := tx.Product.Query().
			Where(product.ProductNameEQ(param.ProductName)).
			Exist(c)
		if err != nil {
			logger.Error("check product name failed", zap.Error(err))
			tx.Rollback()
			return resource.ERR_QUERY_FAILED
		}
		if exist {
			tx.Rollback()
			return resource.ERR_PRODUCT_NAME_EXIST
		}
		update = update.SetProductName(param.ProductName)
	}

	if param.ProductType != _product.ProductType {
		if param.ProductType == "" {
			param.ProductType = product.DefaultProductType
		}
		update = update.SetProductType(param.ProductType)
	}

	if _, err := update.Save(c); err != nil {
		logger.Error("update product failed", zap.Error(err))
		tx.Rollback()
		return resource.ERR_MOD_FAILED
	}

	_newProduct, err := tx.Product.Query().WithManagers(func(pmq *ent.ProductManagerQuery) { pmq.WithUser() }).Where(product.IDEQ(param.ID)).Only(c)
	if err != nil {
		logger.Error("query product failed", zap.Error(err))
		tx.Rollback()
		return resource.ERR_QUERY_FAILED
	}

	// 创建审计日志
	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:    userID,
		Action:    dto.ActionUpdate,
		Module:    dto.ModuleProduct,
		ProductID: param.ID,
		DetailInfo: map[string]interface{}{
			"old_product": _product,
			"new_product": _newProduct,
		},
	})
	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		tx.Rollback()
		return resource.ERR_ADD_LOG_FAILED
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		logger.Error("commit transaction failed", zap.Error(err))
		return resource.ERR_MOD_FAILED
	}

	return resource.CODE_SUCCESS
}

func (s *ProductService) DeleteProduct(c *gin.Context, userID, productID int) resource.RspCode {
	if exist, _ := dto.Client().ProductManager.Query().
		Where( //product_id、user_id,role_main符合
			productmanager.ProductIDEQ(productID),
			productmanager.UserIDEQ(userID),
			productmanager.RoleEQ(productmanager.RoleMain),
		).Exist(c); userID != 1 && !exist {
		return resource.ERR_NO_PERMISSION
	}

	// 获取产品信息
	p, err := dto.Client().Product.Query().
		Where(product.IDEQ(productID)).
		WithSoftwareVersions().
		WithFirmwareVersions().
		WithLicenseTypes().
		WithFeatures().
		Only(c)

	if err != nil {
		logger.Error("query product failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	// 检查是否有关联的软硬件版本、许可证类型和产品功能
	if len(p.Edges.SoftwareVersions) > 0 ||
		len(p.Edges.FirmwareVersions) > 0 ||
		len(p.Edges.LicenseTypes) > 0 ||
		len(p.Edges.Features) > 0 {
		return resource.ERR_PRODUCT_HAS_RELATIONS
	}

	// 事务处理
	tx, err := dto.Client().Tx(c)
	if err != nil {
		logger.Error("begin transaction failed", zap.Any("err:", err))
		return resource.ERR_DEL_FAILED
	}

	// 先删除产品管理员关系
	_, err = tx.ProductManager.Delete().Where(productmanager.ProductIDEQ(productID)).Exec(c)
	if err != nil {
		logger.Error("delete product managers failed", zap.Error(err))
		tx.Rollback()
		return resource.ERR_DEL_FAILED
	}

	// 删除产品本身
	err = tx.Product.DeleteOneID(productID).Exec(c)
	if err != nil {
		logger.Error("delete product failed", zap.Error(err))
		tx.Rollback()
		return resource.ERR_DEL_FAILED
	}

	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:     userID,
		Action:     dto.ActionDelete,
		Module:     dto.ModuleProduct,
		ProductID:  0,
		DetailInfo: p,
	})
	if err != nil {
		logger.Error("create audit log failed", zap.Error(err))
		tx.Rollback()
		return resource.ERR_ADD_LOG_FAILED
	}
	// 提交事务
	if err = tx.Commit(); err != nil {
		logger.Error("commit transaction failed", zap.Error(err))
		return resource.ERR_DEL_FAILED
	}
	return resource.CODE_SUCCESS
}

// AddManager 添加产品管理员
// 参数：当前用户ID、添加管理员请求参数
// 返回：响应状态码
func (s *ProductService) AddManager(c *gin.Context, userID int, param dto.AddManager) resource.RspCode {
	// 1. 权限检查：验证当前用户是否是产品的主管理员
	isMainManager, err := dto.Client().ProductManager.Query().
		Where(
			productmanager.ProductIDEQ(param.ProductID),
			productmanager.UserIDEQ(userID),
			productmanager.RoleEQ(productmanager.RoleMain),
		).Exist(c)

	// 如果不是主管理员且不是超级管理员(ID=1)，则无权添加
	if err != nil || (!isMainManager && userID != 1) {
		logger.Error("permission check failed", zap.Error(err))
		return resource.ERR_NO_PERMISSION
	}

	// 2. 根据邮箱查找用户
	u, err := dto.Client().User.Query().
		Where(user.EmailEQ(param.Email)).
		Only(c)
	if err != nil {
		logger.Error("query user by email failed", zap.String("email", param.Email), zap.Error(err))
		return resource.ERR_USER_NOT_EXIST
	}

	// 3. 检查该用户是否已经是该产品的管理员
	exist, err := dto.Client().ProductManager.Query().
		Where(
			productmanager.ProductIDEQ(param.ProductID),
			productmanager.UserIDEQ(u.ID),
		).Exist(c)
	if err != nil {
		logger.Error("check existing manager failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if exist {
		return resource.ERR_MANAGER_ALREADY_EXIST
	}

	// 4. 验证权限参数合法性
	if param.Permissions == "" {
		param.Permissions = productmanager.PermissionsRead // 默认只读权限
	} else if err := productmanager.PermissionsValidator(param.Permissions); err != nil {
		return resource.ERR_INVALID_PARAMETER
	}

	tx, _ := dto.Client().Tx(c)

	// 5. 添加用户为产品管理员
	_, err = tx.ProductManager.Create().
		SetUserID(u.ID).
		SetProductID(param.ProductID).
		SetRole(productmanager.RoleAssistant). // 新添加的用户始终是副管理员
		SetPermissions(param.Permissions).
		SetRemark(param.Remark). // 设置备注
		Save(c)
	if err != nil {
		logger.Error("add manager failed", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}

	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:     userID,
		Action:     dto.ActionCreate,
		Module:     dto.ModuleProduct,
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

// RemoveManager 删除产品管理员
// 参数：当前用户ID、管理员ID
// 返回：响应状态码
func (s *ProductService) RemoveManager(c *gin.Context, userID int, managerId int) resource.RspCode {
	// 1. 获取要删除的管理员信息
	manager, err := dto.Client().ProductManager.Get(c, managerId)
	if err != nil {
		logger.Error("query manager failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	// 2. 不能删除主管理员
	if manager.Role == productmanager.RoleMain {
		return resource.ERR_NO_PERMISSION
	}

	// 3. 权限检查：验证当前用户是否是产品的主管理员或超级管理员
	if userID != 1 {
		isMainManager, err := dto.Client().ProductManager.Query().
			Where(
				productmanager.ProductIDEQ(manager.ProductID),
				productmanager.UserIDEQ(userID),
				productmanager.RoleEQ(productmanager.RoleMain),
			).Exist(c)

		if err != nil || !isMainManager {
			logger.Error("permission check failed", zap.Error(err))
			return resource.ERR_NO_PERMISSION
		}
	}

	tx, _ := dto.Client().Tx(c)

	// 4. 删除管理员
	err = tx.ProductManager.DeleteOne(manager).Exec(c)
	if err != nil {
		logger.Error("delete manager failed", zap.Error(err))
		return resource.ERR_DEL_FAILED
	}

	err = CreateAuditLog(c, tx, dto.AuditLogData{
		UserID:    userID,
		Action:    dto.ActionDelete,
		Module:    dto.ModuleProduct,
		ProductID: manager.ProductID,
		DetailInfo: map[string]interface{}{
			"removed_manager": manager,
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
		return resource.ERR_DEL_FAILED
	}
	return resource.CODE_SUCCESS
}
