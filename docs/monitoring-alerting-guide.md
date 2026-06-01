# 监控和告警学习说明

本文档用于理解：

- 监控是什么。
- Metric 是什么。
- Prometheus 和 Grafana 是什么。
- QPS、RT、错误率、P95、P99 是什么。
- 告警规则怎么设计。
- 如何通过监控排查 TaskNest 问题。

## 1. 先看结论

监控是持续观察系统运行状态。

可以理解成：

```text
监控 = 服务运行状态的仪表盘
```

后端服务上线后，不能只看“能不能启动”。

还要看：

```text
有没有请求进来
接口是否变慢
错误率是否升高
CPU 和内存是否正常
数据库是否变慢
服务是否频繁重启
```

常见组合：

```text
服务暴露 /metrics
Prometheus 采集指标
Grafana 展示图表
Alertmanager 发送告警
```

## 2. Metric 是什么

Metric 是指标。

它是可以被数字化记录的数据。

常见指标：

```text
请求总数
请求耗时
错误请求数
CPU 使用率
内存使用率
磁盘使用率
数据库连接数
接口 QPS
接口 P95 耗时
接口 P99 耗时
```

例如：

```text
GET /api/todos 1 分钟内请求 300 次
POST /api/todos 平均耗时 80ms
5 分钟内出现 5 次 500 错误
服务当前内存占用 200MB
```

这些都是 Metric。

## 3. Prometheus 是什么

Prometheus 是常用开源监控系统。

它负责：

```text
采集指标
存储指标
查询指标
根据规则触发告警
```

Prometheus 通常是拉取模式。

也就是：

```text
TaskNest 暴露 /metrics
Prometheus 每 15 秒请求一次 /metrics
Prometheus 保存指标
Grafana 查询 Prometheus 展示图表
```

示例配置：

```yaml
scrape_configs:
  - job_name: "tasknest"
    scrape_interval: 15s
    static_configs:
      - targets: ["localhost:8080"]
```

## 4. Grafana 是什么

Grafana 是可视化工具。

它通常不负责采集，而是连接数据源展示图表。

常见数据源：

```text
Prometheus
Elasticsearch
Loki
MySQL
PostgreSQL
```

Grafana 可以展示：

```text
QPS 曲线
接口耗时曲线
错误率趋势
CPU / 内存 / 磁盘
P95 / P99
数据库连接数
```

简单理解：

```text
Prometheus 负责采集和存储。
Grafana 负责展示和看图。
```

## 5. QPS 是什么

QPS 是每秒请求数。

```text
QPS = Queries Per Second
```

例如：

```text
QPS = 10    每秒 10 个请求
QPS = 1000  每秒 1000 个请求
```

QPS 用来衡量服务访问压力。

如果 QPS 突然升高，可能是：

```text
用户访问量增加
前端重复请求
爬虫刷接口
接口被攻击
定时任务集中触发
```

## 6. RT 是什么

RT 是响应时间。

```text
Response Time
```

表示：

```text
请求从进入服务到返回响应花了多久。
```

例如：

```text
创建任务 RT = 50ms
查询任务列表 RT = 120ms
```

RT 升高可能是：

```text
数据库慢
第三方服务慢
代码逻辑变慢
CPU 压力高
锁竞争
网络慢
```

## 7. 错误率是什么

错误率是失败请求占比。

例如：

```text
1 分钟总请求 1000
失败 20
错误率 = 2%
```

常见分类：

```text
4xx 错误率：参数错误、未登录、无权限等
5xx 错误率：服务端异常、数据库失败、代码 panic 等
```

最需要关注的是 5xx。

## 8. P95 和 P99 是什么

平均耗时容易掩盖慢请求。

P95 表示：

```text
95% 的请求耗时都小于等于这个值。
```

P99 表示：

```text
99% 的请求耗时都小于等于这个值。
```

例如：

```text
GET /api/todos P95 = 200ms
GET /api/todos P99 = 800ms
```

意思是：

```text
大部分请求比较快，但少数请求可能很慢。
```

P95 / P99 能帮助发现少数用户体验很差的问题。

## 9. 为什么需要告警

没有告警时，服务出问题可能是用户先发现。

有告警后，可以提前知道：

```text
服务不可用
错误率升高
接口变慢
CPU 太高
内存持续上涨
磁盘快满
数据库连接异常
```

告警的目标是：

```text
让团队及时发现真正影响用户的问题。
```

## 10. 告警规则怎么设计

好的告警应该关注用户体验：

```text
能不能打开
打开快不快
操作会不会失败
数据会不会丢
```

常见规则：

```text
服务连续 1 分钟不可用
5 分钟内 5xx 错误率超过 5%
P95 连续 10 分钟超过 1 秒
CPU 连续 5 分钟超过 90%
内存连续 10 分钟超过 85%
磁盘使用率超过 85%
```

不要太敏感：

```text
一次请求失败就告警
CPU 瞬间超过 90% 就告警
```

告警通常要加持续时间。

## 11. 告警分级

P0 严重故障：

```text
服务完全不可用
核心接口大量 5xx
数据库不可用
数据写入失败
```

P1 重要问题：

```text
接口明显变慢
错误率小幅升高
CPU 长时间过高
内存持续上涨
```

P2 普通风险：

```text
磁盘使用率偏高
某个接口 QPS 异常
日志量突然增加
单个实例偶发重启
```

## 12. 好的告警信息应该有什么

好的告警要能指导行动。

应该包含：

```text
哪个服务
哪个接口
当前指标值
阈值是多少
持续多久
Grafana 链接
日志查询链接
最近发布版本
处理建议
```

不好的告警：

```text
服务异常
```

好的告警：

```text
TaskNest POST /api/todos 5xx 错误率 8.2%，超过阈值 5%，已持续 5 分钟。
请检查最近发布、错误日志、数据库连接状态。
```

## 13. TaskNest 应该监控什么

接口指标：

```text
GET /api/todos QPS
GET /api/todos P95 / P99
POST /api/todos QPS
POST /api/todos 5xx 错误率
DELETE /api/todos/:id 4xx / 5xx
```

服务资源：

```text
CPU
内存
goroutine 数量
GC 次数
数据库连接数
磁盘使用率
```

业务指标：

```text
任务创建成功数
任务创建失败数
任务完成数
任务删除数
```

## 14. 如何用监控排查问题

任务列表变慢：

```text
1. 看 GET /api/todos 的 RT、P95、P99
2. 看 QPS 是否暴涨
3. 看数据库耗时
4. 看 CPU、内存
5. 看是否刚发布新版本
6. 查日志里是否有慢 SQL
```

创建任务失败：

```text
1. 看 POST /api/todos 的 5xx 错误率
2. 查 ERROR 日志
3. 看数据库是否可用
4. 看数据库连接数是否耗尽
5. 看最近配置和发布
```

服务不可用：

```text
1. 看服务实例是否存活
2. 看容器是否重启
3. 看 CPU / 内存是否打满
4. 看磁盘是否满
5. 查启动日志和错误日志
```

## 15. 学习路线

建议顺序：

1. 理解 QPS、RT、错误率、P95、P99。
2. 理解 Prometheus 采集指标。
3. 理解 Grafana 展示图表。
4. 为 TaskNest 设计基础指标。
5. 设计服务不可用、5xx、慢接口告警。
6. 把监控、日志、Trace 串起来排查问题。

