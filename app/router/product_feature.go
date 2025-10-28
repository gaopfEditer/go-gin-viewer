package router

import (
	"cambridge-hit.com/gin-base/activateserver/app/controller"
	"github.com/gin-gonic/gin"
)

func init() {
	Routers = append(Routers, ProductFeatureRouterRegister)
}

func ProductFeatureRouterRegister(r *gin.RouterGroup) {
	productFeatureGroup := r.Group("product-feature")
	productFeatureController := controller.NewProductFeatureController()
	{
		productFeatureGroup.GET("/list", productFeatureController.ListProductFeatures)
		productFeatureGroup.POST("/add", productFeatureController.AddProductFeature)
		productFeatureGroup.GET("/del", productFeatureController.DeleteProductFeature)
	}
} 