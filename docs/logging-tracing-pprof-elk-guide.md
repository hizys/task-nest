# 日志、Trace、PProf 和 ELK 学习说明

本文档用于理解后端服务排查问题时常用的可观测性能力：

- 如何打印日志。
- 为什么要打印日志。
- 日志服务应该接什么服务。
- 什么是 Trace。
- 为什么要有 Trace。
- 如何导出 Go 服务的 PProf 文档。
- 如何通过 PProf 和日志排查服务问题。
- 日志等级有哪些。
- 什么场景打印什么等级的日志比较合理。
- 什么是 ELK。

## 1. 先看结论

后端服务上线后，不能只靠“接口能不能访问”判断系统是否正常。

真正排查问题时，通常需要三类信息：

```text
日志 Log：
记录服务运行时发生了什么。

链路 Trace：
记录一次请求经过了哪些服务、每一步耗时多少。

性能剖析 PProf：
记录程序 CPU、内存、协程、阻塞等运行状态。
```

它们解决的问题不同：

```text
日志：发生了什么？
Trace：一次请求经过哪里，慢在哪里？
PProf：程序本身 CPU、内存、协程哪里异常？
```

ELK 则是常见的日志收集和查询系统：

```text
Elasticsearch  存储和搜索日志
Logstash       收集、过滤、处理日志
Kibana         查询和展示日志
```

## 2. 为什么要打印日志

日志的核心作用是：让服务运行过程可追踪。

本地开发时，你可以通过断点、打印、手动请求排查问题。

但是服务上线后，通常不能随便打断点，也不能直接看到用户请求时内部发生了什么。

这时日志就非常重要。

日志可以帮助你回答：

- 请求有没有进来。
- 参数是什么。
- 执行到了哪一步。
- 依赖的数据库或第三方服务是否报错。
- 业务判断为什么走到了某个分支。
- 接口为什么返回失败。
- 某个问题是什么时间开始出现的。
- 同类错误出现了多少次。

简单理解：

```text
没有日志：
线上服务像黑盒，出问题只能猜。

有日志：
可以根据时间、请求 ID、错误信息一步步还原现场。
```

## 3. 如何打印日志

### 3.1 最简单的打印方式

Go 标准库自带 `log` 包。

示例：

```go
import "log"

func main() {
    log.Println("server started")
}
```

这种方式适合学习和简单项目，但生产项目通常会使用结构化日志库。

### 3.2 推荐使用结构化日志

结构化日志不是只打印一段字符串，而是把关键信息按字段打印出来。

普通日志：

```text
create todo failed, user id is 1001, error is database timeout
```

结构化日志：

```json
{
  "level": "error",
  "message": "create todo failed",
  "user_id": 1001,
  "error": "database timeout"
}
```

结构化日志更适合被日志平台搜索和分析。

Go 里常用日志库包括：

```text
log/slog    Go 标准库结构化日志
zap         Uber 开源，高性能
zerolog     高性能 JSON 日志
logrus      老牌日志库
```

对于新项目，可以优先学习：

```text
log/slog
```

因为它是 Go 标准库的一部分。

### 3.3 slog 示例

```go
import (
    "log/slog"
    "os"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    logger.Info("server started",
        "addr", ":8080",
        "service", "task-nest",
    )

    logger.Error("create todo failed",
        "todo_title", "学习日志",
        "error", "database timeout",
    )
}
```

输出类似：

```json
{"time":"2026-06-02T10:00:00Z","level":"INFO","msg":"server started","addr":":8080","service":"task-nest"}
{"time":"2026-06-02T10:00:01Z","level":"ERROR","msg":"create todo failed","todo_title":"学习日志","error":"database timeout"}
```

## 4. 日志应该打印什么

一条有价值的日志通常包含：

```text
时间
日志等级
服务名
环境
请求 ID
Trace ID
用户 ID 或业务 ID
接口路径
关键参数
执行结果
错误信息
耗时
```

例如：

