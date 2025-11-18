package service

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/post"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/posttagrelation"

	// "cambridge-hit.com/gin-base/activateserver/app/entity/ent/postcategory"
	// "cambridge-hit.com/gin-base/activateserver/app/entity/ent/posttag"
	"time"

	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PostService struct{}

func NewPostService() *PostService {
	return &PostService{}
}

// ListPost 获取文章列表
func (s *PostService) ListPost(c *gin.Context, userID int, query dto.PostQuery) (*dto.PageResult, resource.RspCode) {
	// 构建查询
	q := dto.Client().Post.Query()

	// 添加搜索条件
	if query.Search != "" {
		q = q.Where(post.Or(
			post.TitleContainsFold(query.Search),
			post.ContentContainsFold(query.Search),
		))
	}

	// 添加状态过滤
	if query.Status != "" {
		q = q.Where(post.StatusEQ(query.Status))
	}

	// // 添加类别过滤
	// if query.CategoryID > 0 {
	// 	q = q.Where(post.HasCategoryWith(postcategory.IDEQ(query.CategoryID)))
	// }

	// // 添加标签过滤
	// if query.TagID > 0 {
	// 	q = q.Where(post.HasTagsWith(posttag.IDEQ(query.TagID)))
	// }

	// 计算总数
	total, err := q.Count(c)
	if err != nil {
		logger.Error("count posts failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 执行分页查询
	offset := (query.Page - 1) * query.PageSize
	posts, err := q.
		WithCategory().
		WithTagRelations().
		WithAuthor().
		Limit(query.PageSize).
		Offset(offset).
		Order(ent.Desc(post.FieldCreatedAt)).
		All(c)

	if err != nil {
		logger.Error("query posts failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 构建分页结果
	result := &dto.PageResult{
		Total:    int64(total),
		Page:     query.Page,
		PageSize: query.PageSize,
		List:     posts,
	}

	return result, resource.CODE_SUCCESS
}

// AddPost 添加文章
func (s *PostService) AddPost(c *gin.Context, userID int, param dto.AddPost) resource.RspCode {
	if param.Title == "" {
		return resource.ERR_INVALID_PARAMETER
	}

	// 检查标题是否重复
	exist, err := dto.Client().Post.Query().Where(post.TitleEQ(param.Title)).Exist(c)
	if err != nil {
		logger.Error("check post title failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}
	if exist {
		return resource.ERR_INVALID_PARAMETER // 可以定义专门的错误码
	}

	// 检查别名是否重复
	if param.Slug != "" {
		exist, err = dto.Client().Post.Query().Where(post.SlugEQ(param.Slug)).Exist(c)
		if err != nil {
			logger.Error("check post slug failed", zap.Error(err))
			return resource.ERR_QUERY_FAILED
		}
		if exist {
			return resource.ERR_INVALID_PARAMETER
		}
	}

	// 创建文章
	tx, err := dto.Client().Tx(c)
	if err != nil {
		logger.Error("begin transaction failed", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}
	defer func() {
		if r := recover(); r != nil || err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Error("rollback transaction failed", zap.Error(err))
			}
		}
	}()

	// 设置默认状态
	if param.Status == "" {
		param.Status = "draft"
	}

	// 创建文章
	postCreate := tx.Post.Create().
		SetTitle(param.Title).
		SetContent(param.Content).
		SetExcerpt(param.Excerpt).
		SetSlug(param.Slug).
		SetStatus(param.Status).
		SetAuthorID(userID)

	// 设置类别
	if param.CategoryID > 0 {
		postCreate = postCreate.SetCategoryID(param.CategoryID)
	}

	// 设置发布时间
	if param.Status == "published" {
		postCreate = postCreate.SetPublishedAt(time.Now())
	}

	p, err := postCreate.Save(c)
	if err != nil {
		logger.Error("create post failed", zap.Error(err))
		if err := tx.Rollback(); err != nil {
			logger.Error("rollback transaction failed", zap.Error(err))
		}
		return resource.ERR_ADD_FAILED
	}

	// 添加标签关系
	if len(param.TagIDs) > 0 {
		for _, tagID := range param.TagIDs {
			_, err = tx.PostTagRelation.Create().
				SetPostID(p.ID).
				SetPostTagID(tagID).
				Save(c)
			if err != nil {
				logger.Error("add tag relation to post failed", zap.Error(err))
				if err := tx.Rollback(); err != nil {
					logger.Error("rollback transaction failed", zap.Error(err))
				}
				return resource.ERR_ADD_FAILED
			}
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit transaction failed", zap.Error(err))
		return resource.ERR_ADD_FAILED
	}

	return resource.CODE_SUCCESS
}

// UpdatePost 更新文章
func (s *PostService) UpdatePost(c *gin.Context, userID int, param dto.ModifyPost) resource.RspCode {
	// 获取文章
	p, err := dto.Client().Post.Query().
		Where(post.IDEQ(param.ID)).
		WithAuthor().
		Only(c)
	if err != nil {
		logger.Error("query post failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	// 检查权限（只有作者可以修改）
	if p.Edges.Author.ID != userID {
		return resource.ERR_NO_PERMISSION
	}

	// 开启事务
	tx, err := dto.Client().Tx(c)
	if err != nil {
		logger.Error("begin transaction failed", zap.Error(err))
		return resource.ERR_MOD_FAILED
	}
	defer func() {
		if r := recover(); r != nil || err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Error("rollback transaction failed", zap.Error(err))
			}
		}
	}()

	// 构建更新
	update := tx.Post.UpdateOne(p)

	if param.Title != "" && param.Title != p.Title {
		// 检查标题是否重复
		exist, err := tx.Post.Query().Where(post.TitleEQ(param.Title)).Exist(c)
		if err != nil {
			logger.Error("check post title failed", zap.Error(err))
			tx.Rollback()
			return resource.ERR_QUERY_FAILED
		}
		if exist {
			tx.Rollback()
			return resource.ERR_INVALID_PARAMETER
		}
		update = update.SetTitle(param.Title)
	}

	if param.Content != p.Content {
		update = update.SetContent(param.Content)
	}

	if param.Excerpt != p.Excerpt {
		update = update.SetExcerpt(param.Excerpt)
	}

	if param.Slug != "" && param.Slug != p.Slug {
		// 检查别名是否重复
		exist, err := tx.Post.Query().Where(post.SlugEQ(param.Slug)).Exist(c)
		if err != nil {
			logger.Error("check post slug failed", zap.Error(err))
			tx.Rollback()
			return resource.ERR_QUERY_FAILED
		}
		if exist {
			tx.Rollback()
			return resource.ERR_INVALID_PARAMETER
		}
		update = update.SetSlug(param.Slug)
	}

	if param.Status != "" && param.Status != p.Status {
		update = update.SetStatus(param.Status)
		// 如果状态变为已发布，设置发布时间
		if param.Status == "published" && p.PublishedAt.IsZero() {
			update = update.SetPublishedAt(time.Now())
		}
	}

	if param.CategoryID > 0 {
		update = update.SetCategoryID(param.CategoryID)
	}

	// 更新标签关系
	if len(param.TagIDs) > 0 {
		// 先删除现有的标签关系
		_, err = tx.PostTagRelation.Delete().
			Where(posttagrelation.PostIDEQ(param.ID)).
			Exec(c)
		if err != nil {
			logger.Error("clear tag relations failed", zap.Error(err))
			tx.Rollback()
			return resource.ERR_MOD_FAILED
		}
		
		// 添加新的标签关系
		for _, tagID := range param.TagIDs {
			_, err = tx.PostTagRelation.Create().
				SetPostID(param.ID).
				SetPostTagID(tagID).
				Save(c)
			if err != nil {
				logger.Error("add tag relation failed", zap.Error(err))
				tx.Rollback()
				return resource.ERR_MOD_FAILED
			}
		}
	}

	_, err = update.Save(c)
	if err != nil {
		logger.Error("update post failed", zap.Error(err))
		tx.Rollback()
		return resource.ERR_MOD_FAILED
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit transaction failed", zap.Error(err))
		return resource.ERR_MOD_FAILED
	}

	return resource.CODE_SUCCESS
}

// DeletePost 删除文章
func (s *PostService) DeletePost(c *gin.Context, userID int, postID int) resource.RspCode {
	// 获取文章
	p, err := dto.Client().Post.Query().
		Where(post.IDEQ(postID)).
		WithAuthor().
		Only(c)
	if err != nil {
		logger.Error("query post failed", zap.Error(err))
		return resource.ERR_QUERY_FAILED
	}

	// 检查权限（只有作者可以删除）
	if p.Edges.Author.ID != userID {
		return resource.ERR_NO_PERMISSION
	}

	// 删除文章
	err = dto.Client().Post.DeleteOneID(postID).Exec(c)
	if err != nil {
		logger.Error("delete post failed", zap.Error(err))
		return resource.ERR_DEL_FAILED
	}

	return resource.CODE_SUCCESS
}

// GetPost 获取文章详情
func (s *PostService) GetPost(c *gin.Context, userID int, postID int) (*dto.PostResponse, resource.RspCode) {
	// 获取文章
	p, err := dto.Client().Post.Query().
		Where(post.IDEQ(postID)).
		WithCategory().
		WithTagRelations().
		WithAuthor().
		Only(c)
	if err != nil {
		logger.Error("query post failed", zap.Error(err))
		return nil, resource.ERR_QUERY_FAILED
	}

	// 增加浏览次数
	_, err = dto.Client().Post.UpdateOne(p).AddViewCount(1).Save(c)
	if err != nil {
		logger.Error("update view count failed", zap.Error(err))
		// 不影响主要功能，只记录日志
	}

	// 构建响应
	response := &dto.PostResponse{
		ID:          p.ID,
		Title:       p.Title,
		Content:     p.Content,
		Excerpt:     p.Excerpt,
		Slug:        p.Slug,
		Status:      p.Status,
		ViewCount:   p.ViewCount,
		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	// 设置类别信息
	if p.Edges.Category != nil {
		response.Category = &dto.CategoryResponse{
			ID:          p.Edges.Category.ID,
			Name:        p.Edges.Category.Name,
			Slug:        p.Edges.Category.Slug,
			Description: p.Edges.Category.Description,
			Color:       p.Edges.Category.Color,
			SortOrder:   p.Edges.Category.SortOrder,
			IsActive:    p.Edges.Category.IsActive,
		}
	}

	// 设置标签信息
	if len(p.Edges.TagRelations) > 0 {
		response.Tags = make([]dto.TagResponse, len(p.Edges.TagRelations))
		for i, relation := range p.Edges.TagRelations {
			if relation.Edges.PostTag != nil {
				tag := relation.Edges.PostTag
				response.Tags[i] = dto.TagResponse{
					ID:          tag.ID,
					Name:        tag.Name,
					Slug:        tag.Slug,
					Description: tag.Description,
					Color:       tag.Color,
					PostCount:   tag.PostCount,
					IsActive:    tag.IsActive,
				}
			}
		}
	}

	// 设置作者信息
	if p.Edges.Author != nil {
		response.Author = &dto.AuthorResponse{
			ID:    p.Edges.Author.ID,
			Email: p.Edges.Author.Email,
		}
	}

	return response, resource.CODE_SUCCESS
}
