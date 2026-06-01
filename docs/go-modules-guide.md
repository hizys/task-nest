# Go Modules：go.mod 和 go.sum 学习说明

本文档用于学习 Go 项目中的依赖管理机制，重点解释：

- `go.mod` 是什么。
- `go.sum` 是什么。
- 为什么 Go 需要这两个文件。
- 它们类比 Java/Maven 里的什么东西。
- Go 依赖下载到哪里。
- Go 项目的版本号是什么。
- 如何拉取指定版本依赖。
- 如何下载依赖、整理依赖和维护 `go.sum`。

## 1. 先看结论

在 Go 项目里：

```text
go.mod  管项目是谁、用哪个 Go 版本、依赖哪些包和版本
go.sum  管依赖校验值，确保下载到的依赖没有被篡改
```

类比 Java/Maven：

```text
go.mod  类似 Maven 的 pom.xml
go.sum  类似 Maven 的依赖校验/锁定信息，但不是完全等同于 pom.xml
```

更准确地说：

```text
go.mod  = 项目模块声明 + 依赖清单
go.sum  = 依赖完整性校验清单
```

## 2. 为什么要有 go.mod

早期 Go 项目没有 Go Modules，常用 `GOPATH` 管理依赖。

那时候项目通常要放在固定目录下：

```text
$GOPATH/src/xxx
```

这种方式有几个问题：

- 项目必须放在 GOPATH 下面，不够自由。
- 不同项目依赖同一个包的不同版本时，容易冲突。
- 依赖版本不清晰，别人拉代码后不一定能复现你的环境。
- 项目迁移和团队协作不方便。

Go Modules 出现后，项目可以放在任意目录，并通过 `go.mod` 明确记录依赖。

例如 TaskNest 的 `go.mod`：

```go
module github.com/hizys/task-nest

go 1.26.2

require (
    github.com/gin-gonic/gin v1.12.0
    gorm.io/driver/sqlite v1.6.0
    gorm.io/gorm v1.31.1
)
```

这说明：

```text
module github.com/hizys/task-nest
当前项目的模块名，也可以理解成项目唯一标识。

go 1.26.2
当前项目使用的 Go 语言版本。

require
当前项目依赖哪些第三方包，以及使用哪个版本。
```

## 3. 为什么要有 go.sum

`go.sum` 不是手写的，它由 Go 工具自动生成和维护。

它记录的是依赖包的校验值，类似这样：

```text
github.com/gin-gonic/gin v1.12.0 h1:xxxxx
github.com/gin-gonic/gin v1.12.0/go.mod h1:xxxxx
```

它的作用是：**保证你下载到的依赖和之前确认过的依赖是同一份内容。**

如果有人偷偷改了依赖内容，即使版本号没变，校验值也会变。Go 就能发现异常。

简单理解：

```text
go.mod 说：我要用 gin v1.12.0
go.sum 说：gin v1.12.0 的正确指纹应该是 h1:xxxxx
```

这就像你下载一个安装包，不只看名字，还要看文件指纹是否一致。

## 4. go.mod 和 go.sum 要不要提交到 Git

要提交。

一般 Go 项目都应该提交：

```text
go.mod
go.sum
```

原因：

- 让别人拉代码后能下载相同版本的依赖。
- 让 CI/CD 构建环境能复现依赖。
- 保证依赖完整性和安全性。
- 避免每个人本地依赖版本不一致。

不应该提交：

```text
本地编译产物
本地数据库文件
临时文件
```

## 5. 类比 Java 里的什么

如果你学过 Java，可以这样类比。

### Maven 项目

Java Maven 里常见的是：

```text
pom.xml
```

`pom.xml` 里会写：

```xml
<groupId>com.example</groupId>
<artifactId>task-nest</artifactId>
<version>1.0.0</version>

<dependencies>
    <dependency>
        <groupId>xxx</groupId>
        <artifactId>yyy</artifactId>
        <version>1.2.3</version>
    </dependency>
</dependencies>
```

Go 里的 `go.mod` 类似：

```go
module github.com/hizys/task-nest

require github.com/gin-gonic/gin v1.12.0
```

### Maven 本地仓库

Maven 下载的依赖通常放在：

```text
~/.m2/repository
```

Go 下载的依赖通常放在：

```text
$GOMODCACHE
```

