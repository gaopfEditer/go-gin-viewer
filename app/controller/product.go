package controller

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"strconv"

	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/productmanager"
	"cambridge-hit.com/gin-base/activateserver/app/service"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/auth"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/req-resp/resp"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	s *service.ProductService
}

func NewProductController() *ProductController {
	return &ProductController{s: service.NewProductService()}
}

// ListProduct
// @Tags     product
// @Summary  获取用户的产品列表
// @Produce  application/json
// @Param    page     query    int     false  "页码，从1开始"   default(1)
// @Param    page_size query    int     false  "每页数量"        default(10)
// @Param    Authorization  header    string  true  "Authorization"
// @Success  200      {object}  resp.Response  "获取产品列表"
// @Router   /activate/product/list [get]
func (cl *ProductController) ListProduct(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
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

	result, code := cl.s.ListProduct(c, uai.UserID, page, pageSize)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, result)
}

// AddProduct
// @Tags     Product
// @Summary  增加产品
// @Produce   application/json
// @Param    Authorization header     string true "Authorization"
// @Param    data  body      dto.AddProduct   true  "参数：增加产品"
// @Success  200   {object}  resp.Response{message=string}  "增加产品"
// @Router   /activate/product/add [post]
func (cl *ProductController) AddProduct(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}
	var param dto.AddProduct
	if err := c.ShouldBindJSON(&param); err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}
	err := cl.s.AddProduct(c, uai.UserID, param)
	if err != resource.CODE_SUCCESS {
		resp.Error(c, err)
		return
	}
	resp.Success(c)
}

// ModifyProduct
// @Tags     Product
// @Summary  修改产品
// @Produce  application/json
// @Param    Authorization header     string true "Authorization"
// @Param    data  body      dto.ModifyProduct   true "参数：修改产品"
// @Success  200   {object}  resp.Response{message=string}  "修改产品"
// @Router   /activate/product/put [post]
func (cl *ProductController) ModifyProduct(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	var param dto.ModifyProduct
	if err := c.ShouldBindJSON(&param); err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	if param.ID == 0 {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	for _, ro := range param.ManagerAssistant {
		if err := productmanager.PermissionsValidator(ro.Permission); ro.UserID == 0 || ro.Permission != "" && err != nil {
			resp.Error(c, resource.ERR_INVALID_PARAMETER)
			return
		}
	}

	e := cl.s.ModifyProduct(c, uai.UserID, param)
	if e != resource.CODE_SUCCESS {
		resp.Error(c, e)
		return
	}
	resp.Success(c)
}

// DeleteProduct
// @Tags     Product
// @Summary  删除产品
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    product_id     query     string  true  "产品ID"
// @Success  200    {object}  resp.Response{message=string}  "删除产品"
// @Router   /activate/product/del [get]
func (cl *ProductController) DeleteProduct(c *gin.Context) {
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

	e := cl.s.DeleteProduct(c, uai.UserID, productID)
	if e != resource.CODE_SUCCESS {
		resp.Error(c, e)
		return
	}
	resp.Success(c)
}

// AddManager
// @Tags     Product
// @Summary  添加产品管理员
// @Produce  application/json
// @Param    Authorization header     string true "Authorization"
// @Param    data  body      dto.AddManager   true  "参数：添加产品管理员"
// @Success  200   {object}  resp.Response{message=string}  "添加产品管理员"
// @Router   /activate/product/add-manager [post]
func (cl *ProductController) AddManager(c *gin.Context) {
	// 获取当前用户信息
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	// 绑定请求参数
	var param dto.AddManager
	if err := c.ShouldBindJSON(&param); err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	// 调用服务层方法添加管理员
	e := cl.s.AddManager(c, uai.UserID, param)
	if e != resource.CODE_SUCCESS {
		resp.Error(c, e)
		return
	}

	resp.Success(c)
}

// RemoveManager
// @Tags     Product
// @Summary  删除产品管理员
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    manager_id     query     string  true  "管理员ID"
// @Success  200    {object}  resp.Response{message=string}  "删除产品管理员"
// @Router   /activate/product/remove-manager [get]
func (cl *ProductController) RemoveManager(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	managerId, err := strconv.Atoi(c.Query("manager_id"))
	if err != nil || managerId == 0 {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	e := cl.s.RemoveManager(c, uai.UserID, managerId)
	if e != resource.CODE_SUCCESS {
		resp.Error(c, e)
		return
	}
	resp.Success(c)
}
