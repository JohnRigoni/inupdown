package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"inupdown/internal/api/templates"

	"github.com/a-h/templ"
)

func FlistHandler(w http.ResponseWriter, r *http.Request) {
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
			RawLink:  templ.SafeURL("/api/raw?internal=true&file=" + file.Name()),
			DownLink: templ.SafeURL("/" + file.Name()),
			DelLink:  "/api/delete?internal=true&file=" + file.Name(),
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
