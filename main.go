package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"updowninserve/templates"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

//go:embed assets
var indexHtml embed.FS

func indexHandler(w http.ResponseWriter, r *http.Request) {
	f := r.URL.RequestURI()
	fmt.Println("Root handled:", f)
	if f == "/" {
		http.ServeFile(w, r, "./assets/index.html")
		// p, err := indexHtml.ReadFile("assets/index.html")
		// if err != nil {
		// 	fmt.Println("index is err")
		// 	return
		// }
		// w.Write(p)
		// return
	}
	http.ServeFile(w, r, "."+f)
}

func indexHandler2(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Root2 handled2:", r.URL.RequestURI())
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusInternalServerError)
		return
	}

	asset := r.Form.Get("assets")
	apiRoute := r.Form.Get("api")
	if asset == "" && apiRoute == "" {
		http.Error(w, "Both assets and api missing", http.StatusInternalServerError)
		return
	}

	if asset != "" {
		file := fmt.Sprint("./assets/", asset)
		http.ServeFile(w, r, file)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("apiHandler:", r.URL.RequestURI())
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form: apiHandler", http.StatusInternalServerError)
		return
	}

	apiRoute := r.Form.Get("api")
	if  apiRoute == "" {
		http.Error(w, "api missing", http.StatusInternalServerError)
		return
	}

	switch apiRoute {
	case "upload":
		fmt.Println("uploadddd")
		fileUploadHandler(w, r)
	case "flist":
		fmt.Println("flist")
		helloH(w, r)
	case "delete":
		fmt.Println("delete")
		delHandler(w, r)
	case "writef":
		fmt.Println("writef")
		contactHandler(w, r)

	}

}



func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["file"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Println("Failed to open uploaded file:", err)
			continue
		}
		defer file.Close()

		destFilePath := filepath.Join("./data", fileHeader.Filename)
		destFile, err := os.Create(destFilePath)
		if err != nil {
			fmt.Println("Failed to create destination file:", err)
			continue
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, file)
		if err != nil {
			fmt.Println("Failed to save file:", err)
			continue
		}
		fmt.Printf("Uploaded File: %s\n", fileHeader.Filename)
	}

	if len(files) == 0 {
		fmt.Println("no files. ", r.MultipartForm.File)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File(s) uploaded successfully\n")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
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

	nm := r.Form.Get("name")
	if len(nm) <= 4 || nm[len(nm)-4:] != ".txt" {
		nm += ".txt"
	}

	filePath := "./data/" + nm
	content := r.Form.Get("content") + "\n"
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "PUT request processed successfully\n")
}

func helloH(w http.ResponseWriter, r *http.Request) {
	dir := "./data"
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}
	items := []templates.FtblRowsS{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		finfo, _ := file.Info()

		item := templates.FtblRowsS{
			Link:    templ.SafeURL("/" + file.Name()),
			DelLink: "/?api=delete&file=" + file.Name(),
			Name:    file.Name(),
			Date:    finfo.ModTime().String(),
			Size:    strconv.FormatInt(finfo.Size(), 10),
		}
		items = append(items, item)

	}
	comp := templates.Hello(items)
	comp.Render(context.Background(), w)
}

func delHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	fmt.Println("got file: ", r.Form["file"])
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "PUT request processed successfully\n")
}

func main() {

	r := mux.NewRouter()

	r.PathPrefix("/").Queries("assets", "{assets}").Handler(http.HandlerFunc(indexHandler2))
	r.PathPrefix("/").Queries("api", "{api}").Handler(http.HandlerFunc(apiHandler))
	r.PathPrefix("/").Handler(http.HandlerFunc(indexHandler))

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

	_ = templ.NopComponent
}
