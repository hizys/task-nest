package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/hizys/task-nest/models"
	"github.com/hizys/task-nest/services"
)

// response 是统一的接口响应格式。
// 这样前端或调用方看到的 JSON 结构会更稳定。
type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CreateTodoHandler 处理 POST /api/todos 请求。
func CreateTodoHandler(c *gin.Context) {
	var req models.CreateTodoRequest

	// ShouldBindJSON 会把请求体里的 JSON 绑定到 req 结构体。
	// 如果 title 没传，因为 binding:"required"，这里会返回错误。
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response{Code: 400, Message: "参数错误：title 不能为空"})
		return
	}

	todo, err := services.CreateTodo(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Code: 500, Message: "创建任务失败"})
		return
	}

	c.JSON(http.StatusCreated, response{Code: 201, Message: "创建任务成功", Data: todo})
}

// ListTodosHandler 处理 GET /api/todos 请求。
func ListTodosHandler(c *gin.Context) {
	// Query 用来读取 URL 查询参数，例如 /api/todos?status=done。
	status := c.Query("status")

	todos, err := services.ListTodos(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response{Code: 200, Message: "查询任务列表成功", Data: todos})
}

// GetTodoHandler 处理 GET /api/todos/:id 请求。
func GetTodoHandler(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	todo, err := services.GetTodoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response{Code: 404, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response{Code: 200, Message: "查询任务详情成功", Data: todo})
}

// UpdateTodoHandler 处理 PUT /api/todos/:id 请求。
func UpdateTodoHandler(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req models.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response{Code: 400, Message: "参数错误：title 不能为空"})
		return
	}

	todo, err := services.UpdateTodo(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, response{Code: 404, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response{Code: 200, Message: "更新任务成功", Data: todo})
}

// UpdateTodoStatusHandler 处理 PATCH /api/todos/:id/status 请求。
func UpdateTodoStatusHandler(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req models.UpdateTodoStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response{Code: 400, Message: "参数错误：status 不能为空"})
		return
	}

	todo, err := services.UpdateTodoStatus(id, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response{Code: 200, Message: "更新任务状态成功", Data: todo})
}

// DeleteTodoHandler 处理 DELETE /api/todos/:id 请求。
func DeleteTodoHandler(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	err := services.DeleteTodo(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response{Code: 404, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response{Code: 200, Message: "删除任务成功"})
}

// parseID 把路径参数里的 id 字符串转换成 uint。
// 返回 ok=false 时，表示 id 不合法，并且函数内部已经返回了错误响应。
func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, response{Code: 400, Message: "id 必须是大于 0 的数字"})
		return 0, false
	}

	return uint(id), true
}
