package controller

import (
	"strconv"

	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"cambridge-hit.com/gin-base/activateserver/app/service"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/auth"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/req-resp/resp"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
)

type PostController struct {
	s *service.PostService
}

func NewPostController() *PostController {
	return &PostController{s: service.NewPostService()}
}

// ListPost 获取文章列表
// @Tags     post
// @Summary  获取文章列表
// @Produce  application/json
// @Param    page     query    int     false  "页码，从1开始"   default(1)
// @Param    page_size query    int     false  "每页数量"        default(10)
// @Param    search   query    string  false  "搜索关键字"
// @Param    status   query    string  false  "文章状态"
// @Param    category_id query int     false  "类别ID"
// @Param    tag_id   query    int     false  "标签ID"
// @Param    Authorization  header    string  true  "Authorization"
// @Success  200      {object}  resp.Response  "获取文章列表"
// @Router   /activate/post/list [get]
func (cl *PostController) ListPost(c *gin.Context) {
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

	// 构建查询参数
	query := dto.PostQuery{
		Page:     page,
		PageSize: pageSize,
		Search:   c.Query("search"),
		Status:   c.Query("status"),
	}

	if categoryID, err := strconv.Atoi(c.Query("category_id")); err == nil {
		query.CategoryID = categoryID
	}

	if tagID, err := strconv.Atoi(c.Query("tag_id")); err == nil {
		query.TagID = tagID
	}

	result, code := cl.s.ListPost(c, uai.UserID, query)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, result)
}

// AddPost 添加文章
// @Tags     post
// @Summary  添加文章
// @Produce  application/json
// @Param    Authorization header     string true "Authorization"
// @Param    data  body      dto.AddPost   true  "参数：添加文章"
// @Success  200   {object}  resp.Response{message=string}  "添加文章"
// @Router   /activate/post/add [post]
func (cl *PostController) AddPost(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	var param dto.AddPost
	if err := c.ShouldBindJSON(&param); err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	code := cl.s.AddPost(c, uai.UserID, param)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c)
}

// UpdatePost 更新文章
// @Tags     post
// @Summary  更新文章
// @Produce  application/json
// @Param    Authorization header     string true "Authorization"
// @Param    data  body      dto.ModifyPost   true "参数：更新文章"
// @Success  200   {object}  resp.Response{message=string}  "更新文章"
// @Router   /activate/post/update [post]
func (cl *PostController) UpdatePost(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	var param dto.ModifyPost
	if err := c.ShouldBindJSON(&param); err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	if param.ID == 0 {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	code := cl.s.UpdatePost(c, uai.UserID, param)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c)
}

// DeletePost 删除文章
// @Tags     post
// @Summary  删除文章
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    id     query     string  true  "文章ID"
// @Success  200    {object}  resp.Response{message=string}  "删除文章"
// @Router   /activate/post/delete [delete]
func (cl *PostController) DeletePost(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id == 0 {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	code := cl.s.DeletePost(c, uai.UserID, id)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c)
}

// GetPost 获取文章详情
// @Tags     post
// @Summary  获取文章详情
// @Produce  application/json
// @Param    Authorization  header    string  true  "Authorization"
// @Param    id     query     string  true  "文章ID"
// @Success  200    {object}  resp.Response{message=string}  "获取文章详情"
// @Router   /activate/post/detail [get]
func (cl *PostController) GetPost(c *gin.Context) {
	uai := auth.GetUserAuthInfo(c)
	if uai.UserID == 0 {
		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
		return
	}

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id == 0 {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	result, code := cl.s.GetPost(c, uai.UserID, id)
	if code != resource.CODE_SUCCESS {
		resp.Error(c, code)
		return
	}

	resp.Success(c, result)
}
