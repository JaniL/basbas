package main

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed public
var publicFS embed.FS

var mimeByExt = map[string]string{
	".css": "text/css; charset=utf-8",
	".ico": "image/x-icon",
	".png": "image/png",
}

func serveStatic(w http.ResponseWriter, uri string) {
	name := staticPath(uri)
	if name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := publicFS.ReadFile("public/" + name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", mimeByExt[path.Ext(name)])
	w.Write(data)
}

func staticPath(uri string) string {
	if uri == "/favicon.ico" || uri == "/menu.png" {
		uri = "/public" + uri
	}
	if !strings.HasPrefix(uri, "/public/") {
		return ""
	}
	name := strings.TrimPrefix(uri, "/public/")
	if !fs.ValidPath(name) {
		return ""
	}
	return name
}
