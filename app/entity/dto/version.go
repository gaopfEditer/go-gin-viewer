package dto

import (
	"time"
)

// 添加韧件版本请求
type AddFirmwareVersion struct {
	ProductID   int    `json:"product_id" binding:"required"` // 产品ID
	Version     string `json:"version" binding:"required"`    // 韧件版本号
	ReleaseDate string `json:"release_date" binding:"required"` // 发布日期
	Remark      string `json:"remark"`                        // 备注
}

// 修改韧件版本请求
type ModifyFirmwareVersion struct {
	ID          int    `json:"id" binding:"required"` // 韧件版本ID
	Version     string `json:"version,omitempty"`     // 韧件版本号
	ReleaseDate string `json:"release_date,omitempty"` // 发布日期
	Remark      string `json:"remark,omitempty"`      // 备注
}

// 添加软件版本请求
type AddSoftwareVersion struct {
	ProductID       int       `json:"product_id" binding:"required"`       // 产品ID
	Version         string    `json:"version" binding:"required"`          // 软件版本号
	ReleaseDate     string    `json:"release_date" binding:"required"`     // 发布日期
	UpdateLog       string    `json:"update_log"`                          // 更新日志
	Remark          string    `json:"remark"`                              // 备注
	FeatureIDs      []int     `json:"feature_ids"`                         // 功能ID列表
	FirmwareIDs     []int     `json:"firmware_ids"`                        // 兼容的韧件版本ID列表
}

// 修改软件版本请求
type ModifySoftwareVersion struct {
	ID              int       `json:"id" binding:"required"`          // 软件版本ID
	Version         string    `json:"version,omitempty"`              // 软件版本号
	ReleaseDate     string    `json:"release_date,omitempty"`         // 发布日期
	UpdateLog       string    `json:"update_log,omitempty"`           // 更新日志
	Remark          string    `json:"remark,omitempty"`               // 备注
	FeatureIDs      []int     `json:"feature_ids,omitempty"`          // 功能ID列表
	FirmwareIDs     []int     `json:"firmware_ids,omitempty"`         // 兼容的韧件版本ID列表
}

// 韧件版本响应
type FirmwareVersionResponse struct {
	ID             int    `json:"id"`              // 韧件版本ID
	ProductID      int    `json:"product_id"`      // 产品ID
	Version        string `json:"version"`         // 韧件版本号
	ReleaseDate    string `json:"release_date"`    // 发布日期
	CreatedBy      int    `json:"created_by"`      // 创建者ID
	CreatedByEmail string `json:"created_by_email"` // 创建者邮箱
	CreatedAt      string `json:"created_at"`      // 创建时间
	Remark         string `json:"remark"`          // 备注
}

// 软件版本响应
type SoftwareVersionResponse struct {
	ID             int                 `json:"id"`               // 软件版本ID
	ProductID      int                 `json:"product_id"`       // 产品ID
	Version        string              `json:"version"`          // 软件版本号
	ReleaseDate    string              `json:"release_date"`     // 发布日期
	CreatedBy      int                 `json:"created_by"`       // 创建者ID
	CreatedByEmail string              `json:"created_by_email"` // 创建者邮箱
	CreatedAt      string              `json:"created_at"`       // 创建时间
	UpdateLog      string              `json:"update_log"`       // 更新日志
	Remark         string              `json:"remark"`           // 备注
	Features       []FeatureInfo       `json:"features"`         // 功能列表
	Firmwares      []FirmwareInfo      `json:"firmwares"`        // 兼容的韧件版本列表
}

// 功能信息
type FeatureInfo struct {
	ID          int    `json:"id"`           // 功能ID
	FeatureName string `json:"feature_name"` // 功能名称
	FeatureCode string `json:"feature_code"` // 功能编码
}

// 韧件版本信息
type FirmwareInfo struct {
	ID          int       `json:"id"`           // 韧件版本ID
	Version     string    `json:"version"`      // 韧件版本号
	ReleaseDate time.Time `json:"release_date"` // 发布日期
} 