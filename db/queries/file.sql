-- name: GetDocument :one
SELECT files.path, files.modtime, documents.title, documents.body
FROM
  documents INNER JOIN
  files ON documents.id = files.document_id
WHERE files.path = sqlc.arg(path)
LIMIT 1
;

-- name: Query :many
SELECT path, modtime, title
FROM
  documents INNER JOIN
  files ON documents.id = files.document_id
WHERE
  documents.title MATCH sqlc.arg(query) OR
  documents.body MATCH sqlc.arg(query)
LIMIT 100
;
