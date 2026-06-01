# 认证、登录、Token、JWT 和权限设计学习说明

本文档用于理解：

- 认证是什么。
- 登录是什么。
- Cookie、Session、Token、JWT 是什么。
- 权限和 RBAC 是什么。
- TaskNest 后续如何增加登录和权限。

## 1. 先看结论

认证解决：

```text
你是谁？
```

权限解决：

```text
你能做什么？
```

登录成功后，后端通常会返回一个登录凭证。

前端后续请求带着这个凭证，后端就能识别用户身份。

常见登录凭证：

```text
Cookie
Session ID
Token
JWT
```

## 2. 认证是什么

认证英文是 Authentication。

它用于确认当前请求来自哪个用户。

常见认证方式：

```text
账号密码登录
手机号验证码登录
扫码登录
第三方登录
Session 登录
Token 登录
```

后端最常见流程：

```text
用户提交账号密码
后端校验账号密码
校验通过后生成登录凭证
前端保存登录凭证
后续请求携带登录凭证
后端识别用户身份
```

## 3. 登录是什么

登录是认证的一种流程。

常见步骤：

```text
1. 用户输入用户名和密码
2. 前端调用登录接口
3. 后端查询用户
4. 后端校验密码
5. 校验成功后生成 Token 或 Session
6. 返回登录凭证
7. 前端保存凭证
8. 后续请求带上凭证
```

登录接口示例：

```http
POST /api/login
Content-Type: application/json

{
  "username": "zhangsan",
  "password": "123456"
}
```

响应：

```json
{
  "token": "xxxxx.yyyyy.zzzzz"
}
```

## 4. Cookie 和 Session

Cookie 是浏览器本地保存的一小段数据。

服务端可以返回：

```http
Set-Cookie: session_id=abc123
```

浏览器后续会自动带上：

```http
Cookie: session_id=abc123
```

Session 是服务端保存的登录状态。

流程：

```text
1. 用户登录成功
2. 服务端生成 session_id
3. 服务端把用户信息保存到 Redis 或内存
4. 浏览器保存 session_id
5. 后续请求带 session_id
6. 服务端根据 session_id 找到用户
```

优点：

```text
服务端可控
可以主动让用户下线
敏感信息不暴露给客户端
```

缺点：

```text
服务端需要保存登录状态
多实例部署需要共享 Session
跨域配置较麻烦
```

## 5. Token 是什么

Token 是登录令牌。

登录成功后，后端返回 Token。

前端后续请求带上：

```http
Authorization: Bearer xxxxxx
```

Token 常用于：

```text
前后端分离系统
移动端 App
开放 API
微服务调用
```

优点：

```text
前后端分离友好
移动端使用方便
不依赖浏览器 Cookie
```

缺点：

```text
泄露后有风险
需要过期时间
需要考虑刷新 Token
主动失效较复杂
```

## 6. JWT 是什么

JWT 全称：

```text
JSON Web Token
```

JWT 是 Token 的一种常见格式。

格式：

```text
Header.Payload.Signature
```

也就是：

```text
xxxxx.yyyyy.zzzzz
```

三部分含义：

```text
Header     Token 类型和签名算法
Payload    用户信息、过期时间等
Signature  签名，防止被篡改
```

Payload 示例：

```json
{
  "user_id": 1,
  "username": "zhangsan",
  "role": "admin",
  "exp": 1719999999
}
```

注意：

```text
JWT 的 Payload 默认只是编码，不是加密。
不要放密码、手机号、身份证号等敏感信息。
```

## 7. 权限是什么

权限英文是 Authorization。

它解决：

```text
你有没有资格做这件事？
```

认证和权限不同：

```text
认证：你是谁？
权限：你能做什么？
```

例如：

```text
用户已经登录，但不能删除别人的任务。
```

这就是权限控制。

常见状态码：

```text
401 Unauthorized  未登录或登录凭证无效
403 Forbidden     已登录但没有权限
```

## 8. RBAC 是什么

RBAC 是：

```text
Role-Based Access Control
基于角色的权限控制
```

核心关系：

```text
用户 -> 角色 -> 权限
```

例如：

```text
张三 -> 管理员 -> 可以查看所有任务、删除任务
李四 -> 普通用户 -> 只能管理自己的任务
```

完整表设计可能包括：

```text
users
roles
permissions
user_roles
role_permissions
```

小项目可以先简单做：

```text
users 表加 role 字段
role = admin
role = user
```

## 9. 密码怎么保存

密码不能明文存数据库。

错误：

```text
password = 123456
```

正确：

```text
password_hash = bcrypt 哈希后的字符串
```

登录时不是解密密码，而是：

```text
用户输入密码
用 bcrypt 和数据库中的 hash 比较
匹配则登录成功
```

常见算法：

```text
bcrypt
argon2
```

不要自己发明加密算法。

## 10. Gin 认证中间件

认证逻辑不要每个接口写一遍。

应该放在中间件里。

流程：

```text
请求 -> AuthMiddleware -> Controller -> Service
```

伪代码：

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }

        userID, err := ParseToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        c.Set("user_id", userID)
        c.Next()
    }
}
```

## 11. TaskNest 如何加登录

### 新增 users 表

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    created_at DATETIME,
    updated_at DATETIME
);
```

### todos 表增加 user_id

```sql
ALTER TABLE todos ADD COLUMN user_id INTEGER;
```

查询任务时：

```sql
SELECT * FROM todos WHERE user_id = ?;
```

这样用户只能看到自己的任务。

### 新增接口

```text
POST /api/register
POST /api/login
```

需要登录的接口：

```text
GET    /api/todos
POST   /api/todos
PUT    /api/todos/:id
PATCH  /api/todos/:id/status
DELETE /api/todos/:id
```

## 12. 修改任务时校验归属

用户不能修改别人的任务。

流程：

```text
1. 从 Token 得到当前 user_id
2. 根据 todo_id 查询任务
3. 判断 todo.user_id 是否等于当前 user_id
4. 相等则允许修改
5. 不相等返回 403
```

伪代码：

```go
if todo.UserID != currentUserID {
    return ErrForbidden
}
```

## 13. 学习路线

建议顺序：

1. 理解认证和权限的区别。
2. 理解 Cookie 和 Session。
3. 理解 Token。
4. 理解 JWT。
5. 学习密码哈希。
6. 学习 Gin 中间件。
7. 给 TaskNest 增加用户表。
8. 实现注册和登录。
9. Todo 表增加 `user_id`。
10. 给任务接口加认证和权限校验。

