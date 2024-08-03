package db

//go:generate_input sqlc.yaml schema*.sql ./queries/*.sql
//go:generate_output sqlcgen/*.go
//go:generate sqlc generate
