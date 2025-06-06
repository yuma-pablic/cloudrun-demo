// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: select_user.sql

package query

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getUser = `-- name: GetUser :many
SELECT
    id,
    account_name,
    authority,
    del_flg,
    created_at
FROM users
`

type GetUserRow struct {
	ID          int32            `json:"id"`
	AccountName string           `json:"account_name"`
	Authority   bool             `json:"authority"`
	DelFlg      bool             `json:"del_flg"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
}

func (q *Queries) GetUser(ctx context.Context) ([]GetUserRow, error) {
	rows, err := q.db.Query(ctx, getUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUserRow{}
	for rows.Next() {
		var i GetUserRow
		if err := rows.Scan(
			&i.ID,
			&i.AccountName,
			&i.Authority,
			&i.DelFlg,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
