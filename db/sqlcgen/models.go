// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sqlcgen

import (
	"time"
)

type Document struct {
	Title string
	Body  string
}

type File struct {
	Path       string
	DocumentID int64
	ModifiedAt time.Time
	Size       int64
	UpdatedAt  time.Time
}

type Log struct {
	ID        int64
	Message   string
	CreatedAt time.Time
}
