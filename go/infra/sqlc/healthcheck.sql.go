// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: healthcheck.sql

package query

import (
	"context"
)

const healthcheck = `-- name: Healthcheck :one
SELECT 1
`

func (q *Queries) Healthcheck(ctx context.Context) (int32, error) {
	row := q.db.QueryRow(ctx, healthcheck)
	var column_1 int32
	err := row.Scan(&column_1)
	return column_1, err
}
