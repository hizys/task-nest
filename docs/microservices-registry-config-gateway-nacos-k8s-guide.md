# 微服务、注册中心、配置中心、网关、Nacos 和 K8s 学习说明

本文档用于理解微服务体系里常见的基础设施和理论：

- 什么是注册中心。
- 什么是配置中心。
- 为什么要有注册中心和配置中心。
- 什么是微服务。
- 为什么要微服务化，有什么好处。
- 什么是网关，为什么要有网关。
- 什么是 CAP。
- 什么是 BASE 理论。
- 为什么很多人用 Nacos 做配置中心和注册中心。
- Nacos 有什么优势。
- Nacos 做注册中心和配置中心如何使用。
- 为什么 K8s 也能做配置中心和注册中心。
- K8s 和 Nacos 对比有什么优缺点。

## 1. 先看结论

可以先用一句话理解这些概念：

```text
微服务：
把一个大系统拆成多个小服务，每个服务负责一块业务。

注册中心：
记录每个服务在哪里，方便服务之间互相发现和调用。

配置中心：
集中管理配置，让服务不用重新打包就能调整配置。

网关：
统一接收外部请求，再转发到内部服务。

CAP：
分布式系统里一致性、可用性、分区容错性不能三者同时完美满足。

BASE：
在高可用系统里，允许短暂不一致，最终达到一致。

Nacos：
常用的服务注册发现和配置管理平台。

K8s：
容器编排平台，也能通过 Service、DNS、ConfigMap、Secret 提供注册发现和配置管理能力。
```

如果是单体应用：

```text
用户 -> 一个后端服务 -> 一个数据库
```

如果是微服务：

```text
用户 -> 网关 -> 多个后端服务 -> 各自数据库或共享基础设施
```

微服务变多以后，就需要解决两个关键问题：

```text
服务在哪里？       -> 注册中心
配置怎么统一管理？ -> 配置中心
```

## 2. 什么是微服务

微服务是一种系统架构方式。

它把一个大应用拆成多个小服务，每个服务负责一个相对独立的业务能力。

例如一个任务管理系统后续变复杂后，可能拆成：

```text
用户服务 user-service
任务服务 todo-service
通知服务 notification-service
文件服务 file-service
统计服务 analytics-service
```

每个服务可以：

- 单独开发。
- 单独测试。
- 单独部署。
- 单独扩容。
- 单独回滚。

对比单体应用：

```text
单体应用：
所有功能在一个项目里，一个进程启动。

微服务：
不同功能拆成多个服务，多个进程或容器运行。
```

## 3. 为什么要微服务化

微服务化通常是为了解决大系统的问题。

当系统越来越大，单体应用可能会遇到：

- 代码越来越大，理解困难。
- 多人开发容易互相影响。
- 一个功能发布要带上整个系统。
- 某个模块出问题可能拖垮整个应用。
- 某个热点模块无法单独扩容。
- 技术栈难以局部升级。

微服务化后的好处：

- 服务边界更清楚。
- 团队可以按业务域拆分。
- 可以独立发布。
- 可以独立扩容。
- 故障隔离更好。
- 不同服务可以选择不同技术栈。
- 更适合大型团队协作。

例如：

```text
通知服务请求量很大。
```

单体应用里，你只能扩容整个系统。

微服务里，可以只扩容：

```text
notification-service
```

## 4. 微服务不是越早越好

微服务有好处，也有成本。

微服务会带来：

- 服务调用变复杂。
- 部署变复杂。
- 数据一致性变复杂。
- 日志排查变复杂。
- 需要注册中心、配置中心、网关、链路追踪。
- 需要更成熟的 CI/CD 和监控体系。

所以小项目一开始不一定要微服务。

对于 TaskNest 这种学习项目，初期单体结构更合适：

```text
main.go
routes
controllers
services
models
```

等业务复杂后，再考虑拆服务。

## 5. 什么是注册中心

注册中心是服务目录。

它记录：

```text
有哪些服务
每个服务有哪些实例
每个实例的 IP 和端口是什么
实例是否健康
```

例如：

