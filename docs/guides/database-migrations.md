# Database Migrations Guide

本项目使用 [golang-migrate](https://github.com/golang-migrate/migrate) 工具来管理数据库 Schema 变更（DDL）。

## 1. 简介

`golang-migrate` 是一个流行的 Go 语言数据库迁移工具，支持多种数据库后端。它通过版本化的 SQL 文件来管理数据库状态，确保不同环境（开发、测试、生产）中的数据库结构一致性。

## 2. 安装

请根据你的操作系统安装 CLI 工具。

### macOS (Homebrew)

```bash
brew install golang-migrate
```

### Go Install

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

*注意：如果需要支持其他数据库（如 MySQL），请调整 tags 或参考官方文档。*

## 3. 目录结构

所有迁移文件都存放于项目根目录下的 `migrations/` 文件夹中。

文件命名格式为：`{version}_{title}.{up|down}.sql`

- `version`: 序列号（例如：000001），保证执行顺序。
- `title`: 描述性名称（例如：create_users_table）。
- `up.sql`: 执行迁移时运行的 SQL（创建表、添加字段等）。
- `down.sql`: 回滚迁移时运行的 SQL（删除表、删除字段等）。

示例：

```text
migrations/
├── 000001_create_schemas.up.sql
├── 000001_create_schemas.down.sql
├── 000002_create_users_table.up.sql
└── 000002_create_users_table.down.sql
```

## 4. 常用命令

在使用命令前，请确保你有一个可用的 Postgres 数据库连接字符串，例如：
`postgres://user:password@localhost:5432/dbname?sslmode=disable`

### 4.1 创建新迁移

创建一个新的 SQL 迁移文件对。

```bash
# 用法: migrate create -ext sql -dir <migrations_dir> -seq <name>
migrate create -ext sql -dir migrations -seq add_order_table
```

这将生成：
- `migrations/XXXXXX_add_order_table.up.sql`
- `migrations/XXXXXX_add_order_table.down.sql`

### 4.2 执行迁移 (Up)

将数据库升级到最新版本。

```bash
migrate -path migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up
```

或者升级特定的 N 个版本：

```bash
# 向上执行 1 个版本
migrate -path migrations -database "..." up 1
```

### 4.3 回滚迁移 (Down)

回滚所有迁移（慎用！会清空数据）。

```bash
migrate -path migrations -database "..." down
```

通常我们只回滚最近的一个版本：

```bash
# 向下回滚 1 个版本
migrate -path migrations -database "..." down 1
```

### 4.4 查看当前版本

查看数据库当前的迁移版本状态。

```bash
migrate -path migrations -database "..." version
```

### 4.5 修复脏状态 (Force)

如果迁移过程中失败（例如 SQL 语法错误），数据库会被标记为 "dirty" 状态，此时需要手动修复并强制设定版本号。

1. 手动修改数据库或 SQL 文件修复问题。
2. 强制设置版本为上一个成功的版本号（例如失败在 V2，则强制设为 V1）。

```bash
migrate -path migrations -database "..." force 1
```

## 5. 最佳实践

1.  **不可变性**：一旦迁移文件被合并到主分支或已部署，**切勿**修改现有的迁移文件。如果需要更改，请创建一个新的迁移文件。
2.  **原子性**：尽量保持每个迁移文件专注于一个逻辑变更（如"添加用户表"），这有助于回滚和调试。
3.  **本地验证**：在提交代码前，务必在本地执行 `up` 和 `down` 操作，确保 rollback 逻辑也是正确的。
4.  **团队协作**：如果在你开发期间有新的迁移被合并，请重新拉取代码并基于最新的序列号创建你的迁移。
