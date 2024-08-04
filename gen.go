package biblioseeq

// Generate from module definition.
//go:generate_input gen-from-mod.go gen_from_mod/*.go go.mod *.gen_from_mod.go.tmpl */*.gen_from_mod.go.tmpl
//go:generate_output *.gen_from_mod.go */*.gen_from_mod.go
//go:generate go run ./gen-from-mod/run.go
