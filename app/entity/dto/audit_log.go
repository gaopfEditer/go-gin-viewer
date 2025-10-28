package dto

import "time"

// 定义模块常量
type AuditLogModule string

const (
	ModuleUser            AuditLogModule = "user"
	ModuleAuth            AuditLogModule = "auth"
	ModuleProduct         AuditLogModule = "product"
	ModuleFeature         AuditLogModule = "feature"
	ModuleLicenseType     AuditLogModule = "license_type"
	ModuleFirmwareVersion AuditLogModule = "firmware_version"
	ModuleSoftwareVersion AuditLogModule = "software_version"
	ModuleDevice          AuditLogModule = "device"
)

// 定义操作类型常量
type AuditLogAction string

const (
	ActionCreate   AuditLogAction = "create"
	ActionUpdate   AuditLogAction = "update"
	ActionDelete   AuditLogAction = "delete"
	ActionLogin    AuditLogAction = "login"
	ActionRegister AuditLogAction = "register"
)

type AuditLogData struct {
	UserID     int
	Action     AuditLogAction
	Module     AuditLogModule
	ProductID  int
	DetailInfo interface{}
}

// OperationLogQuery 操作日志查询参数
type OperationLogQuery struct {
	Page      int       `form:"page"`
	PageSize  int       `form:"page_size"`
	Module    string    `form:"module"`
	Operation string    `form:"operation"`
	UserID    int       `form:"user_id"`
	ProductID int       `form:"product_id"`
	StartTime time.Time `form:"start_time"`
	EndTime   time.Time `form:"end_time"`
}

// OperationLogResponse 操作日志响应
type OperationLogResponse struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	UserEmail   string    `json:"user_email"`
	Module      string    `json:"module"`
	Operation   string    `json:"operation"`
	ProductID   int       `json:"product_id,omitempty"`
	ProductName string    `json:"product_name,omitempty"`
	Detail      string    `json:"detail"`
	IPAddress   string    `json:"ip_address"`
	CreatedAt   time.Time `json:"created_at"`
}
