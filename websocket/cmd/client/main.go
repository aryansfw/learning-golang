package main

import (
	"bufio"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	scanner := bufio.NewScanner(os.Stdin)

	u := url.URL{Scheme: "ws", Host: *addr}
	log.Printf("Connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	done := make(chan struct{})
	msg := make(chan string)

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	// input
	go func() {
		for {
			var input strings.Builder
			for {
				scanner.Scan()
				line := scanner.Text()
				if len(line) == 0 {
					break
				}
				input.WriteString(line)
				input.WriteString("\n")
			}
			msg <- input.String()
		}
	}()

	for {
		select {
		case <-done:
			return
		case m := <-msg:
			if err := conn.WriteMessage(websocket.TextMessage, []byte(m)); err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