```text
todo-service
  - 10.0.0.11:8080
  - 10.0.0.12:8080

user-service
  - 10.0.0.21:8080
  - 10.0.0.22:8080
```

服务启动时，会把自己注册到注册中心。

服务关闭或异常时，注册中心会把它下线或标记为不健康。

服务调用其他服务时，不需要写死 IP，而是问注册中心：

```text
user-service 在哪里？
```

注册中心返回可用实例，调用方再发起请求。

## 6. 为什么要有注册中心

如果没有注册中心，服务之间调用可能要写死地址：

```text
http://10.0.0.21:8080/users/1
```

这样问题很多：

- 服务 IP 变化后要改配置。
- 服务扩容后调用方不知道新实例。
- 服务下线后调用方还可能继续访问旧实例。
- 多环境地址难维护。
- 服务数量变多后管理困难。

有注册中心后，调用方式变成：

```text
调用 user-service
```

而不是关心具体 IP。

注册中心解决的是：

```text
服务发现问题。
```

常见注册中心：

```text
Nacos
Eureka
Consul
Zookeeper
Etcd
Kubernetes Service + DNS
```

## 7. 什么是配置中心

配置中心是集中管理配置的系统。

配置包括：

- 数据库地址。
- Redis 地址。
- 第三方接口地址。
- 开关配置。
- 限流阈值。
- 超时时间。
- 日志等级。
- 灰度比例。
- 业务规则参数。

没有配置中心时，配置可能写在：

```text
配置文件
环境变量
代码常量
启动参数
```

配置中心的作用是：

```text
集中存储配置
按环境隔离配置
支持动态更新配置
支持配置版本管理
支持配置回滚
```

## 8. 为什么要有配置中心

如果没有配置中心，修改配置通常要：

```text
改配置文件
重新打包
重新部署
重启服务
```

这对生产环境不方便。

例如你想把日志等级从 `INFO` 改成 `DEBUG`，用于临时排查问题。

没有配置中心：

```text
改文件 -> 发版 -> 重启服务
```

有配置中心：

```text
在配置中心修改日志等级 -> 服务监听配置变化 -> 自动生效
```

配置中心解决的是：

```text
配置集中管理和动态变更问题。
```

## 9. 什么是网关

网关是系统的统一入口。

外部请求先进入网关，再由网关转发到内部服务。

例如：

```text
用户请求
  |
  v
API Gateway
  |
  +--> user-service
  +--> todo-service
  +--> file-service
  +--> notification-service
```

网关通常负责：

- 路由转发。
- 鉴权。
- 限流。
- 熔断。
- 跨域。
- 统一日志。
- 统一 Trace。
- 黑白名单。
- 请求和响应改写。
- 聚合多个服务结果。

常见网关：

```text
Nginx
Kong
Apache APISIX
Spring Cloud Gateway
Traefik
Envoy
Kubernetes Ingress
```

## 10. 为什么要有网关

如果没有网关，前端可能要直接调用多个后端服务：

```text
前端 -> user-service
前端 -> todo-service
前端 -> file-service
```

这样会有问题：

- 前端需要知道很多服务地址。
- 每个服务都要重复做鉴权。
- 每个服务都要处理跨域。
- 限流和安全策略分散。
- 内部服务直接暴露给外部，不安全。

有网关后：

```text
前端只调用网关。
网关负责转发到内部服务。
```

网关解决的是：

```text
统一入口、统一治理、保护内部服务。
```

## 11. 什么是 CAP

CAP 是分布式系统里的一个重要理论。

它说的是，在发生网络分区时，分布式系统不能同时完美满足：

```text
C：Consistency，一致性
A：Availability，可用性
P：Partition Tolerance，分区容错性
```

### 11.1 一致性 C

一致性指所有节点看到的数据是一样的。

例如：

```text
用户刚修改配置为 value=2。
任何服务再读取配置，都应该读到 value=2。
```

### 11.2 可用性 A

可用性指每次请求都能得到响应。

例如：

```text
即使部分节点异常，系统仍然尽量返回结果。
```

### 11.3 分区容错性 P

分区容错性指网络出现问题时，系统仍能继续工作。

例如：

