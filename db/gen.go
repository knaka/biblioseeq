package db

//go:generate -command sqlc go run ../gobin-run.go sqlc

//go:generate_input sqlc.yaml schema*.sql ./queries/*.sql
//go:generate_output sqlcgen/*.go
//go:generate sqlc generate
