# 部署、Docker、K8s、Jenkins 和自动化流水线学习说明

本文档用于理解后端项目从本地开发到服务器部署的常见流程，重点解释：

- 部署流程是什么。
- 为什么要用 Docker。
- 为什么要容器化部署。
- 什么是 K8s。
- 什么是 K8s 集群。
- 什么是 Docker 集群。
- Docker 的原理是什么。
- 如何使用 Docker。
- Dockerfile 怎么写。
- Dockerfile 每个命令是什么意思。
- Docker 怎么做内外端口映射。
- Docker 如何挂载磁盘目录。
- K8s 大概怎么用。
- 什么是 Jenkins。
- 为什么需要 Jenkins。
- 如何通过 A 机器调用 B 机器的终端。
- SSH 是怎么调用远程机器的。
- 什么是自动化流水线。

## 1. 先看结论

部署就是把本地写好的程序放到服务器上运行，让别人可以访问。

最原始的方式是：

```text
开发者手动登录服务器
手动上传代码或二进制文件
手动安装依赖
手动启动服务
```

更现代的方式是：

```text
代码提交到 Git
流水线自动构建
Docker 打包成镜像
推送到镜像仓库
服务器或 K8s 拉取镜像
自动启动新版本服务
```

可以先这样理解：

```text
Docker：
把程序和运行环境打包成一个标准镜像。

K8s：
管理很多容器，负责调度、扩容、重启、发布。

Jenkins：
自动执行构建、测试、打包、部署流程。

SSH：
让你从 A 机器远程登录或执行 B 机器上的命令。

自动化流水线：
把人工部署步骤变成自动执行的一条流程。
```

## 2. 什么是部署流程

部署流程是指：

```text
把代码变成可运行服务，并发布到目标环境的过程。
```

一个后端项目常见部署流程是：

```text
1. 开发完成代码
2. 提交代码到 Git
3. 拉取最新代码
4. 安装依赖
5. 编译或打包程序
6. 准备配置文件
7. 准备数据库
8. 启动服务
9. 验证接口是否正常
10. 如果有问题，回滚到旧版本
```

对于 TaskNest 这种 Go 项目，手动部署可以是：

```text
1. 在服务器安装 Go
2. 拉取项目代码
3. 执行 go mod download
4. 执行 go build -o task-nest
5. 执行 ./task-nest
6. 访问 http://服务器IP:8080/health
```

这种方式能跑，但问题也很多：

- 每台服务器都要手动装 Go。
- 依赖和版本容易不一致。
- 手动操作容易漏步骤。
- 回滚麻烦。
- 多台机器部署更麻烦。
- 开发环境和生产环境可能不一致。

Docker 和自动化流水线就是为了解决这些问题。

## 3. 为什么要用 Docker

Docker 的核心价值是：

```text
把应用程序和运行环境一起打包。
```

传统部署时，你可能会遇到：

```text
我本地能跑，服务器跑不了。
```

常见原因：

- Go 版本不一致。
- 系统库不一致。
- 环境变量不一致。
- 目录结构不一致。
- 依赖没有安装。
- 启动命令不一致。

Docker 解决方式是：

```text
把程序、依赖、运行命令、基础系统环境都写进镜像。
```

这样不同机器只要能运行 Docker，就能用同一个镜像启动服务。

简单理解：

```text
没有 Docker：
每台机器都要手动准备运行环境。

有 Docker：
机器只需要能运行容器，应用环境跟着镜像走。
```

## 4. 为什么要容器化部署

容器化部署就是用容器运行服务。

它的好处包括：

- 环境一致。
- 启动快。
- 发布方便。
- 回滚方便。
- 多服务隔离。
- 资源限制清晰。
- 适合自动化部署。
- 适合 K8s 这类平台统一管理。

例如你把 TaskNest 做成 Docker 镜像：

```text
task-nest:v1
task-nest:v2
task-nest:v3
```

发布新版本就是运行新镜像。

如果新版本有问题，可以回滚到旧镜像：

```text
task-nest:v2
```

容器化后，部署从“在服务器上搭环境”变成了：

```text
服务器拉取镜像
启动容器
```

## 5. Docker 的原理是什么

Docker 不是虚拟机。

虚拟机是模拟一整套操作系统。

