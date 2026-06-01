# 数据库、SQL、事务、索引和慢 SQL 优化学习说明

本文档用于理解后端项目最核心的数据基础：

- 数据库是什么。
- 表、字段、行、主键、外键是什么。
- SQL 是什么。
- Gorm 和 SQL 是什么关系。
- 事务和 ACID 是什么。
- 索引是什么。
- 慢 SQL 是什么，怎么优化。

## 1. 先看结论

后端服务需要保存数据。

TaskNest 至少需要保存：

```text
任务标题
任务内容
任务状态
创建时间
更新时间
```

这些数据不能只放在内存里，因为程序一重启，内存数据就没了。

所以需要数据库。

在 TaskNest 里：

```text
Gin        接收 HTTP 请求
Controller 解析参数
Service    处理业务
Gorm       操作数据库
SQLite     保存任务数据
```

## 2. 数据库是什么

数据库是专门存储和管理数据的软件。

常见数据库：

```text
SQLite
MySQL
PostgreSQL
Redis
MongoDB
```

TaskNest 当前使用 SQLite。

SQLite 的特点是：

```text
轻量
不需要单独启动数据库服务
数据存在一个本地文件里
适合学习、小项目、本地开发
```

MySQL 和 PostgreSQL 更常用于生产环境。

## 3. 表、字段和行

关系型数据库把数据存在表里。

### 表 Table

表是一类数据的集合。

例如：

```text
todos       任务表
users       用户表
categories 任务分类表
```

### 字段 Column

字段是表的一列。

例如 `todos` 表可以有：

```text
id
title
content
status
created_at
updated_at
```

### 行 Row

行是表里的一条具体数据。

例如：

```text
id = 1
title = 学习 Gin
status = pending
```

这就是一条任务记录。

## 4. 主键和外键

### 主键

主键是每条记录的唯一标识。

例如：

```text
todos.id
```

通过主键可以快速找到一条数据：

```sql
SELECT * FROM todos WHERE id = 1;
```

### 外键

外键用来表示两张表之间的关系。

例如以后 TaskNest 增加用户系统：

```text
users 表
todos 表
```

一个用户可以有多个任务。

`todos` 表可以增加：

```text
user_id
```

表示这条任务属于哪个用户。

## 5. SQL 是什么

SQL 全称是：

```text
Structured Query Language
结构化查询语言
```

可以理解成：

```text
SQL 是我们和关系型数据库沟通的语言。
```

常见 SQL：

```text
SELECT       查询数据
INSERT       新增数据
UPDATE       修改数据
DELETE       删除数据
CREATE TABLE 创建表
ALTER TABLE  修改表
DROP TABLE   删除表
```

## 6. 常见 SQL 示例

查询所有任务：

```sql
SELECT * FROM todos;
```

根据 ID 查询：

```sql
SELECT * FROM todos WHERE id = 1;
```

查询未完成任务：

```sql
SELECT * FROM todos WHERE status = 'pending';
```

新增任务：

```sql
INSERT INTO todos (title, content, status)
VALUES ('学习 SQL', '理解增删改查', 'pending');
```

修改任务状态：

```sql
UPDATE todos
SET status = 'done'
WHERE id = 1;
```

删除任务：

```sql
DELETE FROM todos WHERE id = 1;
```

## 7. Gorm 和 SQL 的关系

Gorm 是 Go 语言里的 ORM 框架。

ORM 是：

```text
Object Relational Mapping
对象关系映射
```

简单理解：

```text
Gorm 帮你用 Go 结构体操作数据库，底层仍然会生成 SQL。
```

例如：

```go
db.Find(&todos)
```

大致对应：

```sql
SELECT * FROM todos;
```

```go
db.First(&todo, id)
```

大致对应：

```sql
SELECT * FROM todos WHERE id = ? LIMIT 1;
```

```go
db.Create(&todo)
```

大致对应：

```sql
INSERT INTO todos (...) VALUES (...);
```

所以：

```text
学 Gorm 不能完全不学 SQL。
性能排查、慢查询、索引优化最终都要回到 SQL。
```

## 8. 事务是什么

事务用于保证一组数据库操作：

```text
要么全部成功
要么全部失败
```

经典例子是转账：

```text
1. A 账户扣 100
2. B 账户加 100
```

这两个操作必须一起成功。

如果 A 扣钱成功，B 加钱失败，就会出问题。

所以需要事务：

```text
开始事务
执行操作 1
执行操作 2
提交事务
```

如果中间失败：

```text
回滚事务
```

## 9. ACID 是什么

事务有四个核心特性，叫 ACID。

