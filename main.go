package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	addr := "zarch-mllrlt:8080"
	sm := http.NewServeMux()

	files := generateRoutes(sm)

	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var model tocModel
		for _, file := range files {
			fileName := strings.Split(file.Name(), "-")[1]
			tocItem := tocItem{
				name: fileName,
				link: fmt.Sprintf("/%v", fileName),
			}
			model.items = append(model.items, tocItem)
		}
		tmpl := toc(model)
		w.Header().Set("Content-Type", "text/html")
		tmpl.Render(r.Context(), w)
	})

	log.Printf("running on http://%v", addr)
	err := http.ListenAndServe(addr, sm)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}

func readFiles() (fs.FS, []fs.DirEntry, error) {
	var files []fs.DirEntry
	dir := os.DirFS("./pages")
	files, err := fs.ReadDir(dir, ".")
	if err != nil {
		return dir, files, err
	}

	return dir, files, err
}

func generateRoutes(sm *http.ServeMux) []fs.DirEntry {
	dir, files, err := readFiles()
	for _, file := range files {
		path := fmt.Sprintf("/%v", file.Name())
		sm.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			bytes, _ := fs.ReadFile(dir, file.Name())
			ext := strings.Split(file.Name(), ".")[1]
			if ext == "html" {
				w.Header().Set("Content-Type", "text/html")
			}
			w.Write(bytes)
		})
	}
	if err != nil {
		log.Fatalf("cannot generate toc: %v", err)
	}
	return files
}

func htmlRenderer(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("<h1>%v</h1>", text)))
}
