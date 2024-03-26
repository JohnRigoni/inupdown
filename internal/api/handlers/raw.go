package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
)


func RawHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("rawHandler:", r.URL.RequestURI())
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parrse form", http.StatusInternalServerError)
		return
	}

	filename := r.Form.Get("file")
	if filename == "" {
		http.Error(w, "Query: 'file' missing", http.StatusInternalServerError)
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	contentbytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(contentbytes))
}
