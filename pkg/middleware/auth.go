package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"strings"

	"cambridge-hit.com/gin-base/activateserver/pkg/util/auth"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/req-resp/resp"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
)

// noAuthRouters 不需要认证的路由
var noAuthRouters = []string{
	"/base",
	"/device/activation-file",
}

// JwtAuth jwt认证
// 解析请求头中的token
// 解析成功则通过，失败返回错误信息
func JwtAuth(prefix string) gin.HandlerFunc {
	return func(context *gin.Context) {
		// 前端页面无需鉴权
		if !strings.HasPrefix(context.Request.URL.Path, resource.Conf.App.ApiPrefix) {
			context.Next()
			return
		}

		// trim掉prefix
		test, _ := strings.CutPrefix(context.Request.URL.Path, prefix)
		for _, v := range noAuthRouters {
			if strings.HasPrefix(test, v) {
				context.Next()
				return
			}
		}

		// 从 Authorization 请求头中获取 Bearer Token
		authHeader := context.GetHeader("Authorization")
		if authHeader == "" {
			resp.Error(context, resource.ERR_TOKEN_EXPIRED)
			context.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			resp.Error(context, resource.ERR_TOKEN_EXPIRED)
			context.Abort()
			return
		}

		accessToken := tokenParts[1]

		// 尝试从 Cookie 中获取 refresh_token
		cookie := context.Request.Cookies()
		var refreshToken string
		for _, c := range cookie {
			if c.Name == "refresh_token" {
				refreshToken = c.Value
				break
			}
		}

		// 验证 access_token
		accessTokenObj, err := auth.ValidateToken(accessToken, resource.Conf.JwtConfig.AccessTokenSecret)
		if err == nil && accessTokenObj.Valid {
			//at有效
			setAuthInfo(context, accessTokenObj)
			context.Next()
			return
		}

		// access_token 无效，检查 refresh_token
		refreshTokenObj, err := auth.ValidateToken(refreshToken, resource.Conf.JwtConfig.RefreshTokenSecret)
		if err == nil && refreshTokenObj.Valid {
			//rt有效
			userAuthInfo := setAuthInfo(context, refreshTokenObj)
			// 刷新 access_token
			newAccessToken, newRefreshToken, err := auth.GenAccessTokenAndRefreshToken(*userAuthInfo)
			if err != nil {
				resp.Error(context, resource.ERR_OPERATION_FAILED)
				context.Abort()
				return
			}

			// 更新 context
			context.Header("Authorization", "Bearer "+newAccessToken)
			if newRefreshToken != "" { // 如果 refresh_token 也更新了
				expire := int(resource.Conf.JwtConfig.RefreshExpire)
				domain := strings.Split(context.Request.Host, ":")[0] // 获取主机名，不包含端口号
				//context.SetCookie("refresh_token", newRefreshToken,
				//	expire, "/", domain, true, true)
				context.SetCookie("refresh_token", newRefreshToken,
					expire, "/", domain, false, true)
			}
			context.Next()
			return
		}
		// 均无效
		context.Abort()
	}
}

func setAuthInfo(c *gin.Context, token *jwt.Token) *auth.UserAuthInfo {
	// 获取用户信息体
	userAuthInfo, _ := auth.GetUserInfoFromJwt(token)
	c.Set("userInfo", *userAuthInfo) // 更新用户信息
	return userAuthInfo
}
