package handlers

import (
	"fmt"
	"net/http"
)

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("apiHandler:", r.URL.RequestURI())
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form: apiHandler", http.StatusInternalServerError)
		return
	}

	apiRoute := r.Form.Get("api")
	if apiRoute == "" {
		http.Error(w, "Query: 'api' missing", http.StatusInternalServerError)
		return
	}

	switch apiRoute {
	case "upload":
		fmt.Println("uploadddd-api")
		fileUploadHandler(w, r)
	case "flist":
		fmt.Println("flist-api")
		flistHandler(w, r)
	case "delete":
		fmt.Println("delete-api")
		delHandler(w, r)
	case "writef":
		fmt.Println("writef-api")
		writeFileHandler(w, r)
	case "raw":
		fmt.Println("raw-api")
		rawHandler(w, r)
	default:
		http.Error(w, "api: "+apiRoute+" not found", http.StatusInternalServerError)
	}
}
