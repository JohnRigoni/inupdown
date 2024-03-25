package main

import (
	"embed"
	"inupdown/internal/api/server"
)

//go:embed assets
var assetsFS embed.FS

func main() {
	server.Start(assetsFS)
}

