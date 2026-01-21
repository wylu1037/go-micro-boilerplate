---
trigger: model_decision
---

# Code Style Guide

## Go
You must adhere to the following Go coding standards and best practices in all code generation:

### 1. Concurrency & Safety
- **Zero-value Mutexes:** Use `sync.Mutex` and `sync.RWMutex` directly (no pointers needed unless the struct is large).
- **No Copying Mutexes:** If a struct contains a Mutex, **always** use a pointer receiver (`func (s *S)`) to avoid copying the lock.
- **Atomic Alignment:** For ARM/x86-32 compatibility, keep 64-bit atomic fields at the beginning of the struct.
- **Channel Size:** Default to unbuffered channels or size 1. Avoid large buffers unless strictly necessary for performance metrics.

### 2. Performance & Allocation
- **Pre-allocation:** Always specify capacity for `make([]T, 0, cap)` and `make(map[T]T, cap)` to avoid resizing allocations.
- **Strconv over Fmt:** Use `strconv` instead of `fmt` for converting primitives to strings (it's faster).
- **Avoid string to []byte conversion:** Reuse buffers or use `strings.Builder`.
- **Receivers:**
  - Use **Pointer** receivers for: Mutable objects, structs with Mutex, large structs.
  - Use **Value** receivers for: Small structs (e.g., `time.Time`), basic types, maps/funcs/chans.

### 3. Error Handling
- **No Naked Panics:** Never use `panic` in production code. Return `error` instead.
- **Defer Cleanup:** Handle errors in `defer`. Do not ignore them.
  - *Bad:* `defer file.Close()`
  - *Good:* `defer func() { _ = file.Close() }()` (or log the error).
- **Error Wrapping:** Use `%w` to wrap errors for context. Use `errors.Is` and `errors.As` for checking.

### 4. Code Structure & Syntax
- **JSON Tags:** Always use camelCase for JSON tags.
  - ❌ Bad: `json:"user_id"`
  - ✅ Good: `json:"userId"`
- **Type Aliases:** Use `any` instead of `interface{}` for cleaner code.
- **Variable Shadowing:** **PROHIBITED**. Never redeclare variables in inner scopes that shadow outer scope variables.
  - ❌ Bad: `if err != nil { err := doSomething() }`
  - ✅ Good: `if err != nil { innerErr := doSomething() }`
- **Interface Compliance:** Compile-time check for interface implementation.
  - `var _ MyInterface = (*MyImplementation)(nil)`
- **Enums:** Start `iota` at 1, or ensure the 0-value (default) is handled as "Unknown/Invalid".
- **Functional Options:** Use the Functional Options pattern for complex constructors with optional parameters.
- **Grouping:** Group constants, variables, and imports (`import ( ... )`). Separate standard lib imports from 3rd party.
- **Naked Returns:** Avoid naked returns in non-trivial functions. Explicitly return values for readability.

### 5. Testing
- **Table-Driven Tests:** strictly use table-driven tests (`tests := []struct{...}`) for unit testing logic.
- **Subtests:** Use `t.Run` within loops.

### 6. Logging
- **Field Keys:** Use **snake_case** for all log field keys (compatible with Grafana/Loki/OpenTelemetry conventions).
  - ❌ Bad: `zap.String("traceId", id)`, `zap.String("userId", uid)`
  - ✅ Good: `zap.String("trace_id", id)`, `zap.String("user_id", uid)`
- **Standard Fields:**
  - `trace_id`, `span_id` — OpenTelemetry trace context
  - `request_id` — HTTP request identifier (fallback when OTel unavailable)
  - `user_id`, `service`, `endpoint`, `method`, `duration`
  - `remote_addr`, `user_agent`, `status`, `latency`
- **Trace Context:** Always include `trace_id` and `span_id` from context when available using `tools.ExtractTraceInfo(ctx)`.

## Migration

We use [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations.

### File Naming Convention
```
{version}_{description}.up.sql
{version}_{description}.down.sql
```
- **Version**: Sequential number with leading zeros (e.g., `000001`, `000002`) or Unix timestamp
- **Description**: Snake_case, descriptive (e.g., `create_users_table`, `add_email_index`)

### Examples
```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_add_email_index.up.sql
└── 000002_add_email_index.down.sql
```

### Best Practices
1. **Reversibility**: Every `.up.sql` must have a corresponding `.down.sql`
2. **Idempotency**: Use `IF NOT EXISTS` / `IF EXISTS` to make migrations rerunnable
3. **Data Safety**: Never delete columns/tables directly in production; deprecate first

### Common Commands
```bash
# Create new migration
migrate create -ext sql -dir migrations -seq create_users_table

# Apply all migrations
migrate -path migrations -database "postgres://..." up

# Rollback last migration
migrate -path migrations -database "postgres://..." down 1

# Force version (use with caution)
migrate -path migrations -database "postgres://..." force 5

# Check current version
migrate -path migrations -database "postgres://..." version
```

## Proto (buf)

### Directory Structure
```
proto/<service>/<version>/<file>.proto
proto/common/<version>/          # Shared types
```

### Naming Conventions

| Element | Convention | Example |
|---------|------------|---------|
| Package | `<service>.<version>` | `catalog.v1` |
| Service | PascalCase + `Service` | `CatalogService` |
| Message | PascalCase | `CreateShowRequest` |
| Field | snake_case | `show_id`, `created_at` |
| Enum Value | `TYPE_VALUE` | `SHOW_STATUS_DRAFT` |
| RPC | Verb + Noun | `CreateShow`, `ListShows` |

### Enum Standards
```protobuf
enum ShowStatus {
  SHOW_STATUS_UNSPECIFIED = 0;  // Required
  SHOW_STATUS_DRAFT = 1;
  SHOW_STATUS_ON_SALE = 2;
}
```

### Field Types
- Timestamp: `google.protobuf.Timestamp`
- Money: `int64` (store in cents/smallest unit, avoid floating point)
- ID: `string` (UUID/ULID)
- Optional: `optional` keyword

### Common Commands
```bash
buf lint                              # Check standards
buf breaking --against '.git#branch=main'  # Check compatibility
buf generate                          # Generate code
buf format -w                         # Format files
```