```json
{
  "level": "error",
  "service": "task-nest",
  "env": "prod",
  "trace_id": "abc123",
  "request_id": "req001",
  "method": "POST",
  "path": "/api/todos",
  "todo_title": "学习日志",
  "duration_ms": 120,
  "error": "database timeout",
  "message": "create todo failed"
}
```

注意：不要把敏感信息直接打到日志里。

不应该打印：

- 密码。
- Token。
- 身份证号。
- 银行卡号。
- 完整手机号。
- 用户隐私内容。
- 大段请求体。

## 5. 日志等级分为什么

常见日志等级从低到高是：

```text
TRACE
DEBUG
INFO
WARN
ERROR
FATAL
```

不同库支持的等级可能略有差异。

例如 Go 标准库 `slog` 默认常见等级是：

```text
DEBUG
INFO
WARN
ERROR
```

## 6. 什么场景打印什么等级日志

### 6.1 DEBUG

DEBUG 用于开发和排查问题时查看细节。

适合打印：

- 中间变量。
- 分支判断结果。
- SQL 参数。
- 调用第三方服务的入参摘要。

示例：

```go
logger.Debug("query todos",
    "status", status,
    "page", page,
)
```

生产环境一般不会默认打开大量 DEBUG 日志，否则日志量会太大。

### 6.2 INFO

INFO 用于记录正常关键事件。

适合打印：

- 服务启动。
- 服务关闭。
- 关键业务操作成功。
- 定时任务开始和结束。
- 请求摘要。

示例：

```go
logger.Info("todo created",
    "todo_id", todo.ID,
    "status", todo.Status,
)
```

INFO 是生产环境最常见的日志等级。

### 6.3 WARN

WARN 表示出现异常或不符合预期，但系统还能继续运行。

适合打印：

- 参数不推荐但仍兼容。
- 调用外部服务失败后走了降级。
- 查询结果为空但不影响主流程。
- 重试后成功。
- 用户输入非法导致请求被拒绝。

示例：

```go
logger.Warn("invalid todo status",
    "status", status,
    "path", "/api/todos",
)
```

WARN 表示需要关注，但不一定是系统故障。

### 6.4 ERROR

ERROR 表示发生了真正的错误，需要排查。

适合打印：

- 数据库操作失败。
- 调用外部服务失败且无法恢复。
- 核心业务执行失败。
- 程序出现非预期异常。
- 请求返回 500。

示例：

```go
logger.Error("database create todo failed",
    "title", req.Title,
    "error", err,
)
```

ERROR 日志应该包含足够的上下文，否则只看到 `failed` 没法排查。

### 6.5 FATAL

FATAL 表示严重错误，通常打印后程序会退出。

适合打印：

- 配置文件缺失，服务无法启动。
- 数据库连接初始化失败，服务无法工作。
- 端口监听失败。

示例：

```go
logger.Error("database init failed", "error", err)
os.Exit(1)
```

很多项目不会直接使用 `Fatal` 方法，而是明确打印 ERROR 后退出。

## 7. 日志服务应该接什么服务

本地开发阶段：

```text
直接打印到控制台 stdout
```

单机部署阶段：

```text
打印到文件
配合 logrotate 做日志切割
```

容器或 Kubernetes 部署阶段：

```text
应用打印到 stdout
由容器平台收集日志
```

公司或生产环境常见方案：

```text
ELK / EFK
Loki + Grafana
Datadog
Splunk
阿里云 SLS
腾讯云 CLS
华为云 LTS
AWS CloudWatch
Google Cloud Logging
Azure Monitor
```

如果你是学习项目 TaskNest，可以按阶段理解：

```text
第一阶段：
控制台日志。

第二阶段：
结构化 JSON 日志。

第三阶段：
接入 ELK 或 Loki 这类日志平台。

第四阶段：
日志里加入 trace_id，和 Trace 系统联动。
```

## 8. 什么是 Trace

Trace 是链路追踪。

它用来记录一次请求从入口到结束经过了哪些步骤、每一步耗时多久。

一个 Trace 通常包含多个 Span。

```text
Trace：
一次完整请求链路。

Span：
链路中的一个具体步骤。
```

