package router

import (
	"cambridge-hit.com/gin-base/activateserver/app/controller"
	"github.com/gin-gonic/gin"
)

func init() {
	Routers = append(Routers, MetricsRouterRegister)
}

func MetricsRouterRegister(r *gin.RouterGroup) {
	group := r.Group("metrics")
	ctl := controller.NewMetricsController()
	{
		group.GET("/list", ctl.ListMetrics)
		group.POST("/add", ctl.AddMetric)
		group.POST("/update", ctl.UpdateMetric)
		group.DELETE("/delete", ctl.DeleteMetric)
		group.GET("/detail", ctl.GetMetric)
	}
}

