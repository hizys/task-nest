# HTTP 和 RESTful API 基础学习说明

本文档用于理解后端接口开发最基础的内容：

- HTTP 是什么。
- HTTP 请求和响应由什么组成。
- GET、POST、PUT、PATCH、DELETE 怎么选。
- 常见 HTTP 状态码是什么意思。
- RESTful API 是什么。
- 什么是幂等性。
- TaskNest 的接口应该怎么理解。

## 1. 先看结论

后端服务本质上是在做一件事：

```text
接收请求 -> 处理业务 -> 返回响应
```

在 TaskNest 里，这个流程大概是：

```text
前端或 curl 发 HTTP 请求
   |
   v
Gin 接收请求
   |
   v
Controller 解析参数
   |
   v
Service 处理业务
   |
   v
Gorm 操作数据库
   |
   v
Gin 返回 JSON 响应
```

RESTful API 的核心思想是：

```text
URL 表示资源
HTTP Method 表示动作
HTTP Status Code 表示结果
JSON 表示传输的数据
```

例如任务资源 `todos`：

```text
GET     /api/todos       查询任务列表
GET     /api/todos/1     查询任务详情
POST    /api/todos       创建任务
PUT     /api/todos/1     完整更新任务
PATCH   /api/todos/1     局部更新任务
DELETE  /api/todos/1     删除任务
```

## 2. HTTP 是什么

HTTP 全称是：

```text
HyperText Transfer Protocol
超文本传输协议
```

可以简单理解为：

```text
HTTP 是客户端和服务端沟通的一套规则。
```

客户端可以是：

- 浏览器。
- 前端页面。
- 手机 App。
- curl。
- Postman / Apifox。
- 另一个后端服务。

服务端就是 TaskNest 这种后端 API 服务。

一次 HTTP 通信包括：

```text
客户端发送 Request
服务端返回 Response
```

## 3. HTTP 请求由什么组成

一个 HTTP 请求通常包含：

```text
Method      请求方法
Path        请求路径
Query       查询参数
Header      请求头
Body        请求体
```

### 3.1 Method

Method 表示这次请求要做什么。

```text
GET     查询
POST    新增
PUT     完整更新
PATCH   局部更新
DELETE  删除
```

### 3.2 Path

Path 表示要操作哪个资源。

例如：

```text
/api/todos
/api/todos/1
/api/users/100
```

RESTful 风格里，路径通常用名词，不用动词。

推荐：

```text
GET /api/todos
POST /api/todos
DELETE /api/todos/1
```

不推荐：

```text
GET /api/getTodos
POST /api/createTodo
POST /api/deleteTodo
```

因为动作已经由 Method 表达了。

### 3.3 Query

Query 参数通常用于筛选、分页、搜索和排序。

例如：

```text
GET /api/todos?status=done&page=1&page_size=20
```

这里：

```text
status=done
page=1
page_size=20
```

都是 Query 参数。

### 3.4 Header

Header 用来传递请求附加信息。

常见 Header：

```text
Content-Type: application/json
Authorization: Bearer xxx
Accept: application/json
User-Agent: curl
```

例如提交 JSON 时，需要：

```text
Content-Type: application/json
```

以后做登录时，常见的是：

```text
Authorization: Bearer token
```

### 3.5 Body

Body 是请求体，用来传复杂数据。

创建任务时，请求体可能是：

```json
{
  "title": "学习 HTTP",
  "content": "理解请求和响应"
}
```

Gin 里通常用 `ShouldBindJSON` 把 JSON 绑定到 Go 结构体。

## 4. HTTP 响应由什么组成

HTTP 响应通常包含：

```text
Status Code  状态码
Header       响应头
Body         响应体
```

常见响应体是 JSON。

成功示例：

```json
{
  "id": 1,
  "title": "学习 HTTP",
  "status": "pending"
}
```

失败示例：

```json
{
  "error": "todo not found"
}
```

