package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hizys/task-nest/controllers"
)

// SetupRouter 负责创建 Gin 引擎，并集中注册所有路由。
func SetupRouter() *gin.Engine {
	// gin.Default() 会默认带上日志中间件和错误恢复中间件。
	router := gin.Default()

	// 健康检查接口，用来快速判断服务是否启动成功。
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "TaskNest is running",
		})
	})

	// /api 是统一的接口前缀，方便后续扩展更多模块。
	api := router.Group("/api")
	{
		// todos 是任务模块的路由组。
		todos := api.Group("/todos")
		{
			todos.POST("", controllers.CreateTodoHandler)
			todos.GET("", controllers.ListTodosHandler)
			todos.GET("/:id", controllers.GetTodoHandler)
			todos.PUT("/:id", controllers.UpdateTodoHandler)
			todos.DELETE("/:id", controllers.DeleteTodoHandler)
			todos.PATCH("/:id/status", controllers.UpdateTodoStatusHandler)
		}
	}

	return router
}
