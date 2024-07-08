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

	dir, files, err := readFiles()
	if err != nil {
		log.Panic(err)
	}

	err = generateRoutes(dir, files, sm)
	if err != nil {
		log.Panic(err)
	}

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
	err = http.ListenAndServe(addr, sm)
	if err != nil {
		log.Panic(err)
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

func generateRoutes(dir fs.FS, files []fs.DirEntry, sm *http.ServeMux) error {
	var err error
	for _, file := range files {
		fileName := strings.Split(file.Name(), "-")[1]
		path := fmt.Sprintf("/%v", fileName)
		sm.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			bytes, _ := fs.ReadFile(dir, file.Name())
			ext := strings.Split(file.Name(), ".")[1]
			if ext == "html" {
				w.Header().Set("Content-Type", "text/html")
			}
			_, err := w.Write(bytes)
			if err != nil {
				log.Panic(err)
			}

		})
	}
	return err
}

func htmlRenderer(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("<h1>%v</h1>", text)))
}