```text
机房 A 和机房 B 网络暂时不通。
系统仍然要考虑怎么处理请求。
```

在分布式系统中，网络故障无法完全避免，所以通常必须考虑 P。

因此实际选择往往是在：

```text
CP：优先一致性，牺牲部分可用性。
AP：优先可用性，允许短暂不一致。
```

## 12. 什么是 BASE 理论

BASE 是对强一致性的放松，更适合高可用互联网系统。

BASE 包括：

```text
BA：Basically Available，基本可用
S：Soft State，软状态
E：Eventually Consistent，最终一致性
```

### 12.1 基本可用

系统出现故障时，不一定保证所有功能完美可用，但核心功能尽量可用。

例如：

```text
下单功能可用，但推荐功能暂时不可用。
```

### 12.2 软状态

系统允许中间状态存在。

例如：

```text
订单已创建，但库存扣减消息还在处理中。
```

### 12.3 最终一致性

系统不要求每一刻都强一致，但最终会达到一致。

例如：

```text
订单创建成功后，积分稍后到账。
```

BASE 适合解释很多微服务系统里的异步处理、消息队列、最终一致设计。

## 13. Nacos 是什么

Nacos 是阿里开源的服务发现和配置管理平台。

它常被用作：

```text
注册中心
配置中心
服务管理平台
```

Nacos 支持：

- 服务注册。
- 服务发现。
- 服务健康检查。
- 动态配置。
- 配置分组。
- 配置命名空间。
- 配置版本管理。
- 配置监听和动态刷新。
- Web 控制台管理。

在 Spring Cloud Alibaba 体系里，Nacos 使用非常广。

## 14. 为什么很多人用 Nacos 做配置中心和注册中心

常见原因：

- 同时支持注册中心和配置中心。
- 有 Web 控制台，容易查看和管理。
- 和 Spring Cloud Alibaba 集成好。
- 国内资料多，团队上手快。
- 支持命名空间，方便区分 dev、test、prod。
- 支持配置分组，方便区分不同应用。
- 支持服务健康检查。
- 支持临时实例和持久实例。
- 支持配置动态刷新。
- 部署和使用成本相对低。

相比把注册中心和配置中心拆成两个系统，Nacos 可以一个系统同时解决两类问题。

这对中小团队很有吸引力。

## 15. Nacos 有什么优势

Nacos 的主要优势：

```text
功能整合：
注册中心和配置中心合在一起。

控制台友好：
可以在 Web 页面看服务、改配置。

生态成熟：
和 Spring Cloud Alibaba 配合很好。

动态配置：
配置修改后，客户端可以监听变化。

环境隔离：
Namespace 可以区分不同环境。

配置组织：
Data ID 和 Group 可以管理不同应用和配置。

中文资料多：
学习和排查问题更方便。
```

当然，Nacos 也不是所有场景都必须用。

如果系统已经运行在 K8s 上，并且服务发现和配置需求不复杂，K8s 原生能力也可以满足一部分需求。

## 16. Nacos 做注册中心如何用

Nacos 做注册中心的基本流程：

```text
1. 启动 Nacos Server
2. 服务启动时连接 Nacos
3. 服务把自己的服务名、IP、端口注册到 Nacos
4. 调用方根据服务名从 Nacos 拉取可用实例列表
5. 调用方选择一个实例发起请求
6. Nacos 通过心跳或健康检查判断实例是否可用
```

服务注册后的效果类似：

```text
服务名：todo-service
实例：
  192.168.1.10:8080
  192.168.1.11:8080
```

调用方不写死：

```text
http://192.168.1.10:8080
```

而是调用：

```text
todo-service
```

再由客户端负载均衡选择一个实例。

## 17. Nacos 做配置中心如何用

Nacos 配置中心里常见三个概念：

```text
Namespace：
命名空间，常用于隔离环境，例如 dev、test、prod。

Group：
配置分组，常用于区分项目或业务组。

Data ID：
配置文件 ID，常用于表示某个服务的一份配置。
```

例如：

```text
Namespace: dev
Group: TASK_NEST
Data ID: task-nest.yaml
```

配置内容：

