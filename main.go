package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

func htmlRenderer(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("<h1>%v</h1>", text)))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var fullChatMsg []string

func main() {
	addr := "zarch-mllrlt:8080"
	sm := http.NewServeMux()
	go sm.HandleFunc("/chatroom", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Fatalln(err)
			}
			stringedMsg := string(msg)
			chatMsg := strings.Split(stringedMsg, `"`)
			fullChatMsg = append(fullChatMsg, chatMsg[3])
			fmt.Printf("addr: %s msg: %s \n", conn.RemoteAddr(), chatMsg)
			w.Header().Set("HX-Reswap", "afterbegin")
			err = conn.WriteMessage(msgType, []byte(fmt.Sprintf(`<div id="chat_room" hx-swap-oob="beforeend">%v</div>`, fullChatMsg)))
			if err != nil {
				log.Fatalln(err)
			}
		}
	})
	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pathID := r.PathValue("id")
		intID, _ := strconv.Atoi(pathID)
		templ := content(intID, fullChatMsg)
		templ.Render(r.Context(), w)
	})
	http.ListenAndServe(addr, sm)
}
