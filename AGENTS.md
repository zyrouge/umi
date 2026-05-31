# AGENTS.md

## File Structure

### When to Create a New Directory

- New domain area → new top-level package directory
- New cross-cutting concern → new top-level package directory
- New infrastructure layer → new top-level package directory
- New binary → `cmd/<binary-name>/main.go`
- Submodules within a domain get their own subdirectory with their own `router.go` and `AttachRoutes` if possible
- Never nest domain packages inside each other

### File Naming

All filenames use `snake_case`.

Route packages:
- `router.go` — only `AttachRoutes`
- `route_<name>.go` — one file per handler group
- `middleware_<name>.go` — one file per middleware

Repository package:
- `<entity>.go`
- `<entity>_query.go`
- `<entity>_insert.go`
- `<entity>_update.go`
- `<entity>_delete.go`
- `<entity>_sql.go`
- `<entity>_encryption.go`

Even a single-function operation gets its own file.

---

## Code Conventions

### Naming

- All domain types are prefixed with `Umi`: `UmiUser`, `UmiConfig`, `UmiMemberRole`
- Enumerables use typed aliases: `type UmiMemberRole string`
- Constants are named `<Type><Variant>`: `UmiMemberRoleOwner`, `UmiErrorCodeInternal`
- Package-level cache variables use the `cached` prefix: `cachedConnection`, `configCache`
- Package-level singletons use descriptive exported names: `Manager`, `Logger`, `GlobalValidator`
- No abbreviated identifiers: use `config` not `cfg`, `database` not `db`, `connection` not `conn` — standard Go short-lived variables (`err`, `ok`, `i`) are fine
- Always write explicit types on every parameter: `func(username string, userId string, ...)` not `func(username, userId string, ...)`
- Receiver names are one or two lowercase letters matching the type
- Highly refrain from unexported types, structs, methods, and variables — everything stays public

### Pointers

- Single-entity query functions return `(*T, error)` — `nil, nil` means not found, never return `sql.ErrNoRows` to callers
- List query functions return `([]*T, error)`
- Optional or nullable struct fields use `*T`
- Large structs are passed as pointer parameters
- Config and database singletons are stored and returned as `*T`

### Struct Initialization

- Always initialize as a value then return its address:
  ```go
  value := StructType{Field: x}
  return &value
  ```
- Never `value := &StructType{...}`
- Always use named fields, never positional

### Function Bodies

- No blank lines inside function bodies
- If a blank line feels necessary, split into smaller functions

### Return Values

- `(T, error)` or `(*T, error)` for fallible functions
- Mutation-only functions return `error`
- No named return values
- `AttachRoutes(*mux.Router) error` on every router package
- Always create a subrouter for submodules within `AttachRoutes`, passing it to the submodule's own `AttachRoutes`
- HTTP handlers are void — they write errors directly to the response

### Error Handling

- Return immediately on every error
- Always wrap with context: `fmt.Errorf("context: %w", err)`
- `sql.ErrNoRows` absorbed internally: `errors.Is(err, sql.ErrNoRows)` → return `nil, nil`
- Before returning `ErrorCodeInternal`, always log with all relevant context variables:
  ```go
  utils.Logger.Error().Err(err).Str("userId", userId).Msg("failed to query user")
  utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
  return
  ```

### Request / Response Structs

- Always exported
- Named `<Action><Entity>Request` for request structs
- Named `<RouteName>Output` when there is a single return type
- Named `<RouteName>JsonOutput` / `<RouteName>CookieOutput` when there are multiple distinct output types

### Context

- All typed context keys are defined in the `route_data` package
- Always access via `With<Key>(ctx, value)` and `Get<Key>(ctx)` — never use `context.WithValue` directly outside `route_data`

### Middleware

- Signature: `func XxxMiddleware(next http.Handler) http.Handler`
- Returns `http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { ... })`

### Timestamps

- Always stored as `int64` Unix epoch seconds

### Logging

- Use the single global `utils.Logger` (zerolog)
- `.Info()` for lifecycle events, `.Error()` for unexpected failures

### SQL

- Inline SQL as backtick raw string literals
- Sensitive fields tagged `json:"-"`

### Imports

- Grouped: stdlib → third-party → internal (`zyrouge.me/umi/...`)
