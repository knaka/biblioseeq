-- name: GetFiles :many
SELECT * FROM files;

-- name: AddFtsFile :execlastid
INSERT INTO fts_files (body) VALUES (sqlc.arg(body));

-- name: AddFile :execlastid
INSERT INTO files (
  path,
  fts_file_id,
  modified_at,
  size,
  updated_at
)
VALUES (
  sqlc.arg(path),
  sqlc.arg(fts_file_id),
  sqlc.arg(modified_at),
  sqlc.arg(size),
  current_timestamp
)
;

-- name: GetFile :one
SELECT * FROM files WHERE path = sqlc.arg(path) LIMIT 1;

-- name: GetFileWithBody :one
SELECT sqlc.embed(files), sqlc.embed(fts_files)
FROM
  files INNER JOIN
  fts_files ON files.fts_file_id = fts_files.rowid
WHERE
  files.path = sqlc.arg(path)
LIMIT 1
;

-- name: UpdateFtsFile :exec
UPDATE fts_files
SET body = sqlc.arg(body)
WHERE rowid = (
  SELECT fts_file_id FROM files WHERE path = sqlc.arg(path)
)
;


-- name: DeleteFtsFiles :exec
DELETE FROM fts_files
WHERE rowid IN (
  SELECT fts_file_id
  FROM files
  WHERE
    path = sqlc.arg(path) or
    path LIKE sqlc.arg(path) || '/%' or
    path LIKE sqlc.arg(path) || '\%'
)
;

-- name: DeleteFiles :exec
DELETE FROM files
WHERE
  path = sqlc.arg(path) or
  path LIKE sqlc.arg(path) || '/%' or
  path LIKE sqlc.arg(path) || '\%'
;

-- name: UpdateFile :exec
UPDATE files
SET
  modified_at = sqlc.arg(modified_at),
  size = sqlc.arg(size),
  updated_at = current_timestamp
WHERE
  path = sqlc.arg(path)
;

-- name: Query :many
SELECT sqlc.embed(files), snippet(fts_files, 0, '<b>', '</b>', '...', 30) as snippet
FROM
  files INNER JOIN
  fts_files ON files.fts_file_id = fts_files.rowid
WHERE
  fts_files.body MATCH sqlc.arg(query)
ORDER BY rank
LIMIT CASE WHEN sqlc.arg(limit) > 0 THEN sqlc.arg(limit) ELSE 50 END
;
