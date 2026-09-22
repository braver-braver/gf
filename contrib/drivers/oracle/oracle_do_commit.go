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

// emptyResult implements sql.Result for no-op operations like RELEASE SAVEPOINT.
type emptyResult struct{}

func (emptyResult) LastInsertId() (int64, error) { return 0, nil }
func (emptyResult) RowsAffected() (int64, error) { return 0, nil }

// DoCommit commits current sql and arguments to underlying sql driver.
func (d *Driver) DoCommit(ctx context.Context, in gdb.DoCommitInput) (out gdb.DoCommitOutput, err error) {
	// Oracle releases savepoints when the transaction commits or rolls back and
	// does not support an explicit RELEASE SAVEPOINT statement.
	if strings.HasPrefix(strings.ToUpper(in.Sql), "RELEASE SAVEPOINT") {
		out.Result = emptyResult{}
		out.RawResult = emptyResult{}
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
