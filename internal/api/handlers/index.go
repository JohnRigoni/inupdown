package handlers

import (
	"embed"
	"fmt"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, assetsFS embed.FS) {
	f := r.URL.RequestURI()
	fmt.Println("indexHandler:", f)
	if f == "/" {
		// http.ServeFile(w, r, "./assets/index.html")
		p, err := assetsFS.ReadFile("assets/index.html")
		if err != nil {
			fmt.Println("index is err")
			return
		}
		w.Write(p)
		return
	}
	http.ServeFile(w, r, "."+f)
}
