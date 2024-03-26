package server

import (
	"embed"
	"inupdown/internal/api/handlers"
	"inupdown/internal/api/middleware"

	"github.com/gorilla/mux"
)

func Router(assetsFS embed.FS) *mux.Router {
	assetHandle := middleware.Assets(handlers.AssetsHandler, assetsFS)
	indexHandle := middleware.Assets(handlers.IndexHandler, assetsFS)

	r := mux.NewRouter()

	s := r.PathPrefix("/api").Queries("internal", "true").Subrouter()
	r.PathPrefix("/assets/{route}").Queries("internal", "true").Handler(assetHandle)
	r.PathPrefix("/").Handler(indexHandle)

	s.HandleFunc("/delete", handlers.DelHandler)
	s.HandleFunc("/flist", handlers.FlistHandler)
	s.HandleFunc("/raw", handlers.RawHandler)
	s.HandleFunc("/upload", handlers.FileUploadHandler)
	s.HandleFunc("/write", handlers.WriteFileHandler)

	return r
}
