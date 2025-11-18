package dto

import "time"

// AddPost 添加文章请求参数
type AddPost struct {
	Title      string `json:"title" binding:"required"` // 文章标题
	Content    string `json:"content"`                  // 文章内容
	Excerpt    string `json:"excerpt"`                  // 文章摘要
	Slug       string `json:"slug"`                     // 文章别名
	Status     string `json:"status"`                   // 文章状态
	CategoryID int    `json:"category_id"`              // 类别ID
	TagIDs     []int  `json:"tag_ids"`                  // 标签ID列表
}

// ModifyPost 修改文章请求参数
type ModifyPost struct {
	ID         int    `json:"id" binding:"required"` // 文章ID
	Title      string `json:"title"`                 // 文章标题
	Content    string `json:"content"`               // 文章内容
	Excerpt    string `json:"excerpt"`               // 文章摘要
	Slug       string `json:"slug"`                  // 文章别名
	Status     string `json:"status"`                // 文章状态
	CategoryID int    `json:"category_id"`           // 类别ID
	TagIDs     []int  `json:"tag_ids"`               // 标签ID列表
}

// PostQuery 文章查询参数
type PostQuery struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
	Search     string `form:"search"`      // 搜索关键字
	Status     string `form:"status"`      // 文章状态
	CategoryID int    `form:"category_id"` // 类别ID
	TagID      int    `form:"tag_id"`      // 标签ID
	AuthorID   int    `form:"author_id"`   // 作者ID
}

// PostResponse 文章响应
type PostResponse struct {
	ID          int               `json:"id"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	Excerpt     string            `json:"excerpt"`
	Slug        string            `json:"slug"`
	Status      string            `json:"status"`
	ViewCount   int               `json:"view_count"`
	PublishedAt time.Time         `json:"published_at"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Category    *CategoryResponse `json:"category"`
	Tags        []TagResponse     `json:"tags"`
	Author      *AuthorResponse   `json:"author"`
}

// CategoryResponse 类别响应
type CategoryResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Color       string `json:"color"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

// TagResponse 标签响应
type TagResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Color       string `json:"color"`
	PostCount   int    `json:"post_count"`
	IsActive    bool   `json:"is_active"`
}

// AuthorResponse 作者响应
type AuthorResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// AddPostCategory 添加文章类别请求参数
type AddPostCategory struct {
	Name        string `json:"name" binding:"required"` // 类别名称
	Slug        string `json:"slug"`                    // 类别别名
	Description string `json:"description"`             // 类别描述
	Color       string `json:"color"`                   // 类别颜色
	SortOrder   int    `json:"sort_order"`              // 排序顺序
}

// ModifyPostCategory 修改文章类别请求参数
type ModifyPostCategory struct {
	ID          int    `json:"id" binding:"required"` // 类别ID
	Name        string `json:"name"`                  // 类别名称
	Slug        string `json:"slug"`                  // 类别别名
	Description string `json:"description"`           // 类别描述
	Color       string `json:"color"`                 // 类别颜色
	SortOrder   int    `json:"sort_order"`            // 排序顺序
	IsActive    bool   `json:"is_active"`             // 是否启用
}

// AddPostTag 添加文章标签请求参数
type AddPostTag struct {
	Name        string `json:"name" binding:"required"` // 标签名称
	Slug        string `json:"slug"`                    // 标签别名
	Description string `json:"description"`             // 标签描述
	Color       string `json:"color"`                   // 标签颜色
}

// ModifyPostTag 修改文章标签请求参数
type ModifyPostTag struct {
	ID          int    `json:"id" binding:"required"` // 标签ID
	Name        string `json:"name"`                  // 标签名称
	Slug        string `json:"slug"`                  // 标签别名
	Description string `json:"description"`           // 标签描述
	Color       string `json:"color"`                 // 标签颜色
	IsActive    bool   `json:"is_active"`             // 是否启用
}
