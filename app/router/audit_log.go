package router

import (
	"cambridge-hit.com/gin-base/activateserver/app/controller"
	"github.com/gin-gonic/gin"
)

func init() {
	Routers = append(Routers, AuditLogRouterRegister)
}

func AuditLogRouterRegister(r *gin.RouterGroup) {
	productGroup := r.Group("audit")
	productController := controller.NewAuditLogController()
	{
		productGroup.GET("/list", productController.ListLogs)
	}
}