```yaml
server:
  port: 8080

database:
  driver: sqlite
  dsn: ./tasknest.db

log:
  level: info
```

服务启动时：

```text
连接 Nacos
读取 task-nest.yaml
解析配置
启动服务
监听配置变化
```

如果配置中心里的 `log.level` 从 `info` 改成 `debug`，客户端可以收到变更通知，并动态更新日志等级。

## 18. 为什么 K8s 也能做注册中心

K8s 自带服务发现能力。

在 K8s 中，Pod 的 IP 会变化，不适合直接写死。

K8s 用 Service 给一组 Pod 提供稳定访问入口。

例如：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: todo-service
spec:
  selector:
    app: todo-service
  ports:
    - port: 80
      targetPort: 8080
```

其他服务可以通过 DNS 访问：

```text
http://todo-service
```

或者完整域名：

```text
http://todo-service.default.svc.cluster.local
```

K8s 会把请求转发到匹配 `app=todo-service` 的 Pod。

所以在 K8s 内部，Service + DNS 就承担了一部分注册中心能力。

## 19. 为什么 K8s 也能做配置中心

K8s 提供 ConfigMap 和 Secret。

```text
ConfigMap：
保存普通配置。

Secret：
保存敏感配置，例如密码、Token、证书。
```

ConfigMap 示例：

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: task-nest-config
data:
  APP_ENV: dev
  LOG_LEVEL: info
  DB_DSN: /data/tasknest.db
```

Pod 使用 ConfigMap 作为环境变量：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: task-nest
spec:
  template:
    spec:
      containers:
        - name: task-nest
          image: task-nest:latest
          envFrom:
            - configMapRef:
                name: task-nest-config
```

也可以把 ConfigMap 挂载成文件：

```yaml
volumeMounts:
  - name: config
    mountPath: /app/config
volumes:
  - name: config
    configMap:
      name: task-nest-config
```

所以 K8s 也能集中管理服务配置。

## 20. Nacos 和 K8s 做注册中心的区别

| 对比项 | Nacos 注册中心 | K8s Service / DNS |
| --- | --- | --- |
| 主要定位 | 服务注册发现平台 | 容器编排平台内置服务发现 |
| 使用范围 | 可用于 K8s 内外、虚拟机、物理机 | 主要用于 K8s 集群内部 |
| 注册方式 | 应用主动注册到 Nacos | Pod 由 K8s 管理，Service 通过标签选择 |
| 控制台 | Nacos 有服务列表页面 | K8s 通常通过 kubectl 或平台查看 |
| 客户端依赖 | 应用通常需要接入 Nacos SDK | 应用通常直接访问 DNS 名称 |
| 负载均衡 | 客户端负载均衡较常见 | Service 做转发 |
| 多语言 | 支持多语言客户端 | 语言无关，走网络和 DNS |
| 适合场景 | 微服务框架体系、混合部署 | K8s 内部服务互调 |

简单理解：

```text
Nacos：
应用感知注册中心。

