package pb

//go:generate -command buf go run ../gobin-run.go buf

//go:generate_input buf.yaml buf.gen.yaml */*.proto
//go:generate_output buf.lock bufgen/**/*.go
//go:generate buf generate
