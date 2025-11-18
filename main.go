package main

import (
	"github.com/oskargbc/dws-event-service.git/cmd"
	"embed"
)

//go:embed cmd/*
//go:embed configs/*.go configs/config.yaml
//go:embed internal/*
//go:embed go.mod go.sum main.go Makefile
var EmbedFs embed.FS

func main() {
	cmd.Execute(EmbedFs)
}
