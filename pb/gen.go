package pb

//go:generate_input buf.yaml buf.gen.yaml */*.proto
//go:generate_output buf.lock bufgen/**/*.go
//go:generate buf generate
