package router

import (
	"cambridge-hit.com/gin-base/activateserver/app/controller"
	"github.com/gin-gonic/gin"
)

func init() {
	Routers = append(Routers, PostRouterRegister)
}

func PostRouterRegister(r *gin.RouterGroup) {
	postGroup := r.Group("post")
	postController := controller.NewPostController()
	{
		postGroup.GET("/list", postController.ListPost)
		postGroup.POST("/add", postController.AddPost)
		postGroup.POST("/update", postController.UpdatePost)
		postGroup.DELETE("/delete", postController.DeletePost)
		postGroup.GET("/detail", postController.GetPost)
	}
}
