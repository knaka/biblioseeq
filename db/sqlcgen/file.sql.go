// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: file.sql

package sqlcgen

import (
	"context"
	"time"
)

const getDocument = `-- name: GetDocument :one
SELECT files.path, files.modified_at, documents.title, documents.body
FROM
  documents INNER JOIN
  files ON documents.id = files.document_id
WHERE files.path = ?1
LIMIT 1
`

type GetDocumentRow struct {
	Path       string
	ModifiedAt time.Time
	Title      string
	Body       string
}

func (q *Queries) GetDocument(ctx context.Context, path string) (GetDocumentRow, error) {
	row := q.db.QueryRowContext(ctx, getDocument, path)
	var i GetDocumentRow
	err := row.Scan(
		&i.Path,
		&i.ModifiedAt,
		&i.Title,
		&i.Body,
	)
	return i, err
}

const getDocumentCount = `-- name: GetDocumentCount :one
;

SELECT count(*)
FROM documents
LIMIT 1
`

// Works only when fts5 extension is available.
func (q *Queries) GetDocumentCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getDocumentCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const query = `-- name: Query :many
;

SELECT path, modified_at, title
FROM
  documents INNER JOIN
  files ON documents.id = files.document_id
WHERE
  documents.body MATCH ?1
LIMIT 100
`

type QueryRow struct {
	Path       string
	ModifiedAt time.Time
	Title      string
}

func (q *Queries) Query(ctx context.Context, query string) ([]QueryRow, error) {
	rows, err := q.db.QueryContext(ctx, query, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []QueryRow
	for rows.Next() {
		var i QueryRow
		if err := rows.Scan(&i.Path, &i.ModifiedAt, &i.Title); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
