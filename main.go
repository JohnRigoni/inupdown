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

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Asset handled:", vars["file"])
	p, err := indexHtml.ReadFile("assets/" + vars["file"])
	if err != nil {
		fmt.Println("assets is err")
		return
	}
	w.Header().Add("Content-Type", "application/javascript")
	w.Write(p)
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
			DelLink: "/del?file=" + file.Name(),
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

	r.HandleFunc("/contact/1", contactHandler)
	r.HandleFunc("/assets", assetsHandler)
	r.HandleFunc("/upload", fileUploadHandler)
	// r.HandleFunc("/flist", fListHandler)
	r.HandleFunc("/flist", helloH)
	r.HandleFunc("/del", delHandler)
	// r.HandleFunc("/hello", helloH)

	r.PathPrefix("/").Handler(http.HandlerFunc(indexHandler))
	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

	_ = templ.NopComponent
}
