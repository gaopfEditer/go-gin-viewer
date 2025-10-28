package router

import (
	"cambridge-hit.com/gin-base/activateserver/app/controller"
	"github.com/gin-gonic/gin"
)

func init() {
	Routers = append(Routers, VersionRouterRegister)
}

// VersionRouterRegister 注册版本管理相关路由
func VersionRouterRegister(r *gin.RouterGroup) {
	versionController := controller.NewVersionController()
	firmwareVersionGroup := r.Group("firmware-versions")
	{
		// 韧件版本管理路由
		firmwareVersionGroup.GET("/:product_id", versionController.ListFirmwareVersions)
		firmwareVersionGroup.POST("/", versionController.AddFirmwareVersion)
		firmwareVersionGroup.PUT("/", versionController.ModifyFirmwareVersion)
		firmwareVersionGroup.DELETE("/:id", versionController.DeleteFirmwareVersion)
	}
	softwareVersionGroup := r.Group("software-versions")
	{
		// 软件版本管理路由
		softwareVersionGroup.GET("/:product_id", versionController.ListSoftwareVersions)
		softwareVersionGroup.POST("/", versionController.AddSoftwareVersion)
		softwareVersionGroup.PUT("/", versionController.ModifySoftwareVersion)
		softwareVersionGroup.DELETE("/:id", versionController.DeleteSoftwareVersion)
	}
	
	// 添加其他版本相关路由
	versionGroup := r.Group("version")
	{
		versionGroup.GET("/firmware/all/:product_id", versionController.GetProductFirmwareVersions)
		versionGroup.GET("/features/:product_id", versionController.GetProductFeatures)
	}
}
