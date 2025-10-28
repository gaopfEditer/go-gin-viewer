package controller

import (
	"cambridge-hit.com/gin-base/activateserver/app/entity/dto"
	"strings"
	"time"

	"cambridge-hit.com/gin-base/activateserver/app/service"
	"cambridge-hit.com/gin-base/activateserver/pkg/util"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/req-resp/resp"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
)

type BaseController struct {
	s *service.BaseService
}

func NewBaseController() *BaseController {
	return &BaseController{s: service.NewBaseService()}
}

func (cl *BaseController) TimeStamp(c *gin.Context) {
	// 获取服务器时间
	now := time.Now()
	unix := now.Unix()
	resp.Success(c, map[string]interface{}{
		"Timestamp": unix,
	})
}

// GetCaptcha 生成图片验证码
func (cl *BaseController) GetCaptcha(c *gin.Context) {
	// 生成图片验证码，并返回给客户端
	captcha, err := util.GenerateCaptcha()
	if err != nil {
		resp.Error(c, resource.ERR_SERVER_BUSY)
		return
	}
	resp.Success(c, map[string]interface{}{
		"ID":    captcha.Id,
		"Image": captcha.Image,
	})
}

// VerifyCaptcha 校验图片验证码
func (cl *BaseController) VerifyCaptcha(c *gin.Context) {
	var captchaResult util.CaptchaResult
	err := c.ShouldBind(&captchaResult) // 等价于：context.ShouldBind(&captchaResult)
	if err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}
	// 获取图片验证码的id和b64s,并校验
	result := util.VerifyCaptcha(captchaResult.Id, captchaResult.VerifyValue)
	if !result {
		resp.Error(c, resource.ERR_CAPTCHA_INCORRECT)
		return
	}

	resp.Success(c, resource.CODE_SUCCESS)
	//人机验证码验证通过，十分钟内可调用发送邮件验证码，value为上次发送短信验证码时间，先记录为零值
	//cache.MyRedis.Set(str.GetCaptchaExpireKey(captchaResult.Id), time.Unix(0, 0), time.Minute*10)
}

//func (cl *BaseController) Refresh(c *gin.Context) {
//	token := c.GetHeader("RefreshToken")
//	tokenObj, err := jwt.ValidateToken(token, resource.Conf.JwtConfig.RefreshTokenSecret)
//	if err != nil || !tokenObj.Valid {
//		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
//		return
//	}
//	ua, err := jwt.GetUserInfoFromJwt(tokenObj)
//	if err != nil || ua.UserID == 0 {
//		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
//		return
//	}
//	at, rt, err := jwt.GenAccessTokenAndRefreshToken(*ua)
//	if err != nil {
//		resp.Error(c, resource.ERR_TOKEN_EXPIRED)
//		return
//	}
//	resp.Success(c, map[string]interface{}{
//		"AccessToken":  at,
//		"RefreshToken": rt,
//	})
//}

// UserRegisterByPassword
// @Tags     Base
// @Summary  用户注册
// @Produce   application/json
// @Param    data  body      dto.UserLoginInfo   true  "参数：用户注册"
// @Success  200   {object}  resp.Response{message=string}  "用户注册"
// @Router   /activate/base/register [post]
func (cl *BaseController) UserRegisterByPassword(c *gin.Context) {
	var param dto.UserLoginInfo
	err := c.ShouldBindBodyWithJSON(&param)
	if err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}

	e := cl.s.UserRegisterByPassword(c, param)
	if e != resource.CODE_SUCCESS {
		resp.Error(c, e)
		return
	}

	resp.Success(c)
}

// UserLoginByPassword
// @Tags     Base
// @Summary  用户登录
// @Produce   application/json
// @Param    data  body      dto.UserLoginInfo   true  "参数：用户登录"
// @Success  200   {object}  resp.Response{message=string}  "用户登录"
// @Router   /activate/base/login [post]
func (cl *BaseController) UserLoginByPassword(c *gin.Context) {
	var param dto.UserLoginInfo
	err := c.ShouldBindBodyWithJSON(&param)
	if err != nil {
		resp.Error(c, resource.ERR_INVALID_PARAMETER)
		return
	}
	at, rt, uai, e := cl.s.UserLoginByPassword(c, param)
	if e != resource.CODE_SUCCESS {
		resp.Error(c, e)
		return
	}

	expire := int(resource.Conf.JwtConfig.RefreshExpire)
	c.Header("Authorization", "Bearer "+at)                            // 设置 Authorization 请求头
	domain := strings.Split(c.Request.Host, ":")[0]                    // 获取主机名，不包含端口号
	c.SetCookie("refresh_token", rt, expire, "/", domain, false, true) // 设置 Cookie
	resp.Success(c, uai)
}
