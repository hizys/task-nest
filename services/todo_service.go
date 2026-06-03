package services

import (
	"errors"
	"log"

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
	if err != nil {
		log.Printf("[ERROR] create todo failed: title=%q err=%v", req.Title, err)
		return todo, err
	}

	log.Printf("[INFO] create todo success: id=%d title=%q", todo.ID, todo.Title)
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
	if err != nil {
		log.Printf("[ERROR] list todos failed: status=%q err=%v", status, err)
		return nil, err
	}

	log.Printf("[INFO] list todos success: status=%q count=%d", status, len(todos))
	return todos, err
}

// GetTodoByID 根据 ID 查询单条任务。
func GetTodoByID(id uint) (models.Todo, error) {
	var todo models.Todo

	err := config.DB.First(&todo, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[WARN] get todo not found: id=%d", id)
		return todo, errors.New("任务不存在")
	}
	if err != nil {
		log.Printf("[ERROR] get todo failed: id=%d err=%v", id, err)
		return todo, err
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
	if err != nil {
		log.Printf("[ERROR] update todo failed: id=%d err=%v", id, err)
		return todo, err
	}

	log.Printf("[INFO] update todo success: id=%d title=%q", todo.ID, todo.Title)
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
	if err != nil {
		log.Printf("[ERROR] update todo status failed: id=%d status=%q err=%v", id, status, err)
		return todo, err
	}

	log.Printf("[INFO] update todo status success: id=%d status=%q", todo.ID, todo.Status)
	return todo, err
}

// DeleteTodo 删除任务。
func DeleteTodo(id uint) error {
	todo, err := GetTodoByID(id)
	if err != nil {
		return err
	}

	err = config.DB.Delete(&todo).Error
	if err != nil {
		log.Printf("[ERROR] delete todo failed: id=%d err=%v", id, err)
		return err
	}

	log.Printf("[INFO] delete todo success: id=%d", id)
	return nil
}

// IsValidStatus 判断任务状态是否合法。
func IsValidStatus(status string) bool {
	return status == StatusPending || status == StatusDoing || status == StatusDone
}
