package config

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/hizys/task-nest/models"
)

// DB 是整个项目共用的数据库连接对象。
// 其他包如果需要操作数据库，可以通过 config.DB 拿到它。
var DB *gorm.DB

// InitDatabase 初始化数据库连接，并自动创建 todos 表。
func InitDatabase() {
	var err error

	// 这里使用 SQLite，适合入门学习：不需要额外安装 MySQL，启动项目就能生成 tasknest.db 文件。
	DB, err = gorm.Open(sqlite.Open("tasknest.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败：", err)
	}

	// AutoMigrate 会根据模型结构自动创建或更新表结构。
	// 注意：它适合学习和开发环境，正式生产环境通常会使用迁移脚本管理表结构。
	err = DB.AutoMigrate(&models.Todo{})
	if err != nil {
		log.Fatal("数据库表迁移失败：", err)
	}
}
