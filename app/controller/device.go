package controller

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"cambridge-hit.com/gin-base/activateserver/app/service"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/auth"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/req-resp/resp"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// DeviceController 设备控制器
type DeviceController struct {
	deviceService *service.DeviceService
}

// NewDeviceController 创建设备控制器
func NewDeviceController() *DeviceController {
	return &DeviceController{
		deviceService: service.NewDeviceService(),
	}
}

// ListProducts
// @Tags     device
// @Summary  获取产品列表（带设备数量）
// @Produce  application/json
// @Param    page     query    int     false  "页码，从1开始"   default(1)
// @Param    page_size query    int     false  "每页数量"        default(10)
// @Param    Authorization  header    string  true  "Authorization"
// @Success  200      {object}  resp.Response  "获取产品列表"
// @Router   /activate/device/products [get]
func (c *DeviceController) ListProducts(ctx *gin.Context) {
	uai := auth.GetUserAuthInfo(ctx)
	if uai.UserID == 0 {
		resp.Error(ctx, resource.ERR_TOKEN_EXPIRED)
		return
	}

	// 分页
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	result, code := c.deviceService.ListProducts(ctx, uai.UserID, page, pageSize)
	if code != resource.CODE_SUCCESS {
		resp.Error(ctx, code)
		return
	}

	resp.Success(ctx, result)
}

// ListDevices
// @Tags     device
// @Summary  获取产品下的设备列表
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    product_id     query     int     false  "产品ID"
// @Param    license_type_id query     int     false  "许可证类型ID"
// @Param    sn             query     string  false  "设备序列号"
// @Param    page     query    int     false  "页码，从1开始"   default(1)
// @Param    page_size query    int     false  "每页数量"        default(10)
// @Success  200      {object}  resp.Response  "获取设备列表"
// @Router   /activate/device/list [get]
func (c *DeviceController) ListDevices(ctx *gin.Context) {
	uai := auth.GetUserAuthInfo(ctx)
	if uai.UserID == 0 {
		resp.Error(ctx, resource.ERR_TOKEN_EXPIRED)
		return
	}

	var filter dto.DeviceFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		resp.Error(ctx, resource.ERR_INVALID_PARAMETER)
		return
	}

	result, code := c.deviceService.ListDevices(ctx, uai.UserID, filter)
	if code != resource.CODE_SUCCESS {
		resp.Error(ctx, code)
		return
	}

	resp.Success(ctx, result)
}

// GetDeviceBySN
// @Tags     device
// @Summary  通过SN获取设备信息
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    sn             query     string  true  "设备序列号"
// @Success  200      {object}  resp.Response  "获取设备信息"
// @Router   /activate/device/search [get]
func (c *DeviceController) GetDeviceBySN(ctx *gin.Context) {
	uai := auth.GetUserAuthInfo(ctx)
	if uai.UserID == 0 {
		resp.Error(ctx, resource.ERR_TOKEN_EXPIRED)
		return
	}

	sn := ctx.Query("sn")
	if sn == "" {
		resp.Error(ctx, resource.ERR_INVALID_PARAMETER)
		return
	}

	result, code := c.deviceService.GetDeviceBySN(ctx, uai.UserID, sn)
	if code != resource.CODE_SUCCESS {
		resp.Error(ctx, code)
		return
	}

	resp.Success(ctx, result)
}

// GetLicenseTypes
// @Tags     device
// @Summary  获取产品下的许可证类型列表
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    product_id     query     int     true  "产品ID"
// @Success  200      {object}  resp.Response  "获取许可证类型列表"
// @Router   /activate/device/license-types [get]
func (c *DeviceController) GetLicenseTypes(ctx *gin.Context) {
	uai := auth.GetUserAuthInfo(ctx)
	if uai.UserID == 0 {
		resp.Error(ctx, resource.ERR_TOKEN_EXPIRED)
		return
	}

	productID, err := strconv.Atoi(ctx.Query("product_id"))
	if err != nil || productID <= 0 {
		resp.Error(ctx, resource.ERR_INVALID_PARAMETER)
		return
	}

	result, code := c.deviceService.GetLicenseTypesByProductID(ctx, uai.UserID, productID)
	if code != resource.CODE_SUCCESS {
		resp.Error(ctx, code)
		return
	}

	resp.Success(ctx, result)
}

// AddDevice
// @Tags     device
// @Summary  添加单个设备
// @Produce  application/json
// @Param    Authorization header     string true "Authorization"
// @Param    data  body      dto.DeviceAdd   true  "参数：添加设备"
// @Success  200   {object}  resp.Response{message=string}  "添加设备"
// @Router   /activate/device/add [post]
func (c *DeviceController) AddDevice(ctx *gin.Context) {
	uai := auth.GetUserAuthInfo(ctx)
	if uai.UserID == 0 {
		resp.Error(ctx, resource.ERR_TOKEN_EXPIRED)
		return
	}

	var param dto.DeviceAdd
	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp.Error(ctx, resource.ERR_INVALID_PARAMETER)
		return
	}

	code := c.deviceService.AddDevice(ctx, uai.UserID, param)
	if code != resource.CODE_SUCCESS {
		resp.Error(ctx, code)
		return
	}

	resp.Success(ctx)
}

