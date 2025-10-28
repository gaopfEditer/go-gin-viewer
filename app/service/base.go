package service

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"time"

	"cambridge-hit.com/gin-base/activateserver/app/entity/ent/user"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/auth"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 定义接口

type BaseService struct{}

func NewBaseService() *BaseService {
	return &BaseService{}
}

func (s *BaseService) UserRegisterByPassword(c *gin.Context, param dto.UserLoginInfo) (e resource.RspCode) {
	// 检查用户是否已存在
	exist, err := dto.Client().User.Query().
		Where(user.EmailEQ(param.Email)).
		Exist(c)
	if err != nil {
		logger.Error("查询用户失败")
		return resource.ERR_QUERY_FAILED
	}
	if exist {
		CreateAuditLog(c, nil, dto.AuditLogData{
			UserID:    dto.AnonymousID,
			Action:    dto.ActionRegister,
			Module:    dto.ModuleAuth,
			ProductID: 0,
			DetailInfo: map[string]interface{}{
				"email":  param.Email,
				"status": "failed",
				"reason": "user_exist",
			},
		})
		return resource.ERR_USER_EXIST
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("密码加密失败")
		return resource.ERR_OPERATION_FAILED
	}

	// 创建用户
	now := time.Now()
	_, err = dto.Client().User.Create().
		SetEmail(param.Email).
		SetPassword(string(hashedPassword)).
		SetCreatedAt(now).
		SetLastLoginAt(now).
		Save(c)

	if err != nil {
		logger.Error("创建用户失败")
		return resource.ERR_ADD_FAILED
	}
	CreateAuditLog(c, nil, dto.AuditLogData{
		UserID:    dto.AnonymousID,
		Action:    dto.ActionRegister,
		Module:    dto.ModuleAuth,
		ProductID: 0,
		DetailInfo: map[string]interface{}{
			"email":  param.Email,
			"status": "success",
		},
	})

	return resource.CODE_SUCCESS
}

func (s *BaseService) UserLoginByPassword(c *gin.Context, param dto.UserLoginInfo) (at, rt string, uai auth.UserAuthInfo, e resource.RspCode) {
	// 查找用户
	user_info, err := dto.Client().User.Query().
		Where(user.EmailEQ(param.Email)).
		Only(c)

	if err != nil {
		CreateAuditLog(c, nil, dto.AuditLogData{
			UserID:    dto.AnonymousID,
			Action:    dto.ActionLogin,
			Module:    dto.ModuleAuth,
			ProductID: 0,
			DetailInfo: map[string]interface{}{
				"email":  param.Email,
				"status": "failed",
				"reason": "user_not_found",
			},
		})
		e = resource.ERR_QUERY_FAILED
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user_info.Password), []byte(param.Password)); err != nil {
		CreateAuditLog(c, nil, dto.AuditLogData{
			UserID:    dto.AnonymousID,
			Action:    dto.ActionLogin,
			Module:    dto.ModuleAuth,
			ProductID: 0,
			DetailInfo: map[string]interface{}{
				"email":  param.Email,
				"status": "failed",
				"reason": "incorrect_password",
			},
		})
		e = resource.ERR_INCORRECT_PASSWORD
		return
	}

	// 创建用户认证信息(在更新用户登录时间之前创建)
	authInfo := auth.UserAuthInfo{
		UserID:      user_info.ID,
		Email:       user_info.Email,
		CreatedAt:   user_info.CreatedAt,
		LastLoginAt: user_info.LastLoginAt,
	}

	// 生成JWT token
	at, rt, err = auth.GenAccessTokenAndRefreshToken(authInfo)
	if err != nil {
		logger.Error("生成jwt失败")
		e = resource.ERR_OPERATION_FAILED
		return
	}

	_, err = dto.Client().User.
		UpdateOne(user_info).
		SetLastLoginAt(time.Now()).
		Save(c)

	CreateAuditLog(c, nil, dto.AuditLogData{
		UserID:    dto.AnonymousID,
		Action:    dto.ActionLogin,
		Module:    dto.ModuleAuth,
		ProductID: 0,
		DetailInfo: map[string]interface{}{
			"user_info": authInfo,
			"status":    "success",
		},
	})
	return at, rt, authInfo, resource.CODE_SUCCESS
}
