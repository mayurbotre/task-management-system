package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mayurbotre/task-management-system/internal/middleware"
	"github.com/mayurbotre/task-management-system/internal/service"
)

func SetupRouter(svc service.TaskService) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.RequestID(), middleware.Logger())

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.GET("/ready", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })

	t := r.Group("/tasks")
	{
		t.POST("", createTaskHandler(svc))
		t.GET("", listTasksHandler(svc))
		t.GET("/:id", getTaskHandler(svc))
		t.PUT("/:id", updateTaskHandler(svc))
		t.DELETE("/:id", deleteTaskHandler(svc))
	}
	return r
}
