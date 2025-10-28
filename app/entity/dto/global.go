package dto

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent"
	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/productmanager"
)

var client *ent.Client

func Client() *ent.Client {
	return client
}

func SetClient(c *ent.Client) {
	if client == nil {
		client = c
	}
}

const (
	AnonymousID  = 100000000
	SuperAdminID = 1
)

type UserLoginInfo struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AddProduct struct {
	//ID          int    `json:"id,omitempty"`                    //指定产品id
	Code        string `json:"code" binding:"required"`         // 产品代号
	ProductName string `json:"product_name" binding:"required"` // 产品名称
	ProductType string `json:"product_type" binding:"required"` // 产品名称
}

type productAssistant struct {
	UserID     int                        `json:"user_id"`
	Permission productmanager.Permissions `json:"permission"`
	Remark     string                     `json:"remark"`
}

type ModifyProduct struct { //also used in add_product
	ID int `json:"id" binding:"required"` // 产品id不可空
	//Code        string `json:"code" binding:"required"`         // 产品代号不允许修改
	ProductName      string             `json:"product_name,omitempty"`      // 产品名称
	ProductType      string             `json:"product_type,omitempty"`      // 产品类别
	ManagerMain      int                `json:"manager_main,omitempty"`      // 主管理员
	ManagerAssistant []productAssistant `json:"manager_assistant,omitempty"` // 副管理员
}

// AddManager 添加产品管理员请求参数
type AddManager struct {
	ProductID   int                        `json:"product_id" binding:"required"`  // 产品ID
	Email       string                     `json:"email" binding:"required,email"` // 用户邮箱
	Permissions productmanager.Permissions `json:"permissions"`                    // 权限 (read/full)
	Remark      string                     `json:"remark"`                         // 备注
}

// AddLicenseType 添加许可证类型请求参数
type AddLicenseType struct {
	ProductID   int    `json:"product_id" binding:"required"`   // 产品ID
	TypeName    string `json:"type_name" binding:"required"`    // 许可证类型名称
	LicenseType string `json:"license_type" binding:"required"` // 许可证编码
	FeatureIDs  []int  `json:"feature_ids"`                     // 功能ID列表
}

// AddProductFeature 添加产品功能请求参数
type AddProductFeature struct {
	ProductID   int    `json:"product_id" binding:"required"`   // 产品ID
	FeatureName string `json:"feature_name" binding:"required"` // 功能名称
	FeatureCode string `json:"feature_code" binding:"required"` // 功能编码
}

// UpdateLicenseTypeFeatures 更新许可证类型功能列表请求参数
type UpdateLicenseTypeFeatures struct {
	TypeID     int   `json:"type_id" binding:"required"` // 许可证类型ID
	FeatureIDs []int `json:"feature_ids"`                // 功能ID列表
}

// PageParams 分页参数
type PageParams struct {
	Page     int `form:"page" json:"page"`           // 页码，从1开始
	PageSize int `form:"page_size" json:"page_size"` // 每页数量
}

// PageResult 分页结果
type PageResult struct {
	Total    int64                  `json:"total"`     // 总记录数
	Page     int                    `json:"page"`      // 当前页码
	PageSize int                    `json:"page_size"` // 每页数量
	List     interface{}            `json:"list"`      // 数据列表
	Extra    map[string]interface{} `json:"extra,omitempty"` // 额外数据
}

//// OperationLogQuery 操作日志查询参数
//type OperationLogQuery struct {
//	PageParams
//	Module    string    `form:"module" json:"module"`       // 模块名称
//	Operation string    `form:"operation" json:"operation"` // 操作类型
//	UserID    int       `form:"userId" json:"userId"`       // 操作用户ID
//	StartTime time.Time `form:"startTime" json:"startTime"` // 开始时间
//	EndTime   time.Time `form:"endTime" json:"endTime"`     // 结束时间
//	ProductID int       `form:"productId" json:"productId"` // 产品ID
//}
//
//// OperationLogResponse 操作日志响应数据
//type OperationLogResponse struct {
//	ID          int       `json:"id"`                             // 日志ID
//	UserID      int       `json:"userId"`                         // 操作者ID
//	UserEmail   string    `json:"userEmail"`                      // 操作者邮箱
//	Module      string    `json:"module"`                         // 操作模块
//	Operation   string    `json:"operation"`                      // 具体操作
//	Detail      string    `json:"detail"`                         // 详细信息
//	IPAddress   string    `json:"ipAddress"`                      // IP地址
//	CreatedAt   time.Time `json:"createdAt"`                      // 操作时间
//	ProductID   int       `form:"productId" json:"productId"`     // 产品ID
//	ProductName string    `form:"productName" json:"productName"` // 产品ID
//}
