# dcdump: modernization and hardening notes

Audit of the codebase (~580 lines Go across 5 files) with concrete improvement
ideas, grouped by severity.

## Critical: resource leaks and crash paths

### `defer resp.Body.Close()` inside retry loop (`harvest.go:64`)

Each retry iteration defers a new `resp.Body.Close()`, but none execute until
`HarvestBatch` returns. Failed responses leak their bodies for the duration of
the function. Fix: close the body explicitly at the end of each loop iteration
(or at the top of the next).

```go
// before retrying or breaking, close the current body
resp.Body.Close()
```

### `log.Fatal` in goroutine (`cmd/dcdump/main.go:148`)

A single failed time slice calls `log.Fatal`, which calls `os.Exit(1)` and
kills every in-flight worker without cleanup. Replace with error collection
(e.g. `errgroup`) so partial failures don't abort the entire run.

### potential nil dereference (`harvest.go:38`)

If `maxRetries` is reached and all attempts returned transport errors (where
`resp` is never assigned), `resp.StatusCode` panics. Guard with a nil check or
track last status separately.

## High: deprecated APIs

### replace `io/ioutil` (deprecated since Go 1.16)

| call | file | replacement |
|------|------|-------------|
| `ioutil.TempFile` | `harvest.go:28` | `os.CreateTemp` |
| `ioutil.ReadAll` | `atomic/file.go:15` | `io.ReadAll` |
| `ioutil.TempFile` | `atomic/file.go:26` | `os.CreateTemp` |

Drop the `io/ioutil` import entirely after migrating.

### replace `math/rand` for temp filenames (`atomic/file.go:71`)

`rand.Intn(99999999)` without seeding is deterministic across runs (pre-Go
1.20) and collision-prone. Use `crypto/rand` or `os.CreateTemp` to generate the
intermediate filename instead of manual randomness.

## High: error handling

### no backoff on transport errors (`harvest.go:60-62`)

Transport errors are retried in a tight loop with no sleep or backoff. This can
hammer the server. Add exponential backoff with jitter, or at minimum reuse the
existing `sleep` duration.

### errors lack context

Wrap errors with `fmt.Errorf("...: %w", err)` to preserve the chain. Currently
callers see generic messages without the originating URL or retry count.

### `MustParse` panics (`dateutil/flag.go:101`)

Used at package init time for defaults (`main.go:44`). A panic here produces a
stack trace instead of a user-friendly message. Consider returning an error and
handling it in `main`.

## Medium: concurrency

### use `errgroup` instead of manual `sync.WaitGroup` + semaphore

`golang.org/x/sync/errgroup` gives bounded concurrency, error propagation, and
context cancellation in one package. Replaces the hand-rolled channel semaphore
at `main.go:115` and the `wg` at `main.go:116`.

```go
g, ctx := errgroup.WithContext(context.Background())
g.SetLimit(*workers)
for _, iv := range intervals {
    g.Go(func() error {
        return unrollPages(ctx, iv.Start, iv.End, *directory, *prefix)
    })
}
if err := g.Wait(); err != nil {
    log.Fatalf("harvest failed: %v", err)
}
```

### no context propagation

HTTP requests have no timeout or cancellation. Thread a `context.Context`
through `unrollPages` -> `HarvestBatch` -> `http.NewRequestWithContext` so
workers can be cancelled on signal or on first fatal error.

### no graceful shutdown

No signal handling (SIGINT/SIGTERM). Long-running harvests can't be stopped
cleanly. Add `signal.NotifyContext` and pass the context to workers.

### unrecovered panics in goroutines (`main.go:144`)

If `unrollPages` panics, the goroutine crashes without recovery and the
`wg.Done()` defer fires but other workers continue unaware. Add
`defer recover()` or switch to `errgroup` which handles this.

## Medium: type safety (`api.go`)

The `DOIResponse` struct uses `interface{}` for seven fields (acknowledged by
the existing TODO at line 4). This disables compile-time type checking and makes
downstream consumers cast everything.

Options:

1. **`json.RawMessage`** -- defer parsing without losing type info; cheapest
   fix, keeps the raw JSON bytes for re-serialization.
2. **Proper struct types** -- define concrete types for `Attributes`,
   `AlternateName`, etc. based on the DataCite schema.
3. **Generated types** -- if DataCite publishes an OpenAPI spec, generate
   structs from it.

Option 1 is the pragmatic minimum since `dcdump` mostly passes data through to
ndjson output anyway.

## Medium: TeeReader write-then-parse coupling (`harvest.go:71-81`)

`TeeReader` simultaneously writes to the output file and feeds the JSON
decoder. If decoding fails mid-stream the file already contains partial data.
Better: write to the temp file first, then read back for metadata extraction, or
buffer the response body and write only after successful decode.

## Low: logging inconsistency

`main.go` mixes `log.Printf` (stdlib) with `log.Warnf` (logrus, aliased as
`log`). Since logrus is already imported, use it consistently. The stdlib `log`
calls bypass logrus formatting and level filtering.

## Low: hardcoded constants

| value | location | suggestion |
|-------|----------|------------|
| `maxRetries = 10` | `harvest.go:22` | flag or config |
| page threshold `400` | `harvest.go:91` | named constant |
| time layout `"20060102150405"` | `main.go:66` | package-level `const` |
| `"dcdump"` user-agent | `harvest.go:47` | include version |

## Low: `MoveFile` cleanup gap (`atomic/file.go:76-77`)

If `io.Copy` fails, the temp file `dsttmp` is left on disk. Add a deferred
removal that's cleared on success:

```go
defer func() { if err != nil { os.Remove(dsttmp) } }()
```

## Low: Makefile targets

Add `test`, `lint`, and `vet` targets:

```makefile
.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	go vet ./...
	staticcheck ./...
```

## Low: the "throwaway" comment (`main.go:5`)

The comment "THIS IS THROWAWAY CODE" dates to 2019. The project has been
maintained for 5+ years and harvested 78 GB of data. Remove or replace with
an accurate description of the tool's purpose and status.

## Won't-fix / out of scope

- **Full OpenAPI types for DataCite** -- large effort, low payoff since the tool
  mostly passes JSON through.
- **HTTP/2 or gRPC** -- the DataCite API is REST/HTTP1.1.
- **Prometheus metrics** -- overkill for a batch CLI tool.
- **Config file support** -- flags are sufficient for current usage.

## Suggested order of work

1. Fix the `defer resp.Body.Close()` leak and nil-deref crash (critical, small
   diff).
2. Replace `log.Fatal` in goroutine with `errgroup` (critical, moderate diff).
3. Replace deprecated `io/ioutil` calls (mechanical, small diff).
4. Add `context.Context` plumbing and signal handling.
5. Add backoff on transport errors.
6. Switch `interface{}` fields to `json.RawMessage`.
7. Clean up logging, constants, and temp-file randomness.
8. Add test targets and basic tests.
