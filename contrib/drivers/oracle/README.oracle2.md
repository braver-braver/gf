# Oracle2 compatibility variant

This branch is based on GoFrame Oracle driver `v2.10.2` and keeps the upstream
module path `github.com/gogf/gf/contrib/drivers/oracle/v2`.

In addition to the upstream `oracle` registration, importing this module also
registers the driver name `oracle2`. The variant preserves these application
compatibility fixes:

- ignore unsupported `RELEASE SAVEPOINT` statements;
- preserve quotes around an Oracle column named `SIZE` while retaining the
  driver's legacy quote stripping for other identifiers;
- write go-ora debug traces below `logs/` when that directory is writable.

Consumers should pin a tested commit from this branch with a Go module
`replace` directive. Updating the branch does not update an already pinned
consumer.