例如一次创建 Todo 请求：

```text
Trace: POST /api/todos

Span 1: Gin 接收请求，耗时 120ms
Span 2: Controller 解析参数，耗时 2ms
Span 3: Service 校验业务，耗时 5ms
Span 4: Gorm 写入数据库，耗时 90ms
Span 5: 返回 JSON，耗时 3ms
```

如果是微服务系统，一次请求可能经过多个服务：

```text
用户请求
  |
  v
API Gateway
  |
  v
Task Service
  |
  v
User Service
  |
  v
Database
```

Trace 可以把这些调用串起来。

## 9. 为什么要有 Trace

日志能告诉你某个服务里发生了什么。

但如果一次请求经过多个服务，只看单个服务日志会很难排查。

Trace 可以帮你回答：

- 一次请求经过了哪些服务。
- 每个服务耗时多少。
- 慢在哪个环节。
- 哪个下游服务报错。
- 一个错误是否影响了整条链路。
- 同一个用户请求在多个服务中的日志如何关联。

例如用户反馈：

```text
创建任务很慢。
```

没有 Trace 时，你可能要分别查：

```text
网关日志
后端日志
数据库日志
第三方服务日志
```

有 Trace 时，可以看到：

```text
POST /api/todos 总耗时 1300ms
  Controller: 5ms
  Service: 10ms
  Database Insert: 1250ms
```

这样就能快速判断问题主要在数据库写入。

常见 Trace 系统：

```text
Jaeger
Zipkin
Grafana Tempo
SkyWalking
Datadog APM
OpenTelemetry
```

其中 OpenTelemetry 是现在非常重要的标准，它可以统一采集 Trace、Metric、Log。

## 10. 日志和 Trace 的关系

日志和 Trace 应该配合使用。

常见做法是在每条日志里加入：

```text
trace_id
span_id
request_id
```

这样你可以：

```text
先从 Trace 里找到慢请求
再复制 trace_id 去日志系统查询这次请求的详细日志
```

也可以：

```text
先从 ERROR 日志里发现错误
再通过 trace_id 打开完整请求链路
```

简单理解：

```text
Trace 看整体链路。
日志看具体细节。
```

## 11. 什么是 PProf

PProf 是 Go 官方提供的性能分析工具。

它可以帮助你查看 Go 程序运行时的状态。

常见分析内容：

```text
CPU 使用情况
内存分配情况
Goroutine 数量和堆栈
阻塞情况
锁竞争情况
```

PProf 适合排查：

- CPU 飙高。
- 内存持续上涨。
- Goroutine 泄漏。
- 接口变慢。
- 程序阻塞。
- 锁竞争严重。

## 12. 如何在 Go 服务中开启 PProf

Go 标准库提供了 `net/http/pprof`。

只要导入它，并启动一个 HTTP 服务，就可以访问 PProf 接口。

示例：

```go
import (
    "net/http"
    _ "net/http/pprof"
)

func startPProf() {
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
}
```

然后在 `main` 里调用：

```go
func main() {
    startPProf()

    // 启动业务服务
}
```

启动后可以访问：

```text
http://localhost:6060/debug/pprof/
```

常见地址：

```text
/debug/pprof/
/debug/pprof/profile
/debug/pprof/heap
/debug/pprof/goroutine
/debug/pprof/block
/debug/pprof/mutex
```

重要提醒：

```text
不要把 PProf 直接暴露到公网。
```

生产环境建议：

- 只监听 `localhost`。
- 通过内网访问。
- 加鉴权。
- 通过 VPN 或堡垒机访问。
- Kubernetes 中通过 port-forward 临时访问。

## 13. 如何导出服务的 PProf 文档

### 13.1 导出 CPU profile

CPU profile 通常采样一段时间，例如 30 秒。

