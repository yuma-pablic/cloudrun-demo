-- name: GetUser :many
SELECT
    id,
    account_name,
    authority,
    del_flg,
    created_at
FROM users;
