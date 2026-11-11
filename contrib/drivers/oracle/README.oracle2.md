# Oracle2 compatibility driver

This branch rebases the Oracle compatibility variant onto GoFrame upstream
`master` at `0ced14952e45ded642ec36e09c8200afa5730927` (2026-11-11), retaining the upstream module declaration
`github.com/gogf/gf/contrib/drivers/oracle/v2`.

Compatibility patches retained:
- Register `oracle2` in addition to upstream `oracle`.
- Treat a well-formed `RELEASE SAVEPOINT <name>` as a no-op because Oracle
  does not support releasing individual savepoints.
- When GoFrame's DB debug mode is enabled and `logs/` can be created, place
  go-ora traces at `logs/oracle_trace.log`. If the directory is not
  writable, tracing is skipped; an explicit `TRACE FILE` in `Extra` takes precedence.

Upstream now preserves explicit quoted SQL identifiers in `DoFilter`,
so the old `"SIZE"` placeholder/quote-stripping workaround was removed.
Raw SQL requiring mixed-case identifiers or `"SIZE"` must quote them
explicitly. Framework-generated identifiers follow upstream's unquoted
behavior.

This branch keeps the GoFrame multi-module layout: root and
`contrib/drivers/oracle/` each have their own `go.mod`. The current upstream
Oracle module requires Go 1.25+ and `github.com/sijms/go-ora/v2` v2.9.0.
The module path remains the upstream module path; consumer applications
should use a Go `replace` to this fork and pin a **tested** commit.
Upgrading this branch does not change an already-pinned consumer.

Local checks (from the repository checkout):
```sh
cd contrib/drivers/oracle
go test ./... -run '^TestOracle2' -count=1
```

Database integration tests require an isolated Oracle instance. Do not
interpret Go compilation or unit tests as validation against production 19c.
