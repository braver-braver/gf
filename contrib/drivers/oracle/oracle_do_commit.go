// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

// isReleaseSavepoint recognizes the unsupported Oracle statement without
// accidentally swallowing a second statement or a malformed command.
func isReleaseSavepoint(statement string) bool {
	statement = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(statement), ";"))
	words := strings.Fields(statement)
	return len(words) == 3 &&
		strings.EqualFold(words[0], "RELEASE") &&
		strings.EqualFold(words[1], "SAVEPOINT") &&
		words[2] != ""
}

// DoCommit commits current sql and arguments to underlying sql driver.
func (d *Driver) DoCommit(ctx context.Context, in gdb.DoCommitInput) (out gdb.DoCommitOutput, err error) {
	// Oracle has no RELEASE SAVEPOINT; the savepoint is discarded by
	// transaction COMMIT or ROLLBACK. Preserve GoFrame's nested-tx behavior.
	if isReleaseSavepoint(in.Sql) {
		result := &Result{}
		out.Result = result
		out.RawResult = result
		return out, nil
	}

	out, err = d.Core.DoCommit(ctx, in)
	if err != nil {
		return
	}
	if len(out.Records) > 0 {
		// remove auto added field.
		for i, record := range out.Records {
			delete(record, rowNumberAliasForSelect)
			out.Records[i] = record
		}
	}
	return
}
