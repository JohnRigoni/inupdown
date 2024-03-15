package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	f := r.URL.RequestURI()
	if f == "/" {
		f = "/index.html"
	}
	http.ServeFile(w, r, "."+f)
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
	// Get the uploaded files
	files := r.MultipartForm.File["file"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Println("Failed to open uploaded file:", err)
			continue
		}
		defer file.Close()
		// Create a destination file in the "./data" directory
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

func fListHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("flist")
	dir := "./data" // Directory path
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	outp := `<div class="filebox">`
	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		// fmt.Println(file.Info())	

		outp += `<a href="/data/` + file.Name() + `" rel="nofollow" download> ` + file.Name() + `</a><br>`
	}

	outp += "</div>"
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, outp)

}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data from the request
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

func main() {
	http.HandleFunc("/contact/1", contactHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", fileUploadHandler)
	http.HandleFunc("/flist", fListHandler)

	fmt.Println("Server listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
