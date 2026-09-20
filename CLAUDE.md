# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
make test          # run all tests with race detector and coverage
make bench         # run all benchmarks
make lint          # run golangci-lint (no config file; uses defaults)
```

Single package: `go test -cover -race ./errors/...`
Single test: `go test -cover -race -run TestRetryerDo ./errors/...`

## Architecture

A single Go module (`github.com/go-playground/pkg/v5`) containing independent utility packages that extend the Go standard library. Each package uses an `ext` suffix to avoid naming collisions: `errors/` declares `package errorsext`, `sync/` declares `package syncext`, etc.

Only two external dependencies: `go-playground/assert/v2` (testing) and `go-playground/form/v4` (HTTP form decoding).

### Version-Gated Build Tags

Files targeting different Go versions use build tags. When two implementations exist, the older uses `go1.18 && !go1.21` and the newer uses `go1.21`. Examples: `unsafe/conversions.go` vs `unsafe/conversions_go121.go`, `values/option/option_sql.go` vs `option_sql_go1.22.go`.

### Key Patterns

**Rust-inspired types**: `Result[T, E]` (`values/result/`) and `Option[T]` (`values/option/`) with `Ok`/`Err`/`Some`/`None` constructors and `Unwrap`/`UnwrapOr`/`AndThen` methods. Both support JSON and SQL interfaces.

**Builder pattern with value receivers**: `Retryer` types (in `errors/` and `net/http/`) return copies from configuration methods, so a base retryer can be customized without mutation.

**Guard pattern**: `Mutex2`/`RWMutex2` in `sync/` return guard objects from `Lock()` that bundle the protected value with the unlock mechanism.

**Dot-imports in production code**: Some packages dot-import internal siblings (e.g., `errors/retrier.go` dot-imports `values/result`).

### Inter-Package Dependencies

`net/http/` is the heaviest consumer, depending on `errors/`, `bytes/`, `io/`, `ascii/`, `types/`, and `values/`. The `errors/` package depends on `values/result/`. The `sync/` package depends on both `values/option/` and `values/result/`.

## Testing Conventions

All tests use `github.com/go-playground/assert/v2` via dot-import:

```go
import (
    . "github.com/go-playground/assert/v2"
)

func TestFoo(t *testing.T) {
    Equal(t, actual, expected)
}
```

Tests live in the same package as the code (not `_test` external packages).

## Go Version Policy

Supports the two most recent Go versions (MSGV policy). See CONTRIBUTING.md for details.
