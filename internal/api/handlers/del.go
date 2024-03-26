package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func DelHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	fmt.Println("got file: ", r.Form["file"])
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "PUT request processed successfully\n")
}
