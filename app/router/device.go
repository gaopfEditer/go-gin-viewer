package router

import (
	"cambridge-hit.com/gin-base/activateserver/app/controller"
	"github.com/gin-gonic/gin"
)

func init() {
	Routers = append(Routers, DeviceRouterRegister)
}

func DeviceRouterRegister(r *gin.RouterGroup) {
	deviceGroup := r.Group("device")
	deviceController := controller.NewDeviceController()
	{
		// 产品列表（带设备数量统计）
		deviceGroup.GET("/products", deviceController.ListProducts)

		// 设备列表和查询
		deviceGroup.GET("/list", deviceController.ListDevices)
		deviceGroup.GET("/search", deviceController.GetDeviceBySN)

		// 设备管理
		deviceGroup.POST("/add", deviceController.AddDevice)
		deviceGroup.POST("/batch-add", deviceController.BatchAddDevices)
		deviceGroup.PUT("/update", deviceController.UpdateDevice)
		deviceGroup.DELETE("/:id", deviceController.DeleteDevice)
		deviceGroup.POST("/batch-update-license", deviceController.BatchUpdateLicenseType)

		// 获取设备激活文件
		deviceGroup.GET("/activation-file/:sn", deviceController.GetActivationFile)

		// 许可证类型列表
		deviceGroup.GET("/license-types", deviceController.GetLicenseTypes)
	}
}
