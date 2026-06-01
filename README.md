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

配套说明文档建议按下面的阶段阅读。同一个阶段里的文档可以并行看，不需要完全按一篇一篇死记。

### 第一阶段：先理解项目怎么被开发出来

- [软件项目交付流程说明](docs/software-delivery-flow.md)
- [软件项目交付流程图](docs/assets/software-delivery-flow.png)
- [产品文档和技术设计文档学习说明](docs/product-and-technical-design-docs-guide.md)
- [Git 分支协作策略说明](docs/git-branching-strategies.md)
- [Go Modules 依赖管理说明](docs/go-modules-guide.md)

### 第二阶段：掌握后端接口和项目基础

- [HTTP 和 RESTful API 基础学习说明](docs/http-restful-api-guide.md)
- [Gin 和 Gorm 框架学习说明](docs/gin-gorm-guide.md)
- [数据库、SQL、事务、索引和慢 SQL 优化学习说明](docs/database-sql-transaction-index-guide.md)

### 第三阶段：补齐常见业务能力

- [认证、登录、Token、JWT 和权限设计学习说明](docs/auth-jwt-permission-guide.md)
- [缓存和 Redis 学习说明](docs/cache-redis-guide.md)
- [消息队列 MQ 学习说明](docs/message-queue-guide.md)
- [安全基础学习说明](docs/security-basics-guide.md)

### 第四阶段：学会测试、排查和观测服务

- [Go/Gin 后端测试体系学习说明](docs/testing-guide.md)
- [日志、Trace、PProf 和 ELK 学习说明](docs/logging-tracing-pprof-elk-guide.md)
- [监控和告警学习说明](docs/monitoring-alerting-guide.md)

### 第五阶段：学习部署和自动化交付

- [部署、Docker、K8s、Jenkins 和自动化流水线学习说明](docs/deployment-docker-k8s-jenkins-guide.md)
- [CI/CD 深入学习说明](docs/cicd-advanced-guide.md)

### 第六阶段：理解更大的系统架构

- [系统设计基础学习说明](docs/system-design-basics-guide.md)
- [微服务、注册中心、配置中心、网关、Nacos 和 K8s 学习说明](docs/microservices-registry-config-gateway-nacos-k8s-guide.md)

### 第七阶段：提高 AI 辅助开发效率

- [AI 编程助手核心概念学习说明](docs/ai-agent-concepts-guide.md)

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
