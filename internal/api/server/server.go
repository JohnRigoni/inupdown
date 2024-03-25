package server

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"inupdown/internal/api/middleware"
	"inupdown/internal/api/handlers"

	"github.com/gorilla/mux"
)

func Start(assetsFS embed.FS) error {

	assetHandle := middleware.Assets(handlers.AssetsHandler, assetsFS)
	indexHandle := middleware.Assets(handlers.IndexHandler, assetsFS)
	
	r := mux.NewRouter()
	r.PathPrefix("/").Queries("assets", "{assets}").Handler(assetHandle)

	r.PathPrefix("/").Queries("api", "{api}").Handler(http.HandlerFunc(handlers.ApiHandler))
	r.PathPrefix("/").Handler(http.HandlerFunc(indexHandle))

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8088", r))

	return nil
}
