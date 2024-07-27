package biblioseeq

// Generate from module definition.
//go:generate_input gen_from_mod/*.go go.mod *.gen_from_mod.go.tmpl */*.gen_from_mod.go.tmpl
//go:generate_output *.gen_from_mod.go */*.gen_from_mod.go
//go:generate go run ./gen_from_mod/
