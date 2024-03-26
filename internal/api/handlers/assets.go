package handlers

import (
	"embed"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func AssetsHandler(w http.ResponseWriter, r *http.Request, assetsFS embed.FS) {
	fmt.Println("assetsHandler:", r.URL.RequestURI())
	vars := mux.Vars(r)

	asset := vars["route"]
	if asset == "" {
		http.Error(w, "Query: 'assets' missing", http.StatusInternalServerError)
		return
	}
	fmt.Println("asset: ", asset)

	p, err := assetsFS.ReadFile("assets/" + asset)
	if err != nil {
		fmt.Println("index is err")
		return
	}

	atype := filepath.Ext(asset)[1:]
	if atype == "js" {
		atype = "javascript"
	}

	w.Header().Set("Content-Type", "text/"+atype)
	w.Write(p)
}
