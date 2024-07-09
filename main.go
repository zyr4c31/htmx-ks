package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/zyr4c31/htmx-ks/slideshow"
)

func createSlideShow() *slideshow.Slideshow {
	slideshow := slideshow.Slideshow{}
	slideshow.Add("introduction")
	slideshow.Add("ajax")
	slideshow.Add("triggers")
	slideshow.Add("trigger modifier")
	slideshow.Add("trigger filters")
	slideshow.Add("special events")
	slideshow.Add("polling")
	slideshow.Add("load polling")
	slideshow.Add("indicators")
	slideshow.Add("targets")
	slideshow.Add("swapping")
	slideshow.Add("synchronization")
	slideshow.Add("css transitions")
	slideshow.Add("out of band swaps")
	slideshow.Add("parameters")
	slideshow.Add("confirming")
	slideshow.Add("inheritance")
	slideshow.Add("boosting")
	slideshow.Add("animations")
	slideshow.Add("websockets & SSE")
	slideshow.Add("requests & responses")
	slideshow.Add("validation")
	slideshow.Add("extensions")
	slideshow.Add("events & logging")
	slideshow.Add("debugging")
	slideshow.Add("scripting")
	slideshow.Add("hx-on attribute")
	slideshow.Add("3rd party integration")
	slideshow.Add("Web Components")
	slideshow.Add("caching")
	slideshow.Add("security")
	slideshow.Add("configuring")
	return &slideshow
}

func main() {
	slideshow := createSlideShow()

	hostname, err := os.Hostname()
	if err != nil {
		log.Panic(err)
	}

	addr := fmt.Sprintf("%v:8080", hostname)
	sm := http.NewServeMux()

	fs := http.FileServer(http.Dir("assets/"))
	sp := http.StripPrefix("/assets/", fs)
	sm.Handle("/assets/", sp)

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
			unsafeLink := fmt.Sprintf("/%v", name)
			link := templ.SafeURL(unsafeLink)
			tocItem := tocItem{
				name: name,
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
