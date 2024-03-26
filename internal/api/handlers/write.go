package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func WriteFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusInternalServerError)
		return
	}
	for key, values := range r.Form {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	nm := r.Form.Get("filename")
	if len(nm) <= 4 || nm[len(nm)-4:] != ".txt" {
		nm += ".txt"
	}

	filePath := "./" + nm
	content := r.Form.Get("filecontents") + "\n"
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "PUT request processed successfully\n")
}
