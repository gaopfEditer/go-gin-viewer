//// author:郭健峰
//// time:2025/02/21
//// 增改此文件下的code后请使用generate_code.go生成,例如：
//// go run -mod=mod resource/generate_code.go -input ./resource/code.go -output ./resource/code_name.go -translationDir ./resource/embed/locales
//// 输入文件路径（本文件）,生成文件名（勿直接修改），输出的翻译文件路径（勿直接修改）
//// 请按照返回码类型，编写到对应的分类，勿用iota以免code会随顺序变化

package resource

// swagger:type RspCode
type RspCode int

// 用于自动生成翻译，勿删
var files = []string{"en.json", "zh-hans.json"}

// 用于自动生成翻译，勿删
var translation = map[RspCode]string{
	CODE_SUCCESS:               "Success|成功",
	ERR_SERVER_BUSY:            "Server is busy|系统繁忙",
	ERR_OPERATION_FAILED:       "Operation failed|操作失败",
	ERR_INVALID_PARAMETER:      "Invalid parameter|参数错误",
	ERR_QUERY_FAILED:           "Query failed|查询失败",
	ERR_ADD_FAILED:             "Add failed|新增失败",
	ERR_DEL_FAILED:             "Delete failed|删除失败",
	ERR_MOD_FAILED:             "Modify failed|修改失败",
	ERR_TOKEN_EXPIRED:          "Token expired|Token已过期",
	ERR_NO_PERMISSION:          "No permission|没有权限",
	ERR_CAPTCHA_INCORRECT:      "Captcha incorrect|验证码错误",
	ERR_CAPTCHA_EXPIRED:        "Captcha expired|验证码已过期",
	ERR_EMAIL_EXIST:            "Email already exists|邮箱已存在",
	ERR_LOGIN_FAILED:           "Login failed|登录失败",
	ERR_INCORRECT_PASSWORD:     "Incorrect password|密码错误",
	ERR_USER_EXIST:             "User already exists|用户已存在",
	ERR_USER_NOT_EXIST:         "User does not exist|用户不存在",
	ERR_PRODUCT_CODE_EXIST:     "Product code exists|产品代号已存在",
	ERR_PRODUCT_NAME_EXIST:     "Product name exists|产品名称已存在",
	ERR_MANAGER_ALREADY_EXIST:  "Manager already exists|管理员已存在",
	ERR_LICENSE_TYPE_EXIST:     "License type already exists|许可证类型已存在",
	ERR_FEATURE_CODE_EXIST:     "Feature code already exists|功能编码已存在",
	ERR_FIRMWARE_VERSION_EXIST: "Firmware version already exists|韧件版本已存在",
	ERR_SOFTWARE_VERSION_EXIST: "Software version already exists|软件版本已存在",
	ERR_FIRMWARE_NOT_EXIST:     "Firmware version does not exist|韧件版本不存在",
	ERR_SOFTWARE_NOT_EXIST:     "Software version does not exist|软件版本不存在",
	ERR_LICENSE_CODE_EXIST:     "License code already exists|许可证类型已存在",
	ERR_FEATURE_NAME_EXIST:     "Feature name already exists|功能名称已存在",
	ERR_PRODUCT_HAS_RELATIONS:  "Product has associated data. Please delete all versions, license types and features first.|产品存在关联数据，请先删除所有软硬件版本、许可证类型和功能",
	ERR_ADD_LOG_FAILED:         "Add log failed|新增日志失败",
	ERR_PRODUCT_NOT_EXIST:      "Product does not exist|产品不存在",
	ERR_LICENSE_TYPE_NOT_EXIST: "License type does not exist|许可证类型不存在",
	ERR_DEVICE_SN_EXIST:        "Device SN already exists|设备序列号已存在",
	ERR_DEVICE_NOT_EXIST:       "Device does not exist|设备不存在",
}

// 系统级错误返回码，RspCode不变
const (
	CODE_SUCCESS          RspCode = 200    // 成功
	ERR_SERVER_BUSY       RspCode = 100001 // 系统繁忙
	ERR_OPERATION_FAILED  RspCode = 100002 // 操作失败
	ERR_INVALID_PARAMETER RspCode = 100003 // 参数错误
)

// 数据库错误 格式为200***
const (
	ERR_QUERY_FAILED   RspCode = 200001 + iota // 查询失败
	ERR_ADD_FAILED                             // 新增失败
	ERR_DEL_FAILED                             // 删除失败
	ERR_MOD_FAILED                             // 修改失败
	ERR_ADD_LOG_FAILED                         // 新增日志失败
)

// 用户错误 格式为201*** 具体数值不重要，以返回的消息为准
const (
	ERR_TOKEN_EXPIRED          RspCode = 201000 + iota // token过期
	ERR_NO_PERMISSION                                  // 用户无权访问
	ERR_CAPTCHA_INCORRECT                              // 验证码错误
	ERR_CAPTCHA_EXPIRED                                // 验证码过期
	ERR_EMAIL_EXIST                                    // 邮箱重复
	ERR_LOGIN_FAILED                                   // 登录失败
	ERR_INCORRECT_PASSWORD                             // 密码错误
	ERR_USER_EXIST                                     // 用户已存在
	ERR_PRODUCT_CODE_EXIST                             // 产品编号已存在
	ERR_PRODUCT_NAME_EXIST                             // 产品名已存在
	ERR_MANAGER_ALREADY_EXIST                          // 管理员已存在
	ERR_USER_NOT_EXIST                                 // 用户不存在
	ERR_FIRMWARE_VERSION_EXIST                         // 韧件版本已存在
	ERR_SOFTWARE_VERSION_EXIST                         // 软件版本已存在
	ERR_FIRMWARE_NOT_EXIST                             // 韧件版本不存在
	ERR_SOFTWARE_NOT_EXIST                             // 软件版本不存在
	ERR_LICENSE_TYPE_EXIST                             // 许可证类型已存在
	ERR_LICENSE_CODE_EXIST                             // 许可证类型已存在
	ERR_FEATURE_NAME_EXIST                             // 功能编码已存在
	ERR_FEATURE_CODE_EXIST                             // 功能编码已存在
	ERR_PRODUCT_HAS_RELATIONS                          // 产品存在关联数据
	ERR_PRODUCT_NOT_EXIST                              // 产品不存在
	ERR_LICENSE_TYPE_NOT_EXIST                         // 许可证类型不存在
	ERR_DEVICE_SN_EXIST                                // 设备序列号已存在
	ERR_DEVICE_NOT_EXIST                               // 设备不存在
)
