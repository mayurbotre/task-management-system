package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mayurbotre/task-management-system/internal/models"
	"github.com/mayurbotre/task-management-system/internal/service"
	"github.com/mayurbotre/task-management-system/pkg/pagination"
	"gorm.io/gorm"
)

func createTaskHandler(svc service.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			abortJSON(c, http.StatusBadRequest, err.Error())
			return
		}
		task, err := svc.CreateTask(c, service.CreateTaskInput{
			Title: req.Title, Description: req.Description, Status: req.Status, DueDate: req.DueDate,
		})
		if err != nil {
			abortJSON(c, http.StatusBadRequest, err.Error())
			return
		}
		created(c, task)
	}
}

func parseID(c *gin.Context) (uint, bool) {
	idStr := c.Param("id")
	idU64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		abortJSON(c, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return uint(idU64), true
}

func getTaskHandler(svc service.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, okID := parseID(c)
		if !okID {
			return
		}
		task, err := svc.GetTask(c, id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				abortJSON(c, http.StatusNotFound, "task not found")
				return
			}
			abortJSON(c, http.StatusInternalServerError, err.Error())
			return
		}
		ok(c, task)
	}
}

func updateTaskHandler(svc service.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, okID := parseID(c)
		if !okID {
			return
		}
		var req updateTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			abortJSON(c, http.StatusBadRequest, err.Error())
			return
		}
		task, err := svc.UpdateTask(c, id, service.UpdateTaskInput{
			Title: req.Title, Description: req.Description, Status: req.Status, DueDate: req.DueDate,
		})
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				abortJSON(c, http.StatusNotFound, "task not found")
				return
			}
			abortJSON(c, http.StatusBadRequest, err.Error())
			return
		}
		ok(c, task)
	}
}

func deleteTaskHandler(svc service.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, okID := parseID(c)
		if !okID {
			return
		}
		if err := svc.DeleteTask(c, id); err != nil {
			if err == gorm.ErrRecordNotFound {
				abortJSON(c, http.StatusNotFound, "task not found")
				return
			}
			abortJSON(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func listTasksHandler(svc service.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		var statusPtr *models.Status
		if s := c.Query("status"); s != "" {
			ss := models.Status(s)
			statusPtr = &ss
		}
		items, total, err := svc.ListTasks(c, service.ListFilter{
			Status: statusPtr, Page: page, PageSize: pageSize,
		})
		if err != nil {
			abortJSON(c, http.StatusInternalServerError, err.Error())
			return
		}
		meta := pagination.BuildMeta(page, pageSize, total)
		okWithMeta(c, items, map[string]any{
			"page": meta.Page, "pageSize": meta.PageSize, "total": meta.Total, "totalPages": meta.TotalPages,
		})
	}
}