```text
A Atomicity    原子性
C Consistency  一致性
I Isolation    隔离性
D Durability   持久性
```

### 原子性

事务里的操作要么全成功，要么全失败。

### 一致性

事务执行前后，数据要保持合理状态。

### 隔离性

多个事务同时执行时，不能互相随便干扰。

### 持久性

事务提交后，数据要真正保存下来。

## 10. 什么时候需要事务

简单新增一条任务，通常不需要显式事务。

需要事务的场景通常是：

```text
一次操作要修改多张表
一次操作要写多条关键数据
多个操作必须保持一致
任意一步失败都要全部回滚
涉及金额、库存、订单、账户等核心数据
```

TaskNest 未来可能需要事务的场景：

```text
创建任务，同时写操作日志。
删除用户，同时删除用户的任务。
完成任务，同时更新统计表。
```

Gorm 事务示例：

```go
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&todo).Error; err != nil {
        return err
    }

    if err := tx.Create(&operationLog).Error; err != nil {
        return err
    }

    return nil
})
```

## 11. 索引是什么

索引可以帮助数据库更快查询数据。

可以把索引理解成书的目录。

没有目录时，要一页一页找。

有目录时，可以快速定位。

数据库也是一样。

如果经常按 `status` 查询任务：

```sql
SELECT * FROM todos WHERE status = 'done';
```

可以考虑给 `status` 建索引。

## 12. 常见索引类型

主键索引：

```text
主键默认通常有索引。
按 id 查询通常很快。
```

普通索引：

```sql
CREATE INDEX idx_todos_status ON todos(status);
```

唯一索引：

```sql
CREATE UNIQUE INDEX idx_users_email ON users(email);
```

联合索引：

```sql
CREATE INDEX idx_todos_user_status ON todos(user_id, status);
```

联合索引适合经常一起查询的字段：

```sql
SELECT * FROM todos
WHERE user_id = 100 AND status = 'done';
```

## 13. 索引是不是越多越好

不是。

索引有好处，也有代价。

好处：

```text
加快查询
加快排序
加快关联查询
```

代价：

```text
占用磁盘
新增数据要维护索引
修改数据要更新索引
删除数据要更新索引
索引太多会拖慢写入
```

适合加索引的字段：

```text
经常出现在 WHERE 条件里
经常用于 ORDER BY
经常用于 JOIN
区分度比较高
```

## 14. 慢 SQL 是什么

慢 SQL 指执行时间比较长的 SQL。

常见原因：

```text
没有索引
索引用不上
查询返回数据太多
SELECT * 查询了不需要的字段
分页 offset 太大
ORDER BY 没有合适索引
LIKE 前面带 %
循环里反复查询数据库
```

## 15. 什么是 N+1 查询

N+1 是常见性能问题。

错误做法：

```text
1. 查询 100 条任务
2. 循环 100 次，每次查询任务分类
```

总共查询：

```text
1 + 100 = 101 次
```

更好的做法：

```text
一次查询任务列表
一次批量查询分类
在内存中组装数据
```

## 16. 慢 SQL 怎么优化

常见思路：

1. 看接口日志，确认哪个接口慢。
2. 打开 Gorm SQL 日志。
3. 把 SQL 放到数据库里单独执行。
4. 用 `EXPLAIN` 看执行计划。
5. 判断是否全表扫描。
6. 判断是否缺索引。
7. 减少返回字段。
8. 增加分页。
9. 避免 N+1 查询。
10. 优化索引或 SQL。

MySQL 示例：

```sql
EXPLAIN SELECT * FROM todos WHERE user_id = 100;
```

## 17. TaskNest 查询优化示例

如果以后 TaskNest 支持用户，每个用户只能看自己的任务：

```sql
SELECT id, title, status, created_at
FROM todos
WHERE user_id = 100 AND status = 'pending'
ORDER BY created_at DESC
LIMIT 20 OFFSET 0;
```

可以考虑联合索引：

```sql
CREATE INDEX idx_todos_user_status_created
ON todos(user_id, status, created_at);
```

同时接口要分页：

```text
GET /api/todos?status=pending&page=1&page_size=20
```

不要一次返回所有任务。

## 18. 学习路线

建议顺序：

1. 理解数据库、表、字段、行。
2. 理解主键和外键。
3. 学会 `SELECT`、`INSERT`、`UPDATE`、`DELETE`。
4. 理解 Gorm 和 SQL 的关系。
5. 理解事务和 ACID。
6. 理解索引的作用和代价。
7. 学会识别慢 SQL。
8. 学会用 `EXPLAIN` 分析 SQL。
9. 回到 TaskNest 中观察 Gorm 生成的 SQL。

