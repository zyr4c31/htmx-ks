package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func htmlRenderer(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("<h1>%v</h1>", text)))
}

func main() {
	addr := ":8080"
	sm := http.NewServeMux()
	sm.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		pathID := r.PathValue("id")
		intID, _ := strconv.Atoi(pathID)
		templ := content(intID)
		templ.Render(r.Context(), w)
	})
	http.ListenAndServe(addr, sm)
}
