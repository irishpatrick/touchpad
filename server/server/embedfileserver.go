//go:build release

package server

import (
	"embed"
	"net/http"
)

//go:embed dist/*
var fs embed.FS

func NewFileServer(path string) http.Handler {
	return http.FileServer(http.FS(fs))
}
