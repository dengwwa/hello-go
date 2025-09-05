package main

import (
	"embed"
	"io/fs"
	"testing"
)

//go:embed test.yaml
var file embed.FS

func TestEmbed(t *testing.T) {
	readFile, err := fs.ReadFile(file, "test.yaml")
	if err != nil {
		panic(err)
	}
	println(string(readFile))
}
