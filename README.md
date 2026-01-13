# Go Micro Boilerplate - 演唱会票务系统

基于 go-micro v5 + buf + gRPC 的微服务项目模板。

## 技术栈

| 组件 | 选择 |
|------|------|
| 框架 | go-micro.dev/v5 |
| RPC | gRPC + Protobuf (buf) |
| 数据库 | PostgreSQL |
| 缓存 | Redis |
| 消息队列 | NATS |
| 服务发现 | mDNS (dev) / Kubernetes (prod) |

## 项目结构

```
.
├── proto/                  # Protobuf定义
│   ├── common/v1/          # 公共类型
│   ├── identity/v1/        # 身份服务API
│   ├── catalog/v1/         # 目录服务API
│   ├── booking/v1/         # 订单服务API
│   └── notification/v1/    # 通知服务API
├── gen/go/                 # 生成的Go代码
├── pkg/                    # 共享库
│   ├── config/             # 配置管理
│   ├── db/                 # 数据库连接
│   ├── cache/              # Redis封装
│   ├── auth/               # JWT认证
│   ├── middleware/         # gRPC拦截器
│   ├── errors/             # 错误处理
│   └── logger/             # 日志
├── services/               # 微服务
│   ├── gateway/            # API网关
│   ├── identity/           # 身份服务
│   ├── catalog/            # 目录服务
│   ├── booking/            # 订单服务
│   └── notification/       # 通知服务
├── migrations/             # 数据库迁移
└── deploy/                 # 部署配置
```

## 快速开始

### 环境要求

- Go 1.23+
- Docker & Docker Compose
- [buf](https://buf.build/docs/installation)
- [golang-migrate](https://github.com/golang-migrate/migrate)

### 1. 安装依赖

```bash
# 安装buf
brew install bufbuild/buf/buf

# 安装migrate
brew install golang-migrate
```

### 2. 启动基础设施

```bash
# 启动PostgreSQL, Redis, NATS
docker-compose up -d

# 可选：启动调试工具(pgAdmin, Redis Commander)
docker-compose --profile debug up -d
```

### 3. 生成Proto代码

```bash
make gen
```

### 4. 运行数据库迁移

```bash
make migrate-up
```

### 5. 下载Go依赖

```bash
make deps
```

### 6. 运行服务

```bash
# 在不同终端运行各服务
make run-identity
make run-catalog
make run-booking
make run-notification
make run-gateway
```

## 服务说明

### Identity Service (身份服务)
- 用户注册/登录
- JWT Token管理
- 用户资料管理

### Catalog Service (演出目录服务)
- 演出信息管理
- 场次管理
- 座位区域与票价
- 库存管理

### Booking Service (交易核心服务)
- 订单创建与管理
- 库存预扣 (Redis分布式锁)
- 支付对接
- 订单状态机

### Notification Service (通知服务)
- 事件订阅 (NATS)
- 短信/邮件发送
- 消息模板管理

### Gateway Service (API网关)
- HTTP/REST对外暴露
- JWT认证校验
- 限流

## 环境配置

通过环境变量覆盖配置：

```bash
export TICKETING_DATABASE_HOST=localhost
export TICKETING_DATABASE_PASSWORD=secret
export TICKETING_JWT_SECRET=your-secret-key
export MICRO_REGISTRY=kubernetes  # 生产环境
```

## 常用命令

```bash
make help           # 查看所有命令
make gen            # 生成proto代码
make build          # 构建所有服务
make test           # 运行测试
make lint           # 代码检查
make docker-build   # 构建Docker镜像
```

## 开发指南

### 添加新的API

1. 在 `proto/{service}/v1/` 下定义 `.proto` 文件
2. 运行 `make gen` 生成代码
3. 在 `services/{service}/internal/handler` 实现handler
4. 注册handler到服务

### 添加新服务

1. 在 `proto/` 下创建新的proto目录
2. 在 `services/` 下创建新的服务目录
3. 更新 `go.work` 添加新模块
4. 更新 `Makefile` 添加构建目标

## License

MIT
