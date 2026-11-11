// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.

package oracle

import (
	"context"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
)

func TestOracle2ReleaseSavepointRecognition(t *testing.T) {
	tests := []struct {
		name, sql string
		want      bool
	}{
		{"ordinary", "RELEASE SAVEPOINT sp1", true},
		{"whitespace", "  release  savepoint  sp1;  ", true},
		{"missing_name", "RELEASE SAVEPOINT", false},
		{"extra_statement", "RELEASE SAVEPOINT sp1; DELETE FROM t", false},
		{"ordinary_query", "SELECT 1 FROM DUAL", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isReleaseSavepoint(tt.sql); got != tt.want {
				t.Fatalf("isReleaseSavepoint(%q) = %t, want %t", tt.sql, got, tt.want)
			}
		})
	}
}

func TestOracle2ReleaseSavepointNoOp(t *testing.T) {
	driver := new(Driver)
	out, err := driver.DoCommit(context.Background(), gdb.DoCommitInput{Sql: "RELEASE SAVEPOINT sp1"})
	if err != nil {
		t.Fatalf("DoCommit: %v", err)
	}
	if out.Result == nil || out.RawResult == nil {
		t.Fatal("expected non-nil result and raw result")
	}
	rows, err := out.Result.RowsAffected()
	if err != nil || rows != 0 {
		t.Fatalf("RowsAffected() = %d, %v; want 0, nil", rows, err)
	}
}

func TestOracle2PreservesQuotedIdentifiers(t *testing.T) {
	const query = `INSERT INTO BIO_CUSTOMER_PHOTOS("PHOTO_TYPE","SIZE","CREATED_AT") VALUES(?,?,?) RETURNING "ID" INTO ?`
	got := rewriteQuery(convertPlaceholders(query))
	for _, quoted := range []string{`"PHOTO_TYPE"`, `"SIZE"`, `"CREATED_AT"`, `"ID"`} {
		if !strings.Contains(got, quoted) {
			t.Fatalf("missing %s in rewritten SQL: %s", quoted, got)
		}
	}
	for _, bind := range []string{":v1", ":v2", ":v3", ":v4"} {
		if !strings.Contains(got, bind) {
			t.Fatalf("missing bind %s in rewritten SQL: %s", bind, got)
		}
	}
	if strings.Contains(got, "?") {
		t.Fatalf("unconverted placeholders: %s", got)
	}
}