Docker 容器是共享宿主机内核，在隔离环境里运行进程。

可以这样理解：

```text
虚拟机：
每个应用带一整套操作系统，重但隔离强。

容器：
每个应用是一个隔离进程，共享宿主机内核，轻量、启动快。
```

Docker 主要依赖 Linux 的几类能力：

```text
Namespace：
做隔离，例如进程、网络、文件系统隔离。

Cgroups：
做资源限制，例如限制 CPU、内存。

UnionFS：
做镜像分层，让镜像可以复用和增量构建。
```

Docker 里的几个核心概念：

```text
Image 镜像：
应用和运行环境的打包结果，类似安装包。

Container 容器：
镜像运行起来后的实例，类似正在运行的程序。

Dockerfile：
制作镜像的说明书。

Registry 镜像仓库：
保存镜像的地方，例如 Docker Hub、Harbor。
```

关系是：

```text
Dockerfile -> build -> Image -> run -> Container
```

## 6. 什么是 Docker 集群

Docker 集群通常是指多台机器一起运行和管理 Docker 容器。

常见方式有：

```text
Docker Swarm：
Docker 官方的集群编排方案。

Kubernetes：
更主流的容器编排平台。
```

现在生产环境里提到“容器集群”，大多数时候指 K8s 集群。

Docker 单机适合：

- 本地开发。
- 单机部署。
- 学习容器。
- 小型服务。

K8s 集群适合：

- 多台服务器。
- 多个服务。
- 自动扩缩容。
- 滚动发布。
- 服务发现。
- 高可用部署。

## 7. 如何使用 Docker

常用命令如下。

查看 Docker 版本：

```bash
docker version
```

拉取镜像：

```bash
docker pull nginx
```

查看本地镜像：

```bash
docker images
```

运行容器：

```bash
docker run -d --name my-nginx -p 8080:80 nginx
```

查看运行中的容器：

```bash
docker ps
```

查看所有容器：

```bash
docker ps -a
```

查看容器日志：

```bash
docker logs my-nginx
```

进入容器终端：

```bash
docker exec -it my-nginx sh
```

停止容器：

```bash
docker stop my-nginx
```

删除容器：

```bash
docker rm my-nginx
```

删除镜像：

```bash
docker rmi nginx
```

构建镜像：

```bash
docker build -t task-nest:latest .
```

运行 TaskNest 镜像：

```bash
docker run -d --name task-nest -p 8080:8080 task-nest:latest
```

## 8. Dockerfile 怎么写

Dockerfile 是构建 Docker 镜像的说明书。

对于 TaskNest 这种 Go 项目，可以写一个多阶段构建 Dockerfile。

```dockerfile
# 第一阶段：编译 Go 程序
FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o task-nest .

# 第二阶段：运行 Go 程序
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/task-nest ./task-nest

EXPOSE 8080

CMD ["./task-nest"]
```

这叫多阶段构建：

```text
第一阶段：
用 Go 镜像编译程序。

第二阶段：
用更小的基础镜像运行程序。
```

这样最终镜像不需要包含完整 Go 编译环境，体积会更小。

## 9. Dockerfile 每个命令是什么意思

### 9.1 FROM

指定基础镜像。

```dockerfile
FROM golang:1.26 AS builder
```

意思是：

```text
用 golang:1.26 作为基础环境，并给这个阶段起名 builder。
```

### 9.2 WORKDIR

指定容器里的工作目录。

```dockerfile
WORKDIR /app
```

后续命令都会在 `/app` 目录下执行。

### 9.3 COPY

把宿主机文件复制到镜像里。

```dockerfile
COPY go.mod go.sum ./
```

意思是：

```text
把本地 go.mod 和 go.sum 复制到容器当前目录。
```

```dockerfile
COPY . .
```

意思是：

```text
把当前项目所有文件复制到镜像当前目录。
```

### 9.4 RUN

构建镜像时执行命令。

```dockerfile
RUN go mod download
```

意思是下载 Go 依赖。

```dockerfile
RUN go build -o task-nest .
```

意思是编译 Go 程序，输出可执行文件 `task-nest`。

### 9.5 COPY --from

从前一个构建阶段复制文件。

```dockerfile
COPY --from=builder /app/task-nest ./task-nest
```

意思是：