K8s：
平台层帮你做服务发现。
```

## 21. Nacos 和 K8s 做配置中心的区别

| 对比项 | Nacos 配置中心 | K8s ConfigMap / Secret |
| --- | --- | --- |
| 动态刷新 | 支持客户端监听配置变化 | 默认更多是注入配置，动态生效要应用或重启配合 |
| 控制台 | 有 Web 控制台编辑配置 | 通常通过 YAML、kubectl 或平台管理 |
| 版本管理 | 支持配置历史和回滚 | 原生版本体验较弱，通常依赖 GitOps |
| 环境隔离 | Namespace、Group、Data ID | Namespace、不同资源文件 |
| 使用范围 | K8s 内外都可用 | 主要在 K8s 内部使用 |
| 敏感配置 | 可管理，但需要注意加密和权限 | Secret 专门用于敏感配置 |
| 配置粒度 | 面向应用配置管理 | 面向容器运行配置注入 |

简单理解：

```text
Nacos 更像应用配置管理平台。
K8s ConfigMap 更像容器运行配置注入机制。
```

## 22. Nacos 的优缺点

### 22.1 优点

- 注册中心和配置中心一体化。
- Web 控制台直观。
- 适合 Spring Cloud Alibaba 体系。
- 配置动态刷新能力强。
- 支持命名空间和分组。
- 适合混合部署场景。
- 国内使用多，资料多。

### 22.2 缺点

- 需要额外部署和维护 Nacos 集群。
- 应用通常需要接入 Nacos 客户端。
- 又增加了一个基础设施依赖。
- 如果使用不当，配置权限和配置变更风险需要额外治理。
- 在纯 K8s 云原生体系里，可能与 K8s 原生能力有重叠。

## 23. K8s 原生能力的优缺点

### 23.1 优点

- 不需要额外部署注册中心。
- Service 和 DNS 是 K8s 原生能力。
- 应用可以少依赖特定注册中心 SDK。
- 和容器调度、扩缩容天然结合。
- ConfigMap 和 Secret 与 Deployment 集成自然。
- 适合云原生部署。

### 23.2 缺点

- 主要适合 K8s 集群内部。
- 配置动态刷新不如 Nacos 直观。
- 配置版本管理和审核通常要配合 GitOps。
- Web 配置管理体验不如 Nacos。
- 对不在 K8s 里的服务支持不如 Nacos 直接。
- 如果团队不熟悉 K8s，学习成本较高。

## 24. 什么时候用 Nacos，什么时候用 K8s 原生能力

适合用 Nacos 的场景：

- 使用 Spring Cloud Alibaba。
- 服务既有 K8s 内，也有 K8s 外。
- 希望注册中心和配置中心统一管理。
- 需要配置动态刷新。
- 希望通过 Web 页面管理配置。
- 团队已有 Nacos 运维经验。

适合用 K8s 原生能力的场景：

- 服务都部署在 K8s 内。
- 服务发现通过 Service 和 DNS 已经足够。
- 配置变更走 GitOps 或发布流程。
- 不希望维护额外中间件。
- 应用希望减少注册中心 SDK 依赖。
- 团队熟悉 K8s。

常见现实选择：

```text
传统 Spring Cloud 微服务：
Nacos 很常见。

纯云原生 K8s 服务：
K8s Service + ConfigMap + Secret 更常见。

混合架构：
可能同时使用 K8s 和 Nacos。
```

## 25. TaskNest 可以怎么理解这些概念

TaskNest 当前是单体后端项目。

现在不需要马上引入微服务、Nacos、K8s 网关等复杂设施。

但可以这样学习：

### 25.1 当前阶段

```text
TaskNest 单体服务
本地 SQLite
直接通过 localhost:8080 访问
```

重点学习：

- Gin 路由。
- Gorm 数据库操作。
- 分层结构。
- 基础部署。

### 25.2 服务变多后

如果以后拆成：

```text
user-service
todo-service
notification-service
```

就会遇到：

```text
todo-service 怎么找到 user-service？
配置怎么按 dev/test/prod 管理？
前端到底调用哪个服务？
```

这时就需要：

```text
注册中心
配置中心
网关
```

### 25.3 部署到 K8s 后

如果 TaskNest 部署到 K8s：

```text
Service 负责服务发现
ConfigMap 负责普通配置
Secret 负责敏感配置
Ingress 负责外部入口
```

这时可以先用 K8s 原生能力，不一定马上引入 Nacos。

## 26. 最后总结

可以这样记：

```text
微服务：
把大系统拆成多个小服务。

注册中心：
解决服务在哪里的问题。

配置中心：
解决配置集中管理和动态变更的问题。

网关：
解决统一入口、安全、路由和流量治理问题。

CAP：
分布式系统里 C、A、P 无法同时完美满足。

BASE：
为了高可用，允许短暂不一致，最终一致。

Nacos：
常用的一体化注册中心和配置中心。

K8s：
通过 Service/DNS 提供服务发现，通过 ConfigMap/Secret 提供配置管理。
```

选择 Nacos 还是 K8s 原生能力，不是绝对的。

关键看：

```text
服务是否都在 K8s 内
是否需要动态配置刷新
是否已有 Spring Cloud Alibaba 体系
团队是否愿意维护 Nacos
配置管理是否走 GitOps
```