```bash
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

也可以保存成文件：

```bash
curl -o cpu.pprof "http://localhost:6060/debug/pprof/profile?seconds=30"
```

然后分析：

```bash
go tool pprof cpu.pprof
```

### 13.2 导出内存 heap profile

```bash
curl -o heap.pprof "http://localhost:6060/debug/pprof/heap"
```

分析：

```bash
go tool pprof heap.pprof
```

### 13.3 导出 Goroutine profile

```bash
curl -o goroutine.pprof "http://localhost:6060/debug/pprof/goroutine"
```

也可以直接看文本：

```bash
curl "http://localhost:6060/debug/pprof/goroutine?debug=2"
```

### 13.4 导出阻塞和锁竞争 profile

阻塞：

```bash
curl -o block.pprof "http://localhost:6060/debug/pprof/block"
```

锁竞争：

```bash
curl -o mutex.pprof "http://localhost:6060/debug/pprof/mutex"
```

注意：block 和 mutex 分析通常需要在代码里开启采样。

示例：

```go
runtime.SetBlockProfileRate(1)
runtime.SetMutexProfileFraction(1)
```

## 14. 如何查看 PProf

进入交互模式：

```bash
go tool pprof cpu.pprof
```

常用命令：

```text
top       查看最耗 CPU 或内存的函数
list      查看某个函数的源码级消耗
web       生成调用图
png       导出图片
svg       导出 SVG
pdf       导出 PDF
```

也可以启动网页查看：

```bash
go tool pprof -http=:8081 cpu.pprof
```

然后浏览器打开：

```text
http://localhost:8081
```

如果 `web`、`pdf`、`svg` 不能用，通常是本机没有安装 Graphviz。

macOS 可以安装：

```bash
brew install graphviz
```

## 15. 如何通过日志和 PProf 排查服务问题

### 15.1 接口返回 500

排查顺序：

```text
1. 先查 ERROR 日志
2. 根据 trace_id 或 request_id 找到同一次请求的完整日志
3. 看错误发生在 Controller、Service 还是数据库操作
4. 如果是数据库错误，继续看 SQL、参数、连接状态
5. 修复代码或配置
```

日志重点看：

- 接口路径。
- 请求参数摘要。
- 错误信息。
- 错误堆栈。
- trace_id。

### 15.2 接口变慢

排查顺序：

```text
1. 看接口访问日志里的 duration_ms
2. 找出慢请求的 trace_id
3. 用 Trace 看慢在哪个 Span
4. 如果慢在数据库，检查 SQL 和索引
5. 如果慢在代码计算，导出 CPU profile
6. 如果慢在等待锁或阻塞，查看 block/mutex profile
```

PProf 重点看：

```text
CPU profile
block profile
mutex profile
```

### 15.3 CPU 飙高

排查顺序：

```text
1. 确认 CPU 飙高的时间段
2. 查同一时间段请求量是否异常
3. 查 ERROR/WARN 日志是否异常增多
4. 导出 CPU profile
5. 用 go tool pprof top 查看最耗 CPU 的函数
6. 判断是业务循环、序列化、正则、加密、压缩还是其他计算导致
```

常见原因：

- 死循环。
- 大量 JSON 序列化。
- 正则写得低效。
- 大对象反复处理。
- 请求量突然升高。
- 热点接口没有缓存。

### 15.4 内存持续上涨

排查顺序：

```text
1. 查看内存监控是否持续上涨
2. 查看日志是否出现大量请求或异常
3. 导出 heap profile
4. 用 pprof top 查看哪些对象占用内存最多
5. 多导出几次 heap，对比增长对象
6. 检查是否存在缓存无限增长、大对象未释放、goroutine 泄漏
```

PProf 重点看：

```text
heap profile
goroutine profile
```

### 15.5 Goroutine 泄漏

现象：

```text
服务运行越久 goroutine 越多
内存也可能逐渐上涨
服务最终变慢或崩溃
```

排查：

```bash
curl "http://localhost:6060/debug/pprof/goroutine?debug=2"
```

重点看：

- 很多 goroutine 卡在同一个函数。
- channel 读写没有退出。
- HTTP 请求没有超时。
- 定时任务重复创建。
- 数据库连接等待。

### 15.6 数据库慢

排查顺序：

```text
1. 看接口日志 duration_ms
2. 看 Trace 中 DB Span 耗时
3. 看 Gorm 慢 SQL 日志
4. 检查 where 条件是否走索引
5. 检查是否全表扫描
6. 检查连接池是否耗尽
```

建议打印：

- SQL 摘要。
- 查询条件。
- 耗时。
- 返回条数。
- 错误信息。

## 16. 什么是 ELK

ELK 是一套常见日志平台。

它由三个组件组成：

```text
E：Elasticsearch
L：Logstash
K：Kibana
```

### 16.1 Elasticsearch

Elasticsearch 负责存储和搜索日志。

它适合做：

- 全文搜索。
- 按字段查询。
- 时间范围查询。
- 聚合统计。

例如：

```text
查询过去 1 小时 level=ERROR 的日志
查询 trace_id=abc123 的所有日志
统计某个接口 5xx 错误数量
```

### 16.2 Logstash

Logstash 负责收集、解析、过滤和转发日志。

例如：

```text
读取应用日志文件
解析 JSON 字段
过滤无用字段
把日志写入 Elasticsearch
```

实际生产里也常用 Filebeat 替代一部分 Logstash 的采集工作。

所以你也会看到 EFK：

```text
Elasticsearch
Fluentd / Fluent Bit
Kibana
```

或者：

```text
Filebeat + Elasticsearch + Kibana
```

### 16.3 Kibana

Kibana 负责查询和展示日志。

你可以在 Kibana 里：

- 搜索日志。
- 按 trace_id 查询。
- 看错误趋势。
- 做仪表盘。
- 看某个接口的错误数量。

## 17. 日志进入 ELK 的流程

常见流程是：

```text
应用服务
  |
  v
