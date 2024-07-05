-- name: GetVersion :one
SELECT sqlite_version();

-- name: GetLogs :many
SELECT *
FROM logs
ORDER BY id DESC
LIMIT 10
;
