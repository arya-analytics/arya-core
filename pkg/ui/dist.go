package ui

import (
	"embed"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"io/fs"
	"net/http"
)

//go:embed dist/*
var _dist embed.FS

const distRoot = "dist"

func distFS() fs.FS {
	assets, err := fs.Sub(_dist, distRoot)
	if err != nil {
		panic(err)
	}
	return assets
}

func distHandler() fiber.Handler {
	return filesystem.New(filesystem.Config{
		Root: http.FS(distFS()),
	})
}