## 5. 常见 HTTP 状态码

常见状态码：

```text
200 OK                    请求成功
201 Created               创建成功
400 Bad Request           请求参数错误
401 Unauthorized          未登录或 Token 无效
403 Forbidden             已登录但没有权限
404 Not Found             资源不存在
409 Conflict              资源冲突
429 Too Many Requests     请求太频繁
500 Internal Server Error 服务端内部错误
```

在 TaskNest 里可以这样理解：

```text
GET /api/todos/1
如果任务存在，返回 200
如果任务不存在，返回 404

POST /api/todos
如果参数正确并创建成功，返回 201 或 200
如果 title 为空，返回 400

DELETE /api/todos/1
如果没有登录，返回 401
如果删除别人的任务，返回 403
```

## 6. RESTful API 是什么

RESTful API 是一种接口设计风格。

它强调：

```text
用资源组织 URL。
用 HTTP Method 表示动作。
用状态码表达结果。
用 JSON 交换数据。
```

以 TaskNest 的 Todo 资源为例：

| Method | Path | 含义 |
| --- | --- | --- |
| GET | `/api/todos` | 查询任务列表 |
| GET | `/api/todos/:id` | 查询任务详情 |
| POST | `/api/todos` | 创建任务 |
| PUT | `/api/todos/:id` | 更新任务全部内容 |
| PATCH | `/api/todos/:id/status` | 只更新任务状态 |
| DELETE | `/api/todos/:id` | 删除任务 |

这样接口语义清楚，前后端协作也更容易。

## 7. PUT 和 PATCH 的区别

`PUT` 通常表示完整更新。

例如：

```json
{
  "title": "学习 Go",
  "content": "完整更新任务",
  "status": "doing"
}
```

`PATCH` 通常表示局部更新。

例如只修改状态：

```json
{
  "status": "done"
}
```

TaskNest 当前有：

```text
PUT   /api/todos/:id
PATCH /api/todos/:id/status
```

这个设计是合理的：

```text
PUT 用于改任务内容。
PATCH 用于只改任务状态。
```

## 8. 什么是幂等性

幂等性的意思是：

```text
同一个请求执行一次和执行多次，最终结果一样。
```

例如：

```text
DELETE /api/todos/1
```

第一次删除后任务不存在。

第二次再删，任务还是不存在。

最终结果一样，所以 DELETE 通常是幂等的。

再比如：

```text
PUT /api/todos/1
```

每次都把任务改成同样内容，最终结果一样，所以 PUT 通常是幂等的。

但是：

```text
POST /api/todos
```

每调用一次都可能创建一条新任务，所以通常不是幂等的。

幂等性在支付、订单、消息消费、重试机制里非常重要。

## 9. TaskNest 示例

创建任务：

```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"学习 HTTP","content":"理解 RESTful API"}'
```

查询任务列表：

```bash
curl http://localhost:8080/api/todos
```

按状态筛选：

```bash
curl "http://localhost:8080/api/todos?status=done"
```

修改任务：

```bash
curl -X PUT http://localhost:8080/api/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"学习 Gin","content":"继续完善 TaskNest"}'
```

修改状态：

```bash
curl -X PATCH http://localhost:8080/api/todos/1/status \
  -H "Content-Type: application/json" \
  -d '{"status":"done"}'
```

删除任务：

```bash
curl -X DELETE http://localhost:8080/api/todos/1
```

## 10. 学习路线

建议按这个顺序学：

1. 理解 Request 和 Response。
2. 理解 Method、Path、Header、Body、Query。
3. 记住常见状态码。
4. 学会用 curl / Postman / Apifox 调接口。
5. 理解 RESTful API 风格。
6. 理解幂等性。
7. 回到 Gin 代码里看路由和 Controller。

刚开始重点记住：

```text
GET 查询
POST 新增
PUT 完整更新
PATCH 局部更新
DELETE 删除
```

