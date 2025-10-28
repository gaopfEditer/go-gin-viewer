package auth

import (
	"cambridge-hit.com/gin-base/activateserver/resource"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// genToken 生成token map中key=exp过期时间
func genToken(mapClaims jwt.MapClaims, signed string) (token string, err error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	token, err = claims.SignedString([]byte(signed))
	if err != nil {
		return "", err
	}
	return token, nil
}

func GenAccessTokenAndRefreshToken(info UserAuthInfo) (string, string, error) {
	// 生成accessToken
	accessToken, err := genToken(jwt.MapClaims{
		"userAuthInfo": info,
		"exp":          time.Now().Add(time.Duration(resource.Conf.JwtConfig.AccessExpire) * time.Second).Unix(),
	}, resource.Conf.JwtConfig.AccessTokenSecret)
	if err != nil {
		return "", "", err
	}

	// 生成refreshToken
	refreshToken, err := genToken(jwt.MapClaims{
		"userAuthInfo": info,
		"exp":          time.Now().Add(time.Duration(resource.Conf.JwtConfig.RefreshExpire) * time.Second).Unix(),
	}, resource.Conf.JwtConfig.RefreshTokenSecret)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func ValidateToken(tokenStr string, signed string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signed), nil
	})
	return token, err
}

// GetUserInfoFromJwt 解析token
func GetUserInfoFromJwt(token *jwt.Token) (*UserAuthInfo, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["userAuthInfo"] != nil {
			u, e := ToUserAuthInfo(claims["userAuthInfo"].(map[string]interface{}))
			if e != nil {
				return nil, e
			}
			return &u, nil
		}
	}
	return nil, errors.New("this token is not a validate access token")
}

// ParseUserInfo 解析token
func ParseUserInfo(tokenStr string) (*UserAuthInfo, error) {
	signed := resource.Conf.JwtConfig.AccessTokenSecret
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signed), nil
	})
	if err == nil {
		return GetUserInfoFromJwt(token)
	}
	return nil, errors.New("this token is not a validate access token")
}

func GetUserAuthInfo(context *gin.Context) UserAuthInfo {
	uai, exist := context.Get("userInfo")
	if !exist {
		return UserAuthInfo{
			UserID: 0,
		}
	}
	userInfo, ok := uai.(UserAuthInfo)
	if !ok || userInfo.UserID == 0 {
		return UserAuthInfo{
			UserID: 0,
		}
	}
	return userInfo
}

func GetUserAuthInfoWS(accessToken string) UserAuthInfo {
	obj, err := ValidateToken(accessToken, resource.Conf.AccessTokenSecret)
	userInfo, err1 := GetUserInfoFromJwt(obj)
	if err != nil || err1 != nil || userInfo.UserID == 0 {
		return UserAuthInfo{
			UserID: 0,
		}
	}
	return *userInfo
}