```text
从 builder 阶段的 /app/task-nest 复制到当前镜像的 ./task-nest。
```

### 9.6 EXPOSE

声明容器内服务监听端口。

```dockerfile
EXPOSE 8080
```

注意：

```text
EXPOSE 只是声明，不等于真正把端口暴露到宿主机。
```

真正让外部访问，要在 `docker run` 时用 `-p` 做端口映射。

### 9.7 CMD

容器启动时默认执行的命令。

```dockerfile
CMD ["./task-nest"]
```

意思是：

```text
容器启动后运行 ./task-nest。
```

### 9.8 ENTRYPOINT

`ENTRYPOINT` 也用于设置容器启动命令。

常见区别：

```text
CMD：
更像默认参数，容易被 docker run 后面的命令覆盖。

ENTRYPOINT：
更像固定入口命令。
```

简单项目先用 `CMD` 就够了。

### 9.9 ENV

设置环境变量。

```dockerfile
ENV GIN_MODE=release
```

容器运行时程序可以读取这个环境变量。

也可以在启动容器时传：

```bash
docker run -e GIN_MODE=release task-nest:latest
```

### 9.10 ARG

构建镜像时使用的变量。

```dockerfile
ARG APP_VERSION=dev
```

构建时传入：

```bash
docker build --build-arg APP_VERSION=v1.0.0 -t task-nest:v1.0.0 .
```

### 9.11 VOLUME

声明容器数据卷。

```dockerfile
VOLUME ["/data"]
```

学习阶段更常用的是在 `docker run` 时用 `-v` 明确挂载目录。

## 10. Docker 怎么做内外端口映射

容器有自己的网络空间。

如果程序在容器内监听 `8080`，外部机器不能天然访问它。

需要用端口映射：

```bash
docker run -d --name task-nest -p 8080:8080 task-nest:latest
```

格式是：

```text
-p 宿主机端口:容器端口
```

例如：

```bash
docker run -p 9000:8080 task-nest:latest
```

意思是：

```text
访问宿主机 9000 端口
转发到容器内 8080 端口
```

访问方式：

```text
http://服务器IP:9000
```

容器内部程序仍然监听：

```text
8080
```

## 11. Docker 如何挂载磁盘目录

容器默认是临时环境。

如果容器删除，容器内产生的数据可能也会丢失。

所以数据库文件、日志文件、上传文件等需要挂载到宿主机目录。

Docker 挂载目录常用 `-v`：

```bash
docker run -d \
  --name task-nest \
  -p 8080:8080 \
  -v /data/task-nest:/app/data \
  task-nest:latest
```

格式是：

```text
-v 宿主机目录:容器目录
```

意思是：

```text
容器访问 /app/data
实际读写的是宿主机 /data/task-nest
```

如果 TaskNest 使用 SQLite，可以把数据库文件放在挂载目录里。

例如：

```text
宿主机：/data/task-nest/tasknest.db
容器内：/app/data/tasknest.db
```

这样容器重建后，数据还在宿主机目录里。

也可以挂载配置文件：

```bash
docker run -d \
  -v /data/task-nest/config.yaml:/app/config.yaml \
  task-nest:latest
```

## 12. 什么是 K8s

K8s 是 Kubernetes 的简称。

Kubernetes 是一个容器编排平台。

它主要用来管理大量容器。

Docker 解决的是：

```text
如何把一个应用打包成容器并运行。
```

K8s 解决的是：

```text
很多容器运行在哪些机器上？
容器挂了怎么办？
服务要扩容怎么办？
新版本怎么发布？
服务之间怎么互相访问？
配置和密钥怎么管理？
```

简单理解：

```text
Docker 管单个容器。
K8s 管一群容器。
```

## 13. 什么是 K8s 集群

K8s 集群是一组运行 Kubernetes 的机器。

通常包括两类节点：

```text
Control Plane：
控制面，负责管理整个集群。

Worker Node：
工作节点，负责真正运行容器。
```

K8s 里常见概念：

```text
Pod：
K8s 中最小部署单元，里面可以有一个或多个容器。

Deployment：
管理 Pod 副本数量和发布策略。

Service：
给 Pod 提供稳定访问入口。

Ingress：
管理外部 HTTP/HTTPS 流量入口。

ConfigMap：
保存普通配置。

Secret：
保存敏感配置，例如密码、Token。

Namespace：
做资源隔离，例如 dev、test、prod。
```

