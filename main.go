package main

import (
	"embed"
	"test-data-builder/internal/builder"
)

//go:embed templates/test-builder.tmpl
var f embed.FS

func main() {
	gen := builder.NewGenerator(f)
	gen.Generate()
}
