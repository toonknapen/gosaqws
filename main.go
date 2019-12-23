package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)


func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:9876", Path: "/saqws"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Fatal error in dial", err)
	}
	defer c.Close()

	for i:=0 ; i < 300 ; i++ {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("readmessage error", err)
		}
		log.Printf("recv: %s", message)
	}
}
