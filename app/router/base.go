package router

import (
	"cambridge-hit.com/gin-base/activateserver/app/controller"
	"github.com/gin-gonic/gin"
)

var Routers []func(*gin.RouterGroup)

func init() {
	Routers = append(Routers, BaseRouterRegister)
}

// 注册路由
func BaseRouterRegister(r *gin.RouterGroup) {
	baseGroup := r.Group("base")
	baseController := controller.NewBaseController()
	{
		baseGroup.GET("/getCaptcha", baseController.GetCaptcha)
		baseGroup.POST("/verifyCaptcha", baseController.VerifyCaptcha)
		baseGroup.GET("/timestamp", baseController.TimeStamp)
		//baseGroup.GET("/refresh", baseController.Refresh)
		baseGroup.POST("/login", baseController.UserLoginByPassword)
		baseGroup.POST("/register", baseController.UserRegisterByPassword)
		//baseGroup.POST("/registerByInvitation", baseController.RegisterByInvitation)
		//baseGroup.POST("/sendActivationEmail", baseController.SendActivationEmail)
		//baseGroup.POST("/validateActivationEmail", baseController.ValidateActivationEmail)
		//baseGroup.POST("/resetPassword", baseController.ResetPassword)
	}
}
