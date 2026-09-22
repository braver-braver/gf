package oracle

import (
	"context"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
)

func TestStripOracleIdentifierQuotesPreservesSizeColumn(t *testing.T) {
	sql := `INSERT INTO BIO_CUSTOMER_PHOTOS("PHOTO_TYPE","SIZE","CREATED_AT") VALUES(:v1,:v2,:v3) RETURNING "ID" INTO :v4`

	got := stripOracleIdentifierQuotes(sql)

	if !strings.Contains(got, `"SIZE"`) {
		t.Fatalf("expected SIZE column to remain quoted, got %s", got)
	}
	if strings.Contains(got, `"PHOTO_TYPE"`) || strings.Contains(got, `"CREATED_AT"`) || strings.Contains(got, `"ID"`) {
		t.Fatalf("expected non-reserved identifiers to keep legacy quote stripping, got %s", got)
	}
}

func TestReleaseSavepointIsNoOp(t *testing.T) {
	driver := new(Driver)

	out, err := driver.DoCommit(context.Background(), gdb.DoCommitInput{Sql: "RELEASE SAVEPOINT sp1"})
	if err != nil {
		t.Fatalf("DoCommit returned error: %v", err)
	}
	if out.Result == nil || out.RawResult == nil {
		t.Fatal("expected non-nil no-op results")
	}
}
