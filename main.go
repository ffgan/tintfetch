package main

import (
	"context"
	"embed"
	"os"
	"os/signal"

	"github.com/ffgan/tintfetch/cmd"
)

//go:embed assets/*
var assets embed.FS

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	os.Exit(cmd.Execute(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr, assets))
}
