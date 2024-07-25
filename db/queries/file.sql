-- name: GetDocument :one
SELECT files.path, files.modified_at, documents.title, documents.body
FROM
  documents INNER JOIN
  files ON documents.id = files.document_id
WHERE files.path = sqlc.arg(path)
LIMIT 1
;

-- name: Query :many
SELECT path, modified_at, title
FROM
  documents INNER JOIN
  files ON documents.id = files.document_id
WHERE
  documents.body MATCH sqlc.arg(query)
LIMIT 100
;

-- name: GetDocumentCount :one
-- Works only when fts5 extension is available.
SELECT count(*)
FROM documents
LIMIT 1
;