关系可以这样看：

```text
Deployment 创建和管理 Pod
Pod 里面运行容器
Service 负责访问 Pod
Ingress 负责从集群外访问 Service
```

## 14. K8s 大概怎么用

K8s 通常通过 YAML 文件描述你想要的状态。

例如你希望 TaskNest 运行 2 个副本：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: task-nest
spec:
  replicas: 2
  selector:
    matchLabels:
      app: task-nest
  template:
    metadata:
      labels:
        app: task-nest
    spec:
      containers:
        - name: task-nest
          image: task-nest:latest
          ports:
            - containerPort: 8080
```

再创建一个 Service：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: task-nest
spec:
  selector:
    app: task-nest
  ports:
    - port: 80
      targetPort: 8080
```

应用配置：

```bash
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

查看 Pod：

```bash
kubectl get pods
```

查看 Service：

```bash
kubectl get svc
```

查看日志：

```bash
kubectl logs deployment/task-nest
```

进入 Pod：

```bash
kubectl exec -it deployment/task-nest -- sh
```

滚动更新镜像：

```bash
kubectl set image deployment/task-nest task-nest=task-nest:v2
```

回滚：

```bash
kubectl rollout undo deployment/task-nest
```

## 15. 什么是 Jenkins

Jenkins 是一个自动化任务平台。

它最常用于 CI/CD。

```text
CI：
Continuous Integration，持续集成。

CD：
Continuous Delivery / Continuous Deployment，持续交付 / 持续部署。
```

Jenkins 可以帮你自动执行：

- 拉代码。
- 安装依赖。
- 运行测试。
- 编译程序。
- 构建 Docker 镜像。
- 推送镜像到仓库。
- 登录服务器部署。
- 执行 K8s 发布。
- 通知发布结果。

没有 Jenkins 时：

```text
人手动执行每一步。
```

有 Jenkins 后：

```text
代码提交后，流水线自动执行这些步骤。
```

## 16. 为什么要有 Jenkins

Jenkins 的价值是把重复、容易出错的人工操作自动化。

好处包括：

- 减少手动部署错误。
- 每次发布步骤一致。
- 自动运行测试。
- 自动构建镜像。
- 保留构建记录。
- 出问题可以查看日志。
- 方便多人协作。
- 支持定时任务。
- 支持不同环境发布。

例如 TaskNest 发布流程可以变成：

```text
开发者 push 代码
Jenkins 自动拉代码
Jenkins 运行 go test
Jenkins 执行 go build
Jenkins 构建 Docker 镜像
Jenkins 推送镜像到镜像仓库
Jenkins 触发服务器或 K8s 更新服务
```

这样开发者不需要每次手动登录服务器执行一堆命令。

## 17. 什么是自动化流水线

自动化流水线就是把软件交付过程拆成多个阶段，并让系统自动执行。

常见流水线阶段：

```text
Checkout：
拉取代码。

Test：
运行测试。

Build：
编译或打包。

Package：
构建 Docker 镜像。

Push：
推送镜像到镜像仓库。

Deploy：
部署到服务器或 K8s。

Verify：
发布后验证。
```

流水线可以写成 Jenkinsfile：

```groovy
pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Test') {
            steps {
                sh 'go test ./...'
            }
        }

        stage('Build') {
            steps {
                sh 'go build -o task-nest .'
            }
        }

        stage('Docker Build') {
            steps {
                sh 'docker build -t task-nest:latest .'
            }
        }

        stage('Deploy') {
            steps {
                sh 'docker run -d --name task-nest -p 8080:8080 task-nest:latest'
            }
        }
    }
}
```

实际生产里还会加：

- 镜像版本号。
- 镜像仓库登录。
- 多环境选择。
- 人工审批。
- 回滚步骤。
- 通知机器人。

## 18. 怎么通过 A 机器调用 B 机器的终端

最常用方式是 SSH。

例如你在 A 机器上执行：

```bash
ssh user@192.168.1.20
```

意思是：

```text
用 user 这个用户登录 192.168.1.20 这台 B 机器。
```

登录成功后，你就在 B 机器的终端里执行命令。

例如：

```bash
pwd
ls
docker ps
systemctl status nginx
```

也可以不进入交互终端，直接远程执行命令：

```bash
ssh user@192.168.1.20 "docker ps"
```

这就表示：

```text
A 机器通过 SSH 让 B 机器执行 docker ps，并把结果返回给 A。
```

Jenkins 部署到远程服务器时，经常就是这样做的。

## 19. SSH 是怎么调用远程机器的

SSH 的全称是：

```text
Secure Shell
```

它是安全远程登录协议。

它能做两类常见事情：

```text
远程登录：
ssh user@server

