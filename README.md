# TaskNest

TaskNest 是一个使用 Go、Gin、Gorm 和 HTTP 请求构建的 Todo 待办任务 API 项目。

这是一个适合新手学习的后端项目，重点练习：

- Gin 如何接收 HTTP 请求
- Controller 如何调用 Service
- Service 如何通过 Gorm 操作数据库
- Go 项目如何按目录分层
- RESTful API 的基本写法

## 学习目标

我希望通过 TaskNest 学习一个软件项目从想法到上线的完整过程，而不只是写几行接口代码。

本项目的学习范围包括：

- **技术设计**：学习如何从需求出发，设计接口、数据库表、项目结构和业务流程。
- **软件开发**：学习如何使用 Go、Gin、Gorm 编写一个可运行的后端 API 项目。
- **软件调试**：学习如何通过日志、接口请求、数据库数据和错误信息定位问题。
- **部署上线**：学习如何打包程序、准备运行环境、启动服务，并理解灰度和上线流程。
- **软件项目交付流程**：学习需求评审、技术设计、技术评审、用例评审、技术连调、提测、灰度、上线这一整套流程。

配套说明文档：

- [软件项目交付流程说明](docs/software-delivery-flow.md)
- [软件项目交付流程图](docs/assets/software-delivery-flow.png)
- [Git 分支协作策略说明](docs/git-branching-strategies.md)

## 项目结构

```text
task-nest/
├── main.go                  # 项目入口，启动 HTTP 服务
├── config/
│   └── database.go          # 数据库连接和表结构迁移
├── controllers/
│   └── todo_controller.go   # HTTP 请求处理层
├── models/
│   └── todo.go              # 数据模型和请求参数结构体
├── routes/
│   └── routes.go            # 路由注册
└── services/
    └── todo_service.go      # 业务逻辑和数据库操作
```

## 启动项目

```bash
go run .
```

启动成功后访问：

```text
http://localhost:8080
```

项目默认使用 SQLite，启动后会在本地生成 `tasknest.db` 数据库文件。

## 接口列表

### 健康检查

```bash
curl http://localhost:8080/health
```

### 创建任务

```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"学习 Gin 和 Gorm","content":"完成 TaskNest 第一个接口"}'
```

### 查询任务列表

```bash
curl http://localhost:8080/api/todos
```

按状态筛选：

```bash
curl "http://localhost:8080/api/todos?status=done"
```

状态只能是：

```text
pending
doing
done
```

### 查询任务详情

```bash
curl http://localhost:8080/api/todos/1
```

### 修改任务内容

```bash
curl -X PUT http://localhost:8080/api/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"继续学习 Go","content":"理解 Gin Controller 和 Gorm 操作"}'
```

### 修改任务状态

```bash
curl -X PATCH http://localhost:8080/api/todos/1/status \
  -H "Content-Type: application/json" \
  -d '{"status":"done"}'
```

### 删除任务

```bash
curl -X DELETE http://localhost:8080/api/todos/1
```

## 学习流程

建议按这个顺序看代码：

1. `main.go`：看项目怎么启动
2. `routes/routes.go`：看 URL 和方法绑定到哪个函数
3. `controllers/todo_controller.go`：看 HTTP 请求怎么变成 Go 结构体
4. `services/todo_service.go`：看业务逻辑怎么操作数据库
5. `models/todo.go`：看数据库表字段怎么定义
6. `config/database.go`：看数据库怎么连接和自动建表
