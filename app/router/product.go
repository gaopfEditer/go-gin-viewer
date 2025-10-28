package router

import (
	"cambridge-hit.com/gin-base/activateserver/app/controller"
	"github.com/gin-gonic/gin"
)

func init() {
	Routers = append(Routers, ProductRouterRegister)
}

func ProductRouterRegister(r *gin.RouterGroup) {
	productGroup := r.Group("product")
	productController := controller.NewProductController()
	{
		productGroup.GET("/list", productController.ListProduct)
		productGroup.POST("/add", productController.AddProduct)
		productGroup.GET("/del", productController.DeleteProduct)
		productGroup.POST("/put", productController.ModifyProduct)
		productGroup.POST("/add-manager", productController.AddManager)
	}
}
