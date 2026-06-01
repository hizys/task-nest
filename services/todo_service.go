package services

import (
	"errors"

	"gorm.io/gorm"

	"github.com/hizys/task-nest/config"
	"github.com/hizys/task-nest/models"
)

// 支持的任务状态统一放在这里，避免控制器和业务代码里到处写字符串。
const (
	StatusPending = "pending"
	StatusDoing   = "doing"
	StatusDone    = "done"
)

// CreateTodo 创建一条新的待办任务。
func CreateTodo(req models.CreateTodoRequest) (models.Todo, error) {
	todo := models.Todo{
		Title:   req.Title,
		Content: req.Content,
		Status:  StatusPending,
	}

	// Create 会执行 insert，把 todo 保存到数据库。
	err := config.DB.Create(&todo).Error
	return todo, err
}

// ListTodos 查询任务列表。
// status 为空时查询全部；status 不为空时按状态筛选。
func ListTodos(status string) ([]models.Todo, error) {
	var todos []models.Todo

	query := config.DB.Order("id desc")
	if status != "" {
		if !IsValidStatus(status) {
			return nil, errors.New("任务状态只能是 pending、doing、done")
		}
		query = query.Where("status = ?", status)
	}

	err := query.Find(&todos).Error
	return todos, err
}

// GetTodoByID 根据 ID 查询单条任务。
func GetTodoByID(id uint) (models.Todo, error) {
	var todo models.Todo

	err := config.DB.First(&todo, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return todo, errors.New("任务不存在")
	}

	return todo, err
}

// UpdateTodo 更新任务标题和内容。
func UpdateTodo(id uint, req models.UpdateTodoRequest) (models.Todo, error) {
	todo, err := GetTodoByID(id)
	if err != nil {
		return todo, err
	}

	todo.Title = req.Title
	todo.Content = req.Content

	// Save 会根据主键判断是更新还是新增；这里 todo 已有 ID，所以会执行 update。
	err = config.DB.Save(&todo).Error
	return todo, err
}

// UpdateTodoStatus 只更新任务状态。
func UpdateTodoStatus(id uint, status string) (models.Todo, error) {
	todo, err := GetTodoByID(id)
	if err != nil {
		return todo, err
	}

	if !IsValidStatus(status) {
		return todo, errors.New("任务状态只能是 pending、doing、done")
	}

	todo.Status = status
	err = config.DB.Save(&todo).Error
	return todo, err
}

// DeleteTodo 删除任务。
func DeleteTodo(id uint) error {
	todo, err := GetTodoByID(id)
	if err != nil {
		return err
	}

	return config.DB.Delete(&todo).Error
}

// IsValidStatus 判断任务状态是否合法。
func IsValidStatus(status string) bool {
	return status == StatusPending || status == StatusDoing || status == StatusDone
}
