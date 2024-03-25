package handlers

import (
	"embed"
	"fmt"
	"net/http"
	"path/filepath"
)

func AssetsHandler(w http.ResponseWriter, r *http.Request, assetsFS embed.FS) {
	fmt.Println("assetsHandler:", r.URL.RequestURI())
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to run ParseForm: assetsHandler", http.StatusInternalServerError)
		return
	}

	asset := r.Form.Get("assets")
	if asset == "" {
		http.Error(w, "Query: 'assets' missing", http.StatusInternalServerError)
		return
	}

	// file := fmt.Sprint("./assets/", asset)
	// http.ServeFile(w, r, file)

	p, err := assetsFS.ReadFile("assets/" + asset)
	if err != nil {
		fmt.Println("index is err")
		return
	}

	atype := filepath.Ext(asset)[1:]
	fmt.Println("atype: ", atype)
	if atype == "js" {
		atype = "javascript"
	}

	fmt.Println("atype2: ", atype)

	w.Header().Set("Content-Type", "text/"+atype)
	w.Write(p)
}
