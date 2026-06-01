# TaskNest

TaskNest 是一个使用 Go、Gin、Gorm 和 HTTP 请求构建的 Todo 待办任务 API 项目。

这是一个适合新手学习的后端项目，重点练习：

- Gin 如何接收 HTTP 请求
- Controller 如何调用 Service
- Service 如何通过 Gorm 操作数据库
- Go 项目如何按目录分层
- RESTful API 的基本写法

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
