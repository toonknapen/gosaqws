package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

type PassMark struct {
	Lap   int `json:"lap"`
	Time  int `json:"time"`
	Total int `json:"total"`
}

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Fatal error in dial", err)
	}
	defer c.Close()

	for i := 0; i < 300; i++ {
		var hello = []byte{'h', 'e', 'l'}
		werr := c.WriteMessage(websocket.TextMessage, hello)
		if werr != nil {
			log.Println(werr)
			return
		}
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("readmessage error", err)
		}
		log.Printf("recv: %s", message)

		//var passmarks []PassMark
		//err = json.Unmarshal(message, &passmarks)
		//if err != nil {
		//	panic(err)
		//}
		//log.Printf("f{passmarks}")
	}
}