标准输出或日志文件
  |
  v
Filebeat / Logstash / Fluent Bit
  |
  v
Elasticsearch
  |
  v
Kibana
```

如果 TaskNest 后续接 ELK，建议日志使用 JSON 格式。

例如：

```json
{
  "time": "2026-06-02T10:00:00Z",
  "level": "ERROR",
  "service": "task-nest",
  "env": "dev",
  "trace_id": "abc123",
  "path": "/api/todos",
  "message": "create todo failed",
  "error": "database timeout"
}
```

这样 ELK 能按字段查询：

```text
level = ERROR
service = task-nest
trace_id = abc123
path = /api/todos
```

## 18. TaskNest 后续可以怎么落地

对于当前 TaskNest 项目，可以分阶段接入。

### 18.1 第一阶段：加基础日志

目标：

```text
服务启动时打印日志
每个接口打印请求摘要
错误时打印 ERROR 日志
```

建议内容：

```text
method
path
status_code
duration_ms
error
```

### 18.2 第二阶段：结构化日志

目标：

```text
使用 slog 输出 JSON 日志
统一日志字段
```

建议字段：

```text
time
level
service
env
request_id
method
path
status_code
duration_ms
message
error
```

### 18.3 第三阶段：加入 request_id 和 trace_id

目标：

```text
每个请求生成 request_id
日志中带上 request_id
以后接入 Trace 时再加入 trace_id
```

这样排查单次请求会容易很多。

### 18.4 第四阶段：开启 PProf

目标：

```text
本地或内网开启 PProf
支持导出 CPU、heap、goroutine profile
```

注意：

```text
不要直接暴露到公网。
```

### 18.5 第五阶段：接入日志平台

目标：

```text
把 JSON 日志收集到 ELK、Loki 或云日志服务。
```

学习阶段可以先理解 ELK，不一定马上搭建。

## 19. 最后总结

可以这样记：

```text
日志：
记录发生了什么，适合查错误和业务过程。

Trace：
记录一次请求经过哪里，适合查跨服务调用和慢请求。

PProf：
记录 Go 程序运行状态，适合查 CPU、内存、goroutine、阻塞问题。

ELK：
收集、存储、查询和展示日志的平台。
```

合理的排查方式通常是：

```text
先看监控发现异常
再用 Trace 定位慢在哪段链路
再用日志查看具体错误和业务上下文
如果怀疑程序性能问题，再用 PProf 深入分析
```

