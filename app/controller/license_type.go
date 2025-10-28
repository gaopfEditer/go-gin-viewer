package controller

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"strconv"

	"cambridge-hit.com/gin-base/activateserver/app/service"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/auth"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/req-resp/resp"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
)

type LicenseTypeController struct {
	s *service.LicenseTypeService
}

func NewLicenseTypeController() *LicenseTypeController {
	return &LicenseTypeController{s: service.NewLicenseTypeService()}
}

// ListLicenseTypes
// @Tags     LicenseType
// @Summary  获取许可证类型列表
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    product_id    query     int     true  "产品ID"
// @Param    page          query     int     false "页码，从1开始"   default(1)
// @Param    page_size     query     int     false "每页数量"        default(10)
// @Success  200    {object}  resp.Response  "获取许可证类型列表"
// @Router   /activate/license-type/list [get]
func (cl *LicenseTypeController) ListLicenseTypes(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	productID, err := strconv.Atoi(c.Query("product_id"))
	if err != nil || productID == 0 {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	// 解析分页参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	result, code := cl.s.ListLicenseTypes(c, uai.UserID, productID, page, pageSize)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, result)
}

// AddLicenseType
// @Tags     LicenseType
// @Summary  添加许可证类型
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    data          body      dto.AddLicenseType  true  "添加许可证类型参数"
// @Success  200    {object}  resp.Response  "添加许可证类型"
// @Router   /activate/license-type/add [post]
func (cl *LicenseTypeController) AddLicenseType(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	var param dto.AddLicenseType
	if err := c.ShouldBindJSON(&param); err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	code := cl.s.AddLicenseType(c, uai.UserID, param)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c)
}

// DeleteLicenseType
// @Tags     LicenseType
// @Summary  删除许可证类型
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    type_id       query     string  true  "类型ID"
// @Success  200    {object}  resp.Response  "删除许可证类型"
// @Router   /activate/license-type/del [get]
func (cl *LicenseTypeController) DeleteLicenseType(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	typeID, err := strconv.Atoi(c.Query("type_id"))
	if err != nil || typeID == 0 {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	code := cl.s.DeleteLicenseType(c, uai.UserID, typeID)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c)
}

// UpdateLicenseTypeFeatures
// @Tags     LicenseType
// @Summary  更新许可证类型功能列表
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    data          body      dto.UpdateLicenseTypeFeatures  true  "更新许可证类型功能列表参数"
// @Success  200    {object}  resp.Response  "更新许可证类型功能列表"
// @Router   /activate/license-type/update-features [post]
func (cl *LicenseTypeController) UpdateLicenseTypeFeatures(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	var param dto.UpdateLicenseTypeFeatures
	if err := c.ShouldBindJSON(&param); err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	code := cl.s.UpdateLicenseTypeFeatures(c, uai.UserID, param)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c)
}
