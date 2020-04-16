package main

import (
	"encoding/json"
	"log"

	"github.com/toonknapen/gosaqws"
)

func OnMessage(data json.RawMessage) {
	var msg string
	json.Unmarshal(data, &msg)
	log.Println(msg)
}

func main() {
	var sub gosaqws.Sub
	sub.OnMessageFn = OnMessage
	sub.ConnectSub("ws", "localhost", 9876, "/saqws")
}