可以用命令查看：

```bash
go env GOMODCACHE
```

常见位置是：

```text
~/go/pkg/mod
```

所以类比关系大概是：

```text
Java Maven pom.xml        ≈ Go go.mod
Java Maven ~/.m2          ≈ Go ~/go/pkg/mod
Java Maven dependency     ≈ Go require
Java Maven artifact       ≈ Go module
```

## 6. Go 依赖下载到哪里

Go 的依赖不会下载到你的项目目录里，而是下载到一个统一的模块缓存目录。

查看依赖缓存位置：

```bash
go env GOMODCACHE
```

通常是：

```text
/Users/你的用户名/go/pkg/mod
```

Go 这样设计的原因：

- 多个项目可以复用同一份依赖，节省磁盘空间。
- 项目目录更干净，不会塞满第三方源码。
- 构建时可以从缓存读取依赖，速度更快。
- 依赖版本通过 `go.mod` 控制，不需要把依赖源码放进项目。

这和 Maven 把依赖放到 `~/.m2/repository` 很像。

## 7. Go 从哪里下载依赖

默认情况下，Go 会通过模块代理下载依赖。

查看代理配置：

```bash
go env GOPROXY
```

常见默认值类似：

```text
https://proxy.golang.org,direct
```

意思是：

```text
优先从 Go 官方代理下载
如果代理没有，再直接从源码仓库下载
```

依赖源码可能来自：

```text
GitHub
GitLab
Gitee
官方 Go 模块代理
公司内部私有仓库
```

例如：

```text
github.com/gin-gonic/gin
gorm.io/gorm
```

这些路径既是导入路径，也是 Go 找依赖的模块路径。

## 8. Go 项目的版本号是什么

Go 依赖版本一般使用语义化版本：

```text
v主版本.次版本.修订版本
```

例如：

```text
v1.12.0
v1.31.1
v2.5.0
```

含义：

```text
v1.12.0
1：主版本，可能包含不兼容变更
12：次版本，通常是新增功能
0：修订版本，通常是 bug 修复
```

Go 的版本号通常来自 Git tag。

也就是说，依赖作者在 Git 仓库打了一个 tag：

```bash
git tag v1.12.0
git push origin v1.12.0
```

别人就可以通过 Go 拉取这个版本。

## 9. Go 的 v2 版本为什么特殊

Go Modules 有一个重要规则：

如果模块发布到 `v2` 或更高主版本，模块路径通常要带 `/v2`。

例如：

```text
go.mongodb.org/mongo-driver/v2 v2.5.0
```

这里的 `/v2` 是模块路径的一部分。

原因是：Go 希望不同大版本可以共存，避免 v1 和 v2 API 不兼容时互相冲突。

例如一个项目可以同时依赖：

```text
example.com/lib v1.x.x
example.com/lib/v2 v2.x.x
```

它们在 Go 看来是两个不同模块。

## 10. 怎么拉取依赖

最常用命令是：

```bash
go get 模块名
```

例如添加 Gin：

```bash
go get github.com/gin-gonic/gin
```

添加 Gorm：

```bash
go get gorm.io/gorm
```

添加 SQLite 驱动：

```bash
go get gorm.io/driver/sqlite
```

执行后，Go 会自动更新：

```text
go.mod
go.sum
```

## 11. 怎么拉取指定版本

使用 `@版本号`：

```bash
go get github.com/gin-gonic/gin@v1.12.0
```

指定 Gorm 版本：

```bash
go get gorm.io/gorm@v1.31.1
```

升级到最新版本：

```bash
go get github.com/gin-gonic/gin@latest
```

降级到指定版本：

```bash
go get github.com/gin-gonic/gin@v1.10.0
```

删除一个不再使用的依赖，通常不是手动删 `go.mod`，而是先删除代码里的 import，然后执行：

```bash
go mod tidy
```

## 12. 怎么查看当前依赖

查看所有依赖：

```bash
go list -m all
```

查看某个依赖有哪些可用版本：

```bash
go list -m -versions github.com/gin-gonic/gin
```

查看模块缓存目录：

```bash
go env GOMODCACHE
```

查看当前项目模块信息：

```bash
go env GOMOD
```

## 13. indirect 是什么

你会在 `go.mod` 里看到：

