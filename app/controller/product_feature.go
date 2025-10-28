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

type ProductFeatureController struct {
	s *service.ProductFeatureService
}

func NewProductFeatureController() *ProductFeatureController {
	return &ProductFeatureController{s: service.NewProductFeatureService()}
}

// ListProductFeatures
// @Tags     ProductFeature
// @Summary  获取产品功能列表
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    product_id    query     int     true  "产品ID"
// @Param    page          query     int     false "页码，从1开始"   default(1)
// @Param    page_size     query     int     false "每页数量"        default(10)
// @Success  200    {object}  resp.Response  "获取产品功能列表"
// @Router   /activate/product-feature/list [get]
func (cl *ProductFeatureController) ListProductFeatures(c *gin.Context) {
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

	result, code := cl.s.ListProductFeatures(c, uai.UserID, productID, page, pageSize)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, result)
}

// AddProductFeature
// @Tags     ProductFeature
// @Summary  添加产品功能
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    data          body      dto.AddProductFeature  true  "添加产品功能参数"
// @Success  200    {object}  resp.Response  "添加产品功能"
// @Router   /activate/product-feature/add [post]
func (cl *ProductFeatureController) AddProductFeature(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	var param dto.AddProductFeature
	if err := c.ShouldBindJSON(&param); err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	code := cl.s.AddProductFeature(c, uai.UserID, param)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c)
}

// DeleteProductFeature
// @Tags     ProductFeature
// @Summary  删除产品功能
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    feature_id    query     string  true  "功能ID"
// @Success  200    {object}  resp.Response  "删除产品功能"
// @Router   /activate/product-feature/del [get]
func (cl *ProductFeatureController) DeleteProductFeature(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	featureID, err := strconv.Atoi(c.Query("feature_id"))
	if err != nil || featureID == 0 {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	code := cl.s.DeleteProductFeature(c, uai.UserID, featureID)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c)
}
