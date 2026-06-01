package main

import (
	"log"

	"github.com/hizys/task-nest/config"
	"github.com/hizys/task-nest/routes"
)

func main() {
	// 第一步：初始化数据库连接。
	// 这里会自动创建 tasknest.db 文件，并创建 todos 表。
	config.InitDatabase()

	// 第二步：初始化路由。
	// routes.SetupRouter 会注册所有 HTTP 接口。
	router := routes.SetupRouter()

	// 第三步：启动 HTTP 服务。
	// 访问地址是 http://localhost:8080。
	log.Println("TaskNest 服务启动成功：http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("服务启动失败：", err)
	}
}