// BatchAddDevices
// @Tags     device
// @Summary  批量添加设备
// @Produce  application/json
// @Param    Authorization header     string true "Authorization"
// @Param    data  body      dto.DeviceBatchAdd   true  "参数：批量添加设备"
// @Success  200   {object}  resp.Response{message=string}  "批量添加设备"
// @Router   /activate/device/batch-add [post]
func (c *DeviceController) BatchAddDevices(ctx *gin.Context) {
	uai := auth.GetUserAuthInfo(ctx)
	if uai.UserID == 0 {
		resp.Error(ctx, resource.ERR_TOKEN_EXPIRED)
		return
	}

	var param dto.DeviceBatchAdd
	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp.Error(ctx, resource.ERR_INVALID_PARAMETER)
		return
	}

	code := c.deviceService.BatchAddDevices(ctx, uai.UserID, param)
	if code != resource.CODE_SUCCESS {
		resp.Error(ctx, code)
		return
	}

	resp.Success(ctx)
}

// UpdateDevice
// @Tags     device
// @Summary  更新设备
// @Produce  application/json
// @Param    Authorization header     string true "Authorization"
// @Param    data  body      dto.DeviceUpdate   true  "参数：更新设备"
// @Success  200   {object}  resp.Response{message=string}  "更新设备"
// @Router   /activate/device/update [put]
func (c *DeviceController) UpdateDevice(ctx *gin.Context) {
	uai := auth.GetUserAuthInfo(ctx)
	if uai.UserID == 0 {
		resp.Error(ctx, resource.ERR_TOKEN_EXPIRED)
		return
	}

	var param dto.DeviceUpdate
	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp.Error(ctx, resource.ERR_INVALID_PARAMETER)
		return
	}

	code := c.deviceService.UpdateDevice(ctx, uai.UserID, param)
	if code != resource.CODE_SUCCESS {
		resp.Error(ctx, code)
		return
	}

	resp.Success(ctx)
}

// DeleteDevice
// @Tags     device
// @Summary  删除设备
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    id             path      int     true  "设备ID"
// @Success  200    {object}  resp.Response{message=string}  "删除设备"
// @Router   /activate/device/{id} [delete]
func (c *DeviceController) DeleteDevice(ctx *gin.Context) {
	uai := auth.GetUserAuthInfo(ctx)
	if uai.UserID == 0 {
		resp.Error(ctx, resource.ERR_TOKEN_EXPIRED)
		return
	}

	deviceID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || deviceID <= 0 {
		resp.Error(ctx, resource.ERR_INVALID_PARAMETER)
		return
	}

	code := c.deviceService.DeleteDevice(ctx, uai.UserID, deviceID)
	if code != resource.CODE_SUCCESS {
		resp.Error(ctx, code)
		return
	}

	resp.Success(ctx)
}

// BatchUpdateLicenseType
// @Tags     device
// @Summary  批量更新设备许可证类型
// @Produce  application/json
// @Param    Authorization header     string true "Authorization"
// @Param    data  body      dto.DeviceBatchUpdateLicense   true  "参数：批量更新许可证类型"
// @Success  200   {object}  resp.Response{message=string}  "批量更新许可证类型"
// @Router   /activate/device/batch-update-license [post]
func (c *DeviceController) BatchUpdateLicenseType(ctx *gin.Context) {
	uai := auth.GetUserAuthInfo(ctx)
	if uai.UserID == 0 {
		resp.Error(ctx, resource.ERR_TOKEN_EXPIRED)
		return
	}

	var param dto.DeviceBatchUpdateLicense
	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp.Error(ctx, resource.ERR_INVALID_PARAMETER)
		return
	}

	code := c.deviceService.BatchUpdateLicenseType(ctx, uai.UserID, param)
	if code != resource.CODE_SUCCESS {
		resp.Error(ctx, code)
		return
	}

	resp.Success(ctx)
}

// GetActivationFile
// @Tags     device
// @Summary  获取设备激活文件
// @Produce  application/octet-stream
// @Param    productID      path      string  true  "产品ID"
// @Param    sn             path      string  true  "设备序列号"
// @Success  200      {file}   string  "激活文件"
// @Router   /activate/device/activation-file/{sn} [get]
func (c *DeviceController) GetActivationFile(ctx *gin.Context) {
	//uai := auth.GetUserAuthInfo(ctx)
	//if uai.UserID == 0 {
	//	resp.Error(ctx, resource.ERR_TOKEN_EXPIRED)
	//	return
	//}

	sn := ctx.Param("sn")
	//_productID := ctx.Param("productID")
	//productID, err := strconv.Atoi(_productID)
	//if err != nil || sn == "" || productID == 0 {
	//	resp.Error(ctx, resource.ERR_INVALID_PARAMETER)
	//	return
	//}

	result, code := c.deviceService.GetActivationFile(ctx, sn)
	if code != resource.CODE_SUCCESS {
		resp.Error(ctx, code)
		return
	}
	// 设置响应头，使浏览器下载文件
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.lic"`, sn))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Length", fmt.Sprint(len(result)))
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")
	// 直接写入文件内容
	ctx.Data(200, "application/octet-stream", result)
}