```go
github.com/go-playground/validator/v10 v10.30.1 // indirect
```

`indirect` 表示：这个依赖不是你的代码直接 import 的，而是你的直接依赖间接需要的。

例如：

```text
TaskNest 直接依赖 gin
gin 又依赖 validator
所以 validator 对 TaskNest 来说就是间接依赖
```

依赖关系类似：

```text
TaskNest
  -> gin
      -> validator
      -> sonic
      -> sse
```

## 14. go mod tidy 是什么

`go mod tidy` 是整理依赖的命令。

它会做两件事：

```text
1. 代码里用到了但 go.mod 没写的依赖，加进去
2. go.mod 里写了但代码已经不用的依赖，删掉
```

同时它也会更新 `go.sum`。

常用时机：

- 新增 import 后。
- 删除某个依赖后。
- 提交代码前。
- 拉取别人代码后依赖有变化时。

命令：

```bash
go mod tidy
```

## 15. go mod download 是什么

`go mod download` 只负责下载依赖，不会像 `tidy` 那样根据代码整理依赖。

命令：

```bash
go mod download
```

常用场景：

- CI/CD 构建前先下载依赖。
- Docker 构建镜像时利用缓存。

例如 Docker 里常见写法：

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app .
```

这样可以利用 Docker 缓存：只要 `go.mod` 和 `go.sum` 没变，就不用重复下载依赖。

## 16. go.sum 怎么整理

通常不手动改 `go.sum`。

正确做法：

```bash
go mod tidy
```

如果你想重新校验依赖：

```bash
go mod verify
```

如果你想清理本地模块缓存：

```bash
go clean -modcache
```

注意：`go clean -modcache` 会删除本地所有 Go 依赖缓存，下次构建会重新下载。

## 17. go.mod 里的 replace 是什么

有时候你想临时使用本地版本的依赖，可以用 `replace`。

例如：

```go
replace example.com/common => ../common
```

意思是：

```text
原本要下载 example.com/common
现在改用本地 ../common 目录
```

常见用途：

- 本地同时开发两个 Go 模块。
- 临时修复第三方库问题。
- 公司内部模块本地调试。

注意：提交代码前要确认 `replace` 是否应该保留，否则别人可能找不到你的本地路径。

## 18. TaskNest 当前依赖解释

TaskNest 当前直接依赖：

```text
github.com/gin-gonic/gin v1.12.0
gorm.io/driver/sqlite v1.6.0
gorm.io/gorm v1.31.1
```

含义：

```text
gin
负责 HTTP 路由、请求参数绑定、JSON 响应。

gorm
负责 ORM 数据库操作，把 Go 结构体映射到数据库表。

gorm sqlite driver
让 Gorm 可以连接 SQLite 数据库。
```

间接依赖是这些库内部需要的依赖，例如 JSON 解析、参数校验、网络处理等。

## 19. 常用命令速查

初始化 Go Module：

```bash
go mod init github.com/hizys/task-nest
```

添加依赖：

```bash
go get github.com/gin-gonic/gin
```

添加指定版本：

```bash
go get github.com/gin-gonic/gin@v1.12.0
```

升级到最新版本：

```bash
go get github.com/gin-gonic/gin@latest
```

整理依赖：

```bash
go mod tidy
```

下载依赖：

```bash
go mod download
```

校验依赖：

```bash
go mod verify
```

查看所有依赖：

```bash
go list -m all
```

查看可用版本：

```bash
go list -m -versions github.com/gin-gonic/gin
```

查看依赖缓存目录：

```bash
go env GOMODCACHE
```

清理本地依赖缓存：

```bash
go clean -modcache
```

## 20. 简单记忆

```text
go.mod：我要什么依赖、什么版本
go.sum：这些依赖的指纹是什么
go get：添加或修改依赖版本
go mod tidy：按代码实际使用情况整理依赖
go mod download：只下载依赖
go mod verify：校验依赖有没有被改过
GOMODCACHE：Go 依赖下载后的本地缓存目录
```

如果用 Java 类比：

```text
go.mod         ≈ pom.xml
go.sum         ≈ 依赖校验记录
GOMODCACHE     ≈ ~/.m2/repository
go get         ≈ 添加 dependency
go mod tidy    ≈ 自动整理 dependency
go mod download ≈ 下载 Maven 依赖
```
