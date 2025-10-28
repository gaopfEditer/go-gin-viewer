package controller

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"cambridge-hit.com/gin-base/activateserver/app/service"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/auth"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/req-resp/resp"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"strconv"

	"github.com/gin-gonic/gin"
)

// VersionController 版本管理控制器
type VersionController struct {
	s *service.VersionService
}

// NewVersionController 创建版本管理控制器
func NewVersionController() *VersionController {
	return &VersionController{
		s: service.NewVersionService(),
	}
}

// ListFirmwareVersions
// @Tags     version
// @Summary  获取韧件版本列表
// @Produce  application/json
// @Param    product_id  path     int     true  "产品ID"
// @Param    page        query    int     false "页码，从1开始"   default(1)
// @Param    page_size   query    int     false "每页数量"        default(10)
// @Param    Authorization  header    string  true  "Authorization"
// @Success  200      {object}  resp.Response{data=dto.PageResult{list=[]dto.FirmwareVersionResponse}}  "获取韧件版本列表"
// @Router   /activate/firmware-versions/{product_id} [get]
func (cl *VersionController) ListFirmwareVersions(c *gin.Context) {
	// 获取产品ID
	productIDStr := c.Param("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	// 获取分页参数
	var page, pageSize int
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if pageSizeStr != "" {
		pageSize, _ = strconv.Atoi(pageSizeStr)
	}
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID <= 0 {
		resp.Error(c, resource.ERR_NO_PERMISSION)
		return
	}

	// 调用服务层方法
	result, code := cl.s.ListFirmwareVersions(c, uai.UserID, productID, page, pageSize)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, result)
}

// AddFirmwareVersion
// @Tags     version
// @Summary  添加韧件版本
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    data          body      dto.AddFirmwareVersion  true  "参数：添加韧件版本"
// @Success  200   {object}  resp.Response{message=string}  "添加韧件版本"
// @Router   /activate/firmware-versions [post]
func (cl *VersionController) AddFirmwareVersion(c *gin.Context) {
	// 解析请求参数
	var param dto.AddFirmwareVersion
	if err := c.ShouldBindJSON(&param); err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	// 获取用户ID
	uai := auth.GetUserAuthInfo(c)
	userID := uai.UserID
	if userID <= 0 {
		resp.Error(c, resource.ERR_NO_PERMISSION)
		return
	}

	// 调用服务层方法
	code := cl.s.AddFirmwareVersion(c, userID, param)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, nil)
}

// ModifyFirmwareVersion
// @Tags     version
// @Summary  修改韧件版本
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    data          body      dto.ModifyFirmwareVersion  true  "参数：修改韧件版本"
// @Success  200   {object}  resp.Response{message=string}  "修改韧件版本"
// @Router   /activate/firmware-versions [put]
func (cl *VersionController) ModifyFirmwareVersion(c *gin.Context) {
	// 解析请求参数
	var param dto.ModifyFirmwareVersion
	if err := c.ShouldBindJSON(&param); err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}
	uai := auth.GetUserAuthInfo(c)
	// 获取用户ID
	userID := uai.UserID
	if userID <= 0 {
		resp.Error(c, resource.ERR_NO_PERMISSION)
		return
	}

	// 调用服务层方法
	code := cl.s.ModifyFirmwareVersion(c, userID, param)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, nil)
}

// DeleteFirmwareVersion
// @Tags     version
// @Summary  删除韧件版本
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    id            path      int     true  "韧件版本ID"
// @Success  200   {object}  resp.Response{message=string}  "删除韧件版本"
// @Router   /activate/firmware-versions/{id} [delete]
func (cl *VersionController) DeleteFirmwareVersion(c *gin.Context) {
	// 获取版本ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}
	uai := auth.GetUserAuthInfo(c)
	userID := uai.UserID
	if userID <= 0 {
		resp.Error(c, resource.ERR_NO_PERMISSION)
		return
	}

	// 调用服务层方法
	code := cl.s.DeleteFirmwareVersion(c, userID, id)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, nil)
}

