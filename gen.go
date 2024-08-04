package biblioseeq

// Generate from module definition.
//go:generate_input gen-from-mod.go gen-from-mod/*.go go.mod *.gen-from-mod.go.tmpl */*.gen-from-mod.go.tmpl
//go:generate_output *.gen-from-mod.go */*.gen-from-mod.go
//go:generate go run ./gen-from-mod/run.go
