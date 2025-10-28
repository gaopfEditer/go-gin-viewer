package dto

import "time"

// DeviceFilter 设备查询过滤条件
type DeviceFilter struct {
	ProductID     int    `json:"product_id" form:"product_id"`
	LicenseTypeID int    `json:"license_type_id" form:"license_type_id"`
	SN            string `json:"sn" form:"sn"`
	OEMTag        string `json:"oem_tag" form:"oem_tag"`
	Page          int    `json:"page" form:"page" binding:"required,min=1"`
	PageSize      int    `json:"page_size" form:"page_size" binding:"required,min=1,max=100"`
}

// DeviceAdd 添加设备请求
type DeviceAdd struct {
	ProductID     int    `json:"product_id" binding:"required"`
	SN            string `json:"sn" binding:"required"`
	LicenseTypeID int    `json:"license_type_id" binding:"required"`
	OEMTag        string `json:"oem_tag"`
	Remark        string `json:"remark"`
}

// DeviceBatchAdd 批量添加设备请求
type DeviceBatchAdd struct {
	ProductID     int      `json:"product_id" binding:"required"`
	SNs           []string `json:"sns" binding:"required"`
	LicenseTypeID int      `json:"license_type_id" binding:"required"`
	OEMTag        string   `json:"oem_tag"`
	Remark        string   `json:"remark"`
}

// DeviceUpdate 更新设备请求
type DeviceUpdate struct {
	ID            int    `json:"id" binding:"required"`
	LicenseTypeID int    `json:"license_type_id" binding:"required"`
	OEMTag        string `json:"oem_tag"`
	Remark        string `json:"remark"`
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	ID              int       `json:"id"`
	SN              string    `json:"sn"`
	ProductID       int       `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductCode     string    `json:"product_code"`
	LicenseTypeID   int       `json:"license_type_id"`
	LicenseTypeName string    `json:"license_type_name"`
	LicenseTypeCode string    `json:"license_type_code"`
	OEMTag          string    `json:"oem_tag"`
	Remark          string    `json:"remark"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       int       `json:"created_by"`
	CreatedByEmail  string    `json:"created_by_email"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       int       `json:"updated_by"`
	UpdatedByEmail  string    `json:"updated_by_email"`
}

// DeviceSummary 设备简要信息
type DeviceSummary struct {
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Count       int    `json:"count"`
}

// DeviceBatchUpdateLicense 批量更新许可证类型请求
type DeviceBatchUpdateLicense struct {
	DeviceIDs     []int  `json:"device_ids" binding:"required"`
	LicenseTypeID int    `json:"license_type_id" binding:"required"`
	Remark        string `json:"remark"`
}

// ActivationData 激活数据
type ActivationData struct {
	SN           string   `json:"sn"`            // 设备序列号
	ProductID    int      `json:"product_id"`    // 产品ID
	LicenseType  int      `json:"license_type"`  // 许可证类型ID
	OEMTag       string   `json:"oem_tag"`       // OEM标签
	CreatedAt    int64    `json:"created_at"`    // 创建时间
	FeatureCodes []string `json:"feature_codes"` // 功能编码列表
}

// ActivationFile 激活文件
type ActivationFile struct {
	Data      ActivationData `json:"data"`      // 激活数据
	Signature []byte         `json:"signature"` // RSA签名
}
