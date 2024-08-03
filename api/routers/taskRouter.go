package routers

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/whoisaditya/golang-task-management-system/api/controllers"
	"github.com/whoisaditya/golang-task-management-system/api/middleware"
)

type Task struct {
	Name  string
	start time.Time
	end   time.Time
}

func TaskRoutes(router *gin.Engine) {
	tasks := router.Group("/task")
	{
		tasks.POST("/create", middleware.AuthMiddleware, controllers.CreateTask)
		tasks.POST("/bulkupload", middleware.AuthMiddleware, controllers.CreateTaskBulk)
		tasks.GET("/", controllers.GetTasks)
		tasks.PUT("/update/", middleware.AuthMiddleware, controllers.UpdateTask)
		tasks.DELETE("/delete/", middleware.AuthMiddleware, controllers.DeleteTask)
	}
}
