package server

import (
	"embed"
	"fmt"
	"log"
	"net/http"
)

func Start(assetsFS embed.FS) error {
	router := Router(assetsFS)
	url := "0.0.0.0:8088"
	fmt.Println("Server listening on:", url)
	log.Fatal(http.ListenAndServe(url, router))

	return nil
}
