package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func server(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
	}

	go handleConnection(c)
}

func handleConnection(c *websocket.Conn) {
	defer c.Close()
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", msg)
		// if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
		// 	log.Println("write:", err)
		// 	return
		// }
	}
}

func main() {
	http.HandleFunc("/", server)
	log.Println("Server started on :8080")
	_ = http.ListenAndServe(":8080", nil)
}
