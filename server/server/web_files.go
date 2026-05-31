package server

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"zyrouge.me/umi/application"
)

type WebFilesHandler struct {
	FileSystem http.FileSystem
	RawHandler http.Handler
}

func NewWebFilesHandler(dir string) *WebFilesHandler {
	fileSystem := http.Dir(dir)
	handler := WebFilesHandler{
		FileSystem: fileSystem,
		RawHandler: http.FileServer(fileSystem),
	}
	return &handler
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

func AttachWebFilesRoutes(router *mux.Router) error {
	config, err := application.GetConfig()
	if err != nil {
		return err
	}
	if config.Server.WebFiles == "" {
		return nil
	}
	router.PathPrefix("/").Handler(NewWebFilesHandler(config.Server.WebFiles))
	return nil
}
