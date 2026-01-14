---
trigger: always_on
---

# code-style-guide

## migration

- use [migrate](https://github.com/golang-migrate/migrate)

## proto (buf)

### 文件组织

#### 目录结构
- proto 文件按服务和版本组织：`proto/<service>/<version>/<file>.proto`
- 示例：`proto/catalog/v1/catalog.proto`
- 共享类型放在 `proto/common/<version>/` 目录下

#### 文件命名
- 使用小写字母和下划线：`snake_case.proto`
- 文件名应该描述其主要内容
- 一个文件通常包含一个主要 service 及其相关 message

#### Package 命名
- 格式：`<service>.<version>`
- 示例：`package catalog.v1;`
- 版本号使用 `v1`, `v2` 等，不使用 `v1.0`

### 命名规范

#### Service 命名
- 使用 PascalCase
- 以 `Service` 结尾
- 示例：`CatalogService`, `BookingService`

#### Message 命名
- 使用 PascalCase
- Request 消息以 `Request` 结尾
- Response 消息以 `Response` 结尾
- 示例：`CreateShowRequest`, `ListShowsResponse`

#### Field 命名
- 使用 snake_case
- 使用描述性名称
- 示例：`show_id`, `created_at`, `poster_url`
- ID 字段命名：`<resource>_id`（如 `show_id`，而不是 `id` 或 `showId`）

#### Enum 命名
- Enum 类型使用 PascalCase
- Enum 值使用 UPPER_SNAKE_CASE
- Enum 值必须以类型名作为前缀
- 第一个值必须是 `<TYPE>_UNSPECIFIED = 0`
- 示例：
  ```protobuf
  enum ShowStatus {
    SHOW_STATUS_UNSPECIFIED = 0;
    SHOW_STATUS_DRAFT = 1;
    SHOW_STATUS_ON_SALE = 2;
  }
  ```

#### RPC 方法命名
- 使用 PascalCase
- 使用动词开头，描述操作
- 常见动词：`Create`, `Get`, `Update`, `Delete`, `List`
- 示例：`CreateShow`, `GetShow`, `ListShows`, `UpdateShow`

### 编码规范

#### 基本规则
- 必须使用 `syntax = "proto3";`
- 每个 proto 文件必须声明 `package`
- 必须声明 `go_package` option（由 buf managed mode 自动管理）

#### 字段编号
- 字段编号从 1 开始
- 1-15 使用 1 字节编码，保留给最常用字段
- 16-2047 使用 2 字节编码
- 不要使用 19000-19999（protobuf 保留）
- 已删除的字段编号不要重用，使用 `reserved` 标记

#### 注释规范
- Service、Message、Enum、Field 都应该有注释
- 使用 `//` 单行注释，放在定义之前
- 注释应该描述用途和行为，而不是重复名称
- 示例：
  ```protobuf
  // CatalogService manages shows, sessions, and seat inventory
  service CatalogService {
    // Creates a new show with the provided details
    rpc CreateShow(CreateShowRequest) returns (Show);
  }
  ```

#### Import 顺序
- 标准库 import（如 `google/protobuf/timestamp.proto`）
- 第三方依赖 import（如 `buf.build/bufbuild/protovalidate`）
- 项目内部 import（如 `common/v1/pagination.proto`）
- 各组之间用空行分隔

#### 字段类型选择
- 时间使用 `common.v1.Timestamp`（自定义类型）或 `google.protobuf.Timestamp`
- 金额使用 `int64` 存储分（cents），避免浮点数
- 布尔值使用 `bool`
- 可选字段使用 `optional` 关键字（proto3）
- 列表使用 `repeated`
- ID 字段使用 `string` 类型（UUID、ULID 等）

### Buf 配置

#### Lint 规则
项目使用 buf 的 `STANDARD` lint 规则集，但排除了 `PACKAGE_DIRECTORY_MATCH`：
- 启用所有标准 lint 检查
- 不允许使用注释忽略 lint 错误（`disallow_comment_ignores: true`）
- 运行 `buf lint` 检查代码规范

#### Breaking Change 检测
- 项目启用了 `FILE` 级别的 breaking change 检测
- 在修改 proto 文件前，运行 `buf breaking --against '.git#branch=main'` 检查兼容性
- 避免的 breaking changes：
  - 删除或重命名字段
  - 修改字段类型
  - 修改字段编号
  - 删除或重命名 service/rpc
  - 修改 rpc 的请求/响应类型

#### 代码生成
- 使用 `buf generate` 生成代码
- 项目配置了 managed mode，自动管理 `go_package`
- 生成的 Go 代码位于 `gen/go/` 目录
- 使用 `paths=source_relative` 保持目录结构

### 最佳实践

#### 版本管理
- 使用语义化版本：`v1`, `v2`, `v3`
- 新版本创建新的 package，不修改旧版本
- 保持旧版本的 proto 文件不变，确保向后兼容
- 示例：`proto/catalog/v1/` 和 `proto/catalog/v2/` 并存

#### 向后兼容性
- 只添加新字段，不删除旧字段
- 使用 `reserved` 标记已废弃的字段编号和名称
- 新增可选字段使用 `optional` 关键字
- 示例：
  ```protobuf
  message UpdateShowRequest {
    reserved 9;  // 已删除的字段
    reserved "old_field_name";

    string show_id = 1;
    optional string title = 2;
  }
  ```

#### API 设计原则
- 使用资源导向的设计（RESTful 风格）
- CRUD 操作使用标准动词：`Create`, `Get`, `Update`, `Delete`, `List`
- List 操作应该支持分页（使用 `common.v1.PaginationRequest/Response`）
- Update 操作的字段应该使用 `optional`，支持部分更新
- Delete 操作返回简单的成功响应

#### 错误处理
- 使用标准的 gRPC status codes
- 在 Response 中包含 `error_message` 字段（如果需要）
- 不要在 proto 中定义复杂的错误类型，使用 gRPC 的错误机制

#### 常用命令
```bash
# 检查 lint 规范
buf lint

# 检查 breaking changes（对比 main 分支）
buf breaking --against '.git#branch=main'

# 生成代码
buf generate

# 格式化 proto 文件
buf format -w

# 更新依赖
buf dep update
```