远程执行命令：
ssh user@server "command"
```

SSH 登录通常有两种认证方式：

```text
密码登录：
输入远程机器用户密码。

密钥登录：
使用私钥和公钥认证，更适合自动化。
```

密钥登录大概流程：

```text
1. A 机器生成一对密钥：私钥和公钥
2. 私钥保存在 A 机器
3. 公钥放到 B 机器的 ~/.ssh/authorized_keys
4. A 机器 SSH 到 B 机器时，用私钥证明自己身份
5. B 机器验证通过后允许登录
```

生成密钥：

```bash
ssh-keygen -t ed25519 -C "your_email@example.com"
```

复制公钥到服务器：

```bash
ssh-copy-id user@192.168.1.20
```

测试登录：

```bash
ssh user@192.168.1.20
```

使用 SSH 执行部署命令：

```bash
ssh user@192.168.1.20 "cd /data/task-nest && docker compose pull && docker compose up -d"
```

## 20. Jenkins 如何通过 SSH 部署到服务器

一种常见方式：

```text
Jenkins 机器保存 SSH 私钥
B 服务器保存 Jenkins 公钥
Jenkins 流水线通过 ssh 执行远程部署命令
```

示例流程：

```text
1. Jenkins 拉取代码
2. Jenkins 构建 Docker 镜像
3. Jenkins 推送镜像到镜像仓库
4. Jenkins SSH 到服务器
5. 服务器拉取新镜像
6. 服务器停止旧容器
7. 服务器启动新容器
8. Jenkins 调用健康检查接口
```

示例命令：

```bash
ssh deploy@server "docker pull registry.example.com/task-nest:v1.0.0"
ssh deploy@server "docker stop task-nest || true"
ssh deploy@server "docker rm task-nest || true"
ssh deploy@server "docker run -d --name task-nest -p 8080:8080 registry.example.com/task-nest:v1.0.0"
```

注意：

```text
Jenkins 里的 SSH 私钥要用凭据管理，不要写死在代码仓库里。
```

## 21. TaskNest 可以怎么部署

对于当前 TaskNest，可以按三个阶段学习。

### 21.1 第一阶段：手动部署

```text
服务器安装 Go
拉代码
go build
手动启动
curl /health 验证
```

适合学习最基础部署过程。

### 21.2 第二阶段：Docker 部署

```text
编写 Dockerfile
docker build 构建镜像
docker run 启动容器
使用 -p 映射端口
使用 -v 挂载 SQLite 数据目录
```

适合理解容器化。

### 21.3 第三阶段：Jenkins 自动化部署

```text
代码 push 后触发 Jenkins
Jenkins 自动测试
Jenkins 自动构建镜像
Jenkins SSH 到服务器部署
Jenkins 自动健康检查
```

适合理解 CI/CD。

### 21.4 第四阶段：K8s 部署

```text
把镜像推送到镜像仓库
编写 Deployment 和 Service YAML
kubectl apply 发布
用 kubectl logs 查看日志
用 rollout 更新和回滚
```

适合理解容器编排。

## 22. 最后总结

可以这样记：

```text
部署：
把程序发布到服务器并运行。

Docker：
把程序和环境打包成镜像，用容器运行。

容器化：
让部署环境标准化，减少“我本地能跑”的问题。

K8s：
管理很多容器，让服务能扩容、重启、滚动发布。

K8s 集群：
由控制节点和工作节点组成的一组机器。

Docker 集群：
多台机器共同运行容器，常见方案是 Docker Swarm 或 K8s。

Jenkins：
自动执行构建、测试、打包、部署。

SSH：
从一台机器安全登录或控制另一台机器。

自动化流水线：
把软件交付步骤变成可重复、可追踪、可自动执行的流程。
```

