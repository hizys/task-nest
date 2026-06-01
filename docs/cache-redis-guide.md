# 缓存和 Redis 学习说明

本文档用于理解：

- 缓存是什么。
- Redis 是什么。
- 为什么后端系统需要缓存。
- 缓存常见使用方式。
- 缓存穿透、击穿、雪崩是什么。
- TaskNest 后续可以怎么使用 Redis。

## 1. 先看结论

缓存的核心思想是：

```text
把经常访问、变化不频繁、计算成本高的数据，临时存到更快的地方。
```

Redis 是后端系统里非常常用的缓存组件。

它通常运行在内存中，读写很快。

常见用途：

```text
缓存查询结果
保存登录 Token
保存验证码
做计数器
做限流
做排行榜
做分布式锁
做简单消息队列
```

## 2. 为什么需要缓存

没有缓存时，请求可能每次都查数据库：

```text
用户请求 -> Go 服务 -> 数据库 -> 返回结果
```

如果请求很多，数据库压力会变大。

有缓存后：

```text
用户请求 -> Go 服务 -> Redis
                |
                +-- 有缓存：直接返回
                |
                +-- 没缓存：查数据库，再写入 Redis
```

缓存可以：

- 提高接口响应速度。
- 减少数据库压力。
- 降低重复计算成本。
- 提高系统抗并发能力。

## 3. Redis 是什么

Redis 是一个高性能的内存数据存储系统。

常见数据结构：

```text
String      字符串，最常用
Hash        哈希表，适合存对象字段
List        列表，适合队列
Set         集合，适合去重
Sorted Set  有序集合，适合排行榜
Bitmap      位图，适合签到、状态标记
```

常用命令：

```bash
SET name zhangsan
GET name
DEL name
EXPIRE name 60
TTL name
HSET user:1 name zhangsan
HGET user:1 name
LPUSH queue task1
RPOP queue
```

## 4. 哪些数据适合缓存

适合缓存的数据：

```text
读多写少
变化不频繁
计算成本高
数据库查询较慢
允许短时间不一致
```

例如：

- 用户基本信息。
- 任务分类列表。
- 首页统计结果。
- 热门数据。
- 系统配置。
- 登录 Token。

不适合随便缓存的数据：

```text
强一致要求特别高
变化非常频繁
数据体积特别大
每个用户只访问一次
```

## 5. 缓存查询结果

常见流程：

```text
1. 生成缓存 key
2. 查询 Redis
3. Redis 有数据，直接返回
4. Redis 没数据，查询数据库
5. 把数据库结果写入 Redis
6. 返回数据
```

伪代码：

```go
cacheKey := "todo:list:user:1"

data := redis.Get(cacheKey)
if data != "" {
    return data
}

todos := db.FindTodosByUserID(1)
redis.Set(cacheKey, todos, 5*time.Minute)

return todos
```

## 6. 缓存过期时间

缓存通常要设置过期时间。

例如：

```text
任务列表缓存 1 到 5 分钟
验证码缓存 5 分钟
登录 Token 缓存 7 天
统计数据缓存 1 分钟
```

如果永不过期，可能长期返回旧数据。

## 7. 修改数据后怎么办

如果用户新增、修改、删除任务，原来的缓存可能变旧。

常见做法：

```text
先更新数据库
再删除相关缓存
下次查询时重新加载缓存
```

例如：

```text
创建任务成功
删除 todo:list:user:1:all
删除 todo:list:user:1:status:pending
```

不要只更新数据库，不处理缓存。

否则用户可能看到旧数据。

## 8. 缓存穿透

缓存穿透指：

```text
请求一个数据库里根本不存在的数据。
Redis 查不到。
数据库也查不到。
每次请求都打到数据库。
```

例如一直请求：

```text
GET /api/todos/999999999
```

解决方式：

```text
缓存空结果
参数校验
布隆过滤器
限制恶意请求
```

简单做法：

```text
数据库查不到，也缓存一个空值 30 秒。
```

## 9. 缓存击穿

缓存击穿指：

```text
某个热点 key 过期的一瞬间，大量请求同时打到数据库。
```

解决方式：

```text
热点 key 设置较长过期时间
互斥锁，只允许一个请求查数据库
后台提前刷新缓存
```

## 10. 缓存雪崩

缓存雪崩指：

```text
大量缓存 key 同一时间过期，导致请求同时打到数据库。
```

解决方式：

```text
过期时间加随机值
热点数据后台刷新
Redis 高可用部署
限流和降级
```

例如：

```text
5 分钟 + 0 到 60 秒随机值
```

不要让大量 key 同一秒过期。

## 11. 大 Key 问题

如果一个 Redis key 存了特别大的数据，比如一次存几万个任务，会导致：

```text
读写慢
网络传输慢
阻塞 Redis
删除成本高
```

解决方式：

```text
分页缓存
拆分 key
只缓存必要字段
控制单个 key 大小
```

## 12. TaskNest 示例

### 缓存任务列表

缓存 key 可以设计成：

```text
todo:list:user:1:all
todo:list:user:1:status:pending
todo:list:user:1:status:done
```

查询流程：

```text
1. 根据用户和筛选条件生成 key
2. 先查 Redis
3. 有数据直接返回
4. 没数据查数据库
5. 写入 Redis
6. 返回结果
```

### 修改任务后删除缓存

这些接口修改数据后要删除相关缓存：

```text
POST   /api/todos
PUT    /api/todos/:id
PATCH  /api/todos/:id/status
DELETE /api/todos/:id
```

### 保存登录态

以后 TaskNest 增加登录后，可以用 Redis 保存 Token：

```text
key: auth:token:xxx
value: user_id
ttl: 7 天
```

请求接口时：

```text
1. 从 Header 读取 Token
2. 去 Redis 查询 Token
3. 存在则已登录
4. 不存在则返回 401
```

## 13. 学习路线

建议顺序：

1. 理解缓存思想。
2. 学习 Redis 基础命令。
3. 学习 Redis 数据结构。
4. 在 Go 中连接 Redis。
5. 给任务列表接口加缓存。
6. 修改任务后删除缓存。
7. 理解穿透、击穿、雪崩。
8. 学习 Redis 高可用和监控。

