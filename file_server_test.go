package goweb

import (
	"embed"
	"io/fs"
	"net/http"
	"testing"
)

// in order to use go:mebed command directive you cannot add space after //, so this won't work:
// go:embed [file/dir name], it has to be like below
//
//go:embed resources
var resources embed.FS

func TestFileServer(t *testing.T) {
	directory, _ := fs.Sub(resources, "resources")
	fileServer := http.FileServerFS(directory)

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
