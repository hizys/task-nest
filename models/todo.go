package models

import "time"

// Todo 表示一条待办任务。
// Gorm 会根据这个结构体生成 todos 表。
type Todo struct {
	// ID 是主键，Gorm 看到 primaryKey 后会把它作为表的唯一标识。
	ID uint `json:"id" gorm:"primaryKey"`

	// Title 是任务标题，not null 表示不能为空。
	Title string `json:"title" gorm:"type:varchar(100);not null"`

	// Content 是任务详细内容，适合放更长的说明。
	Content string `json:"content" gorm:"type:text"`

	// Status 是任务状态，例如 pending、doing、done。
	Status string `json:"status" gorm:"type:varchar(20);not null;default:pending"`

	// CreatedAt 和 UpdatedAt 是 Gorm 的约定字段。
	// 创建和更新数据时，Gorm 会自动维护这两个时间。
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateTodoRequest 表示创建任务时，前端或接口调用方需要传入的数据。
type CreateTodoRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

// UpdateTodoRequest 表示更新任务内容时允许修改的数据。
type UpdateTodoRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

// UpdateTodoStatusRequest 表示只更新任务状态时传入的数据。
type UpdateTodoStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
