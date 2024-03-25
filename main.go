package main

import (
	"embed"
	"fmt"
	"inupdown/templates"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

//go:embed assets
var assetsFS embed.FS

func indexHandler(w http.ResponseWriter, r *http.Request) {
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

func assetsHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "text/" + atype)
	w.Write(p)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
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

		destFilePath := filepath.Join("./", fileHeader.Filename)
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

func writeFileHandler(w http.ResponseWriter, r *http.Request) {
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

func flistHandler(w http.ResponseWriter, r *http.Request) {
	dir := "./"
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
		size := computeByteString(finfo.Size())

		item := templates.FtblRowsS{
			RawLink:  templ.SafeURL("/?api=raw&file=" + file.Name()),
			DownLink: templ.SafeURL("/" + file.Name()),
			DelLink:  "/?api=delete&file=" + file.Name(),
			Name:     file.Name(),
			Date:     finfo.ModTime().Format("Jan 2"),
			Size:     size,
		}
		items = append(items, item)
	}
	fmt.Println("got items: ", items)

	comp := templates.Hello(items)
	if comp == nil {
		fmt.Println("comp is nil")
		return
	}

	err = comp.Render(r.Context(), w)
	if err != nil {
		log.Fatalf("Failed to render template: %v", err)
	}
}

func computeByteString(bytes int64) string {
	output := ""
	suffixes := []string{"B", "KB", "MB", "GB", "TB"}

	for _, s := range suffixes {
		if bytes < 1000 {
			output = fmt.Sprintf("%d%s", bytes, s)
			return output
		}
		bytes = bytes / 1000
	}

	return fmt.Sprintf("%d%s", bytes, "PB")
}

func rawHandler(w http.ResponseWriter, r *http.Request) {
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
	r.PathPrefix("/").Queries("assets", "{assets}").Handler(http.HandlerFunc(assetsHandler))
	r.PathPrefix("/").Queries("api", "{api}").Handler(http.HandlerFunc(apiHandler))
	r.PathPrefix("/").Handler(http.HandlerFunc(indexHandler))

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8088", r))
}
