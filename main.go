package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Panic(err)
	}

	addr := fmt.Sprintf("%v:8080", hostname)
	sm := http.NewServeMux()

	fs := http.FileServer(http.Dir("assets"))
	sp := http.StripPrefix("/assets", fs)
	sm.Handle("/assets", sp)

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
			name := strings.Split(fileName, ".")[0]
			link := templ.SafeURL(fmt.Sprintf("/%v", name))
			tocItem := tocItem{
				name: fileName,
				link: link,
			}
			model.items = append(model.items, tocItem)
		}
		tmpl := content("toc", model)
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
		fileNameExtSplit := strings.Split(fileName, ".")
		name := fileNameExtSplit[0]
		ext := fileNameExtSplit[1]
		path := fmt.Sprintf("/%v", name)
		sm.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			bytes, _ := fs.ReadFile(dir, file.Name())
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
