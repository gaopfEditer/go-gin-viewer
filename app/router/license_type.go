package router

import (
	"cambridge-hit.com/gin-base/activateserver/app/controller"
	"github.com/gin-gonic/gin"
)

func init() {
	Routers = append(Routers, LicenseTypeRouterRegister)
}

func LicenseTypeRouterRegister(r *gin.RouterGroup) {
	licenseTypeGroup := r.Group("license-type")
	licenseTypeController := controller.NewLicenseTypeController()
	{
		licenseTypeGroup.GET("/list", licenseTypeController.ListLicenseTypes)
		licenseTypeGroup.POST("/add", licenseTypeController.AddLicenseType)
		licenseTypeGroup.GET("/del", licenseTypeController.DeleteLicenseType)
		licenseTypeGroup.POST("/update-features", licenseTypeController.UpdateLicenseTypeFeatures)
	}
}
