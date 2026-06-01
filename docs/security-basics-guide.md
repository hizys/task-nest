# 安全基础学习说明

本文档用于理解：

- SQL 注入是什么。
- XSS 是什么。
- CSRF 是什么。
- HTTPS / TLS 是什么。
- 密钥管理是什么。
- 敏感信息脱敏是什么。
- TaskNest 应该注意哪些安全问题。

## 1. 先看结论

后端安全的目标是保护：

```text
系统
数据
用户
密钥
权限
```

常见风险：

```text
数据库被注入攻击
用户 Token 被盗
密码或密钥泄露
接口被越权调用
日志泄露敏感信息
HTTP 明文传输被监听
```

新手先重点理解：

```text
SQL 注入
XSS
CSRF
HTTPS/TLS
密钥管理
敏感信息脱敏
```

## 2. SQL 注入是什么

SQL 注入是攻击者通过输入特殊内容，改变 SQL 原本含义。

错误示例：

```go
sql := "SELECT * FROM todos WHERE title = '" + title + "'"
```

正常输入：

```text
学习 Go
```

SQL：

```sql
SELECT * FROM todos WHERE title = '学习 Go'
```

恶意输入：

```text
' OR '1'='1
```

SQL 可能变成：

```sql
SELECT * FROM todos WHERE title = '' OR '1'='1'
```

条件永远成立，可能查出所有数据。

## 3. 如何防 SQL 注入

不要拼接 SQL。

错误：

```go
db.Raw("SELECT * FROM todos WHERE title = '" + title + "'").Scan(&todos)
```

正确：

```go
db.Where("title = ?", title).Find(&todos)
```

这里 `?` 是占位符。

Gorm 会把用户输入当成参数，而不是 SQL 语句。

还要做：

```text
校验 ID 必须是数字
限制字段长度
使用最小权限数据库账号
避免把用户输入直接拼进 SQL
```

## 4. XSS 是什么

XSS 是跨站脚本攻击。

攻击者把恶意 JavaScript 注入页面，让其他用户打开页面时执行。

例如任务标题：

```html
<script>alert('attack')</script>
```

如果前端直接用 HTML 渲染：

```js
element.innerHTML = todo.title
```

浏览器可能执行脚本。

恶意脚本可能：

```text
偷 Cookie
偷 Token
伪造用户操作
跳转钓鱼网站
修改页面内容
```

## 5. 如何防 XSS

前端不要把用户输入当 HTML 渲染。

危险：

```js
element.innerHTML = todo.title
```

安全：

```js
element.textContent = todo.title
```

后端也可以：

```text
限制标题长度
过滤控制字符
输出时转义特殊字符
不要允许普通字段存 HTML
```

如果使用 Cookie 存登录态，可以设置：

```text
HttpOnly
Secure
SameSite
```

降低被脚本读取的风险。

## 6. CSRF 是什么

CSRF 是跨站请求伪造。

攻击者诱导已登录用户访问恶意页面，借用户身份发请求。

如果用户已经登录 TaskNest，浏览器有 Cookie。

恶意页面可能偷偷发请求：

```html
<img src="https://tasknest.com/api/todos/delete?id=1">
```

如果后端用 GET 删除任务，且只依赖 Cookie，可能误删数据。

## 7. 如何防 CSRF

不要用 GET 做修改操作。

错误：

```text
GET /api/todos/delete?id=1
```

正确：

```text
DELETE /api/todos/1
```

其他防护：

```text
CSRF Token
SameSite Cookie
校验 Origin / Referer
使用 Authorization Header
重要操作二次确认
```

前后端分离系统如果使用：

```text
Authorization: Bearer token
```

CSRF 风险通常比纯 Cookie 低一些。

## 8. HTTPS / TLS 是什么

HTTP 是明文传输。

中间人可能看到：

```text
账号
密码
Token
请求内容
响应内容
```

HTTPS 是 HTTP 的安全版本。

它通过 TLS 加密传输内容。

TLS 主要解决：

```text
加密：防止内容被偷看
完整性：防止内容被篡改
身份认证：确认服务器是真的
```

生产环境基本都应该使用 HTTPS。

常见部署：

```text
用户 -> HTTPS -> Nginx -> HTTP -> TaskNest
```

Nginx 负责 TLS 证书，Go 服务在内网运行。

## 9. 密钥管理是什么

密钥管理是安全保存和使用敏感配置。

常见密钥：

```text
数据库密码
Redis 密码
JWT Secret
第三方 API Key
短信服务密钥
对象存储 AccessKey
SSH 私钥
TLS 证书私钥
```

不要写死在代码里。

错误：

```go
const DBPassword = "123456"
const JWTSecret = "my-secret"
```

更好：

```text
环境变量
配置中心
Kubernetes Secret
云厂商 Secret Manager
```

并且：

```text
不同环境使用不同密钥
密钥不要提交到 Git
密钥泄露后能快速轮换
```

## 10. 哪些文件不要提交

通常不要提交：

```text
.env
.env.local
config.prod.yaml
id_rsa
*.pem
*.key
包含真实密码的配置文件
```

可以提交：

```text
.env.example
config.example.yaml
```

里面只放示例，不放真实密钥。

## 11. 敏感信息脱敏是什么

脱敏是把敏感数据的一部分隐藏。

手机号：

```text
13812345678 -> 138****5678
```

身份证：

```text
110101199001011234 -> 110101********1234
```

邮箱：

```text
test@example.com -> t***@example.com
```

Token：

```text
eyJhbGciOi... -> eyJ***masked***
```

## 12. 哪些信息需要脱敏

常见敏感信息：

```text
密码
Token
手机号
邮箱
身份证号
银行卡号
地址
姓名
API Key
Cookie
Authorization Header
数据库连接串
验证码
```

特别注意：

```text
密码永远不应该明文保存。
密码通常也不应该打印到日志。
Token 不应该打印完整值。
```

## 13. TaskNest 安全检查

SQL 查询：

```go
db.Where("title = ?", title).Find(&todos)
```

不要：

```go
db.Raw("SELECT * FROM todos WHERE title = '" + title + "'")
```

接口设计：

```text
GET 只查询
POST 创建
PUT/PATCH 修改
DELETE 删除
```

不要：

```text
GET /api/todos/delete?id=1
```

日志：

```text
不要打印 password
不要打印完整 token
不要打印 Authorization Header
不要打印数据库密码
```

错误返回：

```text
不要把完整 SQL、堆栈、数据库连接串返回给前端。
```

权限：

```text
后端必须校验权限。
不能只靠前端隐藏按钮。
```

## 14. 常见安全误区

```text
把密码明文存在数据库。
把 JWT Secret 写死在代码里。
把 .env 提交到 Git。
用 GET 做删除操作。
只在前端做权限判断。
日志里打印完整 Token。
错误响应暴露数据库细节。
生产环境不用 HTTPS。
```

这些都要避免。

## 15. 学习路线

建议顺序：

1. 理解 SQL 注入。
2. 理解参数化查询。
3. 理解 XSS 和前端转义。
4. 理解 CSRF 和 Method 设计。
5. 理解 HTTPS/TLS。
6. 理解密码哈希和 Token 安全。
7. 学会用环境变量管理密钥。
8. 学会日志脱敏。
9. 给 TaskNest 做安全检查清单。