// ListSoftwareVersions
// @Tags     version
// @Summary  获取软件版本列表
// @Produce  application/json
// @Param    product_id  path     int     true  "产品ID"
// @Param    page        query    int     false "页码，从1开始"   default(1)
// @Param    page_size   query    int     false "每页数量"        default(10)
// @Param    Authorization  header    string  true  "Authorization"
// @Success  200      {object}  resp.Response{data=dto.PageResult{list=[]dto.SoftwareVersionResponse}}  "获取软件版本列表"
// @Router   /activate/software-versions/{product_id} [get]
func (cl *VersionController) ListSoftwareVersions(c *gin.Context) {
	// 获取产品ID
	productIDStr := c.Param("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	// 获取分页参数
	var page, pageSize int
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if pageSizeStr != "" {
		pageSize, _ = strconv.Atoi(pageSizeStr)
	}

	uai := auth.GetUserAuthInfo(c)
	userID := uai.UserID
	if userID <= 0 {
		resp.Error(c, resource.ERR_NO_PERMISSION)
		return
	}

	// 调用服务层方法
	result, code := cl.s.ListSoftwareVersions(c, userID, productID, page, pageSize)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, result)
}

// AddSoftwareVersion
// @Tags     version
// @Summary  添加软件版本
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    data          body      dto.AddSoftwareVersion  true  "参数：添加软件版本"
// @Success  200   {object}  resp.Response{message=string}  "添加软件版本"
// @Router   /activate/software-versions [post]
func (cl *VersionController) AddSoftwareVersion(c *gin.Context) {
	// 解析请求参数
	var param dto.AddSoftwareVersion
	if err := c.ShouldBindJSON(&param); err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	uai := auth.GetUserAuthInfo(c)
	userID := uai.UserID
	if userID <= 0 {
		resp.Error(c, resource.ERR_NO_PERMISSION)
		return
	}

	// 调用服务层方法
	code := cl.s.AddSoftwareVersion(c, userID, param)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, nil)
}

// ModifySoftwareVersion
// @Tags     version
// @Summary  修改软件版本
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    data          body      dto.ModifySoftwareVersion  true  "参数：修改软件版本"
// @Success  200   {object}  resp.Response{message=string}  "修改软件版本"
// @Router   /activate/software-versions [put]
func (cl *VersionController) ModifySoftwareVersion(c *gin.Context) {
	// 解析请求参数
	var param dto.ModifySoftwareVersion
	if err := c.ShouldBindJSON(&param); err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	uai := auth.GetUserAuthInfo(c)
	userID := uai.UserID
	if userID <= 0 {
		resp.Error(c, resource.ERR_NO_PERMISSION)
		return
	}

	// 调用服务层方法
	code := cl.s.ModifySoftwareVersion(c, userID, param)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, nil)
}

// DeleteSoftwareVersion
// @Tags     version
// @Summary  删除软件版本
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    id            path      int     true  "软件版本ID"
// @Success  200   {object}  resp.Response{message=string}  "删除软件版本"
// @Router   /activate/software-versions/{id} [delete]
func (cl *VersionController) DeleteSoftwareVersion(c *gin.Context) {
	// 获取版本ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	uai := auth.GetUserAuthInfo(c)
	userID := uai.UserID
	if userID <= 0 {
		resp.Error(c, resource.ERR_NO_PERMISSION)
		return
	}

	// 调用服务层方法
	code := cl.s.DeleteSoftwareVersion(c, userID, id)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, nil)
}

// GetProductFirmwareVersions
// @Tags     version
// @Summary  获取产品的所有韧件版本
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    product_id    path      int     true  "产品ID"
// @Success  200   {object}  resp.Response{data=[]dto.FirmwareInfo}  "获取产品韧件版本列表"
// @Router   /activate/version/firmware/all/{product_id} [get]
func (cl *VersionController) GetProductFirmwareVersions(c *gin.Context) {
	// 获取产品ID
	productIDStr := c.Param("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	// 调用服务层方法
	result, code := cl.s.GetProductFirmwareVersions(c, productID)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, result)
}

// GetProductFeatures
// @Tags     version
// @Summary  获取产品的所有功能
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    product_id    path      int     true  "产品ID"
// @Success  200   {object}  resp.Response{data=[]dto.FeatureInfo}  "获取产品功能列表"
// @Router   /activate/version/features/{product_id} [get]
func (cl *VersionController) GetProductFeatures(c *gin.Context) {
	// 获取产品ID
	productIDStr := c.Param("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	// 调用服务层方法
	result, code := cl.s.GetProductFeatures(c, productID)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, result)
}
