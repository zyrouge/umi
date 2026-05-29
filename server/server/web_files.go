package server

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"zyrouge.me/umi/config"
)

type WebFilesHandler struct {
	FileSystem http.FileSystem
	RawHandler http.Handler
}

func NewWebFilesHandler(dir string) *WebFilesHandler {
	fileSystem := http.Dir(dir)
	return &WebFilesHandler{
		FileSystem: fileSystem,
		RawHandler: http.FileServer(fileSystem),
	}
}

func (webFilesHandler *WebFilesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Clean(r.URL.Path)
	f, err := webFilesHandler.FileSystem.Open(path)
	if err == nil {
		f.Close()
		webFilesHandler.RawHandler.ServeHTTP(w, r)
		return
	}
	r2 := r.Clone(r.Context())
	r2.URL.Path = "/"
	webFilesHandler.RawHandler.ServeHTTP(w, r2)
}

func AttachWebFilesRoutes(router *mux.Router) {
	config, err := config.GetConfig()
	if err != nil {
		return
	}
	if config.Server.WebFiles == "" {
		return
	}
	router.PathPrefix("/").Handler(NewWebFilesHandler(config.Server.WebFiles))
}
