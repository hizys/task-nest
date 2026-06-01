# Go/Gin 后端测试体系学习说明

本文档用于理解：

- 测试是什么。
- 为什么后端项目需要测试。
- 单元测试、接口测试、集成测试是什么。
- Mock 是什么。
- Go 项目怎么写测试。
- TaskNest 应该怎么补测试。

## 1. 先看结论

测试是用代码或工具验证程序行为是否符合预期。

它回答这些问题：

```text
代码真的能正常工作吗？
改完以后有没有把旧功能弄坏？
异常情况会不会崩？
接口返回是否符合预期？
```

Go 内置测试能力。

测试文件必须以 `_test.go` 结尾。

运行测试：

```bash
go test ./...
```

## 2. 为什么需要测试

只靠手动测试不可靠。

问题包括：

```text
容易漏场景
重复劳动多
多人协作不稳定
上线前验证成本高
改代码时心里没底
```

自动化测试可以：

- 更快发现 bug。
- 更放心重构。
- 减少重复手动验证。
- 接入 CI/CD 自动检查。
- 防止旧功能被改坏。

## 3. 单元测试是什么

单元测试测试一个很小的代码单元。

例如：

```text
一个函数
一个方法
一个 Service
一个工具方法
```

特点：

```text
范围小
速度快
失败容易定位
尽量不依赖真实外部服务
```

示例：

```go
func TestValidateStatus(t *testing.T) {
    err := ValidateStatus("done")
    if err != nil {
        t.Fatalf("expected nil, got %v", err)
    }
}
```

## 4. 表格驱动测试

Go 项目常用表格驱动测试。

适合一次测试多组输入。

```go
func TestValidateStatus(t *testing.T) {
    tests := []struct {
        name    string
        status  string
        wantErr bool
    }{
        {"pending 合法", "pending", false},
        {"doing 合法", "doing", false},
        {"done 合法", "done", false},
        {"空状态非法", "", true},
        {"未知状态非法", "xxx", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateStatus(tt.status)
            if tt.wantErr && err == nil {
                t.Fatalf("expected error")
            }
            if !tt.wantErr && err != nil {
                t.Fatalf("expected nil, got %v", err)
            }
        })
    }
}
```

## 5. 接口测试是什么

接口测试直接测试 HTTP API。

它关注：

```text
路由是否正确
请求参数是否正确绑定
状态码是否正确
响应 JSON 是否正确
错误场景是否正确返回
```

Gin 接口测试常用：

```go
net/http/httptest
```

示例：

```go
func TestHealth(t *testing.T) {
    router := gin.Default()
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    resp := httptest.NewRecorder()

    router.ServeHTTP(resp, req)

    if resp.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", resp.Code)
    }
}
```

## 6. 集成测试是什么

集成测试验证多个模块组合后是否正常。

例如：

```text
Gin Router + Controller + Service + Gorm + SQLite
```

它能发现：

```text
路由注册错误
参数绑定错误
Service 调用错误
数据库字段不匹配
迁移漏字段
```

缺点是：

```text
速度比单元测试慢
环境准备更复杂
失败定位范围更大
```

## 7. 测试数据库

测试不要直接用真实开发数据库。

推荐：

```text
SQLite 内存数据库
独立测试数据库
每个测试初始化自己的数据
测试结束后自动清理
```

SQLite 内存数据库示例：

```go
db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
if err != nil {
    t.Fatalf("failed to open test db: %v", err)
}

if err := db.AutoMigrate(&models.Todo{}); err != nil {
    t.Fatalf("failed to migrate test db: %v", err)
}
```

## 8. Mock 是什么

Mock 是模拟对象。

如果测试依赖外部服务，可以用 Mock 替代。

例如：

```text
真实数据库
真实 Redis
真实短信服务
真实邮件服务
真实第三方 API
```

测试时可以模拟：

```text
发送短信成功
发送短信失败
Redis 超时
第三方服务返回错误
```

Mock 的好处：

```text
测试更快
测试更稳定
不依赖外部服务
更容易覆盖异常场景
```

新手可以先手写 Fake，不急着上复杂 Mock 框架。

## 9. 测试应该覆盖什么

正常场景：

```text
创建任务成功
查询任务成功
修改任务成功
删除任务成功
```

参数错误：

```text
title 为空
id 不是数字
JSON 格式错误
status 非法
```

数据不存在：

```text
查询不存在任务
修改不存在任务
删除不存在任务
```

权限错误：

```text
未登录访问
Token 过期
用户修改别人的任务
```

外部依赖异常：

```text
数据库失败
Redis 超时
MQ 发送失败
```

## 10. Go 测试常用命令

运行所有测试：

```bash
go test ./...
```

详细输出：

```bash
go test -v ./...
```

只测试某个包：

```bash
go test ./services
```

只运行某个测试：

```bash
go test -run TestCreateTodo ./services
```

查看覆盖率：

```bash
go test -cover ./...
```

生成覆盖率文件：

```bash
go test -coverprofile=coverage.out ./...
```

浏览器查看覆盖率：

```bash
go tool cover -html=coverage.out
```

注意：

```text
覆盖率高不等于测试质量一定高。
测试是否真正验证了业务行为更重要。
```

## 11. TaskNest 怎么补测试

Service 层：

```text
CreateTodo 正常创建
CreateTodo 标题为空返回错误
ListTodos 按状态筛选
GetTodo 查询不存在返回错误
UpdateTodo 修改内容
UpdateTodoStatus 状态非法返回错误
DeleteTodo 删除后查询不到
```

Controller 层：

```text
GET /health 返回 200
POST /api/todos 参数正确返回成功
POST /api/todos 缺 title 返回 400
GET /api/todos?status=xxx 返回 400
GET /api/todos/:id 不存在返回 404
```

集成测试：

```text
初始化测试数据库
注册路由
发 HTTP 请求
检查响应
检查数据库数据
```

未来登录后还要测：

```text
未登录返回 401
无效 Token 返回 401
用户只能看自己的任务
用户不能删别人的任务
管理员可以管理所有任务
```

## 12. 学习路线

建议顺序：

1. 学会 `_test.go` 和 `go test ./...`。
2. 给简单函数写测试。
3. 学表格驱动测试。
4. 给 Service 写单元测试。
5. 用 `httptest` 写 Gin 接口测试。
6. 用 SQLite 内存数据库写集成测试。
7. 学 Mock 和 Fake。
8. 把 `go test ./...` 接入 CI/CD。

