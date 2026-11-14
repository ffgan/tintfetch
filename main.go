package main

import (
	"embed"

	"github.com/ffgan/tintfetch/cmd"
)

//go:embed assets/*
var assets embed.FS

func main() {
	cmd.Run(assets)
}
