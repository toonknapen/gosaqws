package main

import (
	"gosaqws/gosaqws"
	"log"
)

func OnMessage(data []byte) {
	log.Println(string(data))
}

func main() {
	var sub gosaqws.Sub
	sub.OnMessageFn = OnMessage
	sub.ConnectSub("ws", "localhost", 9876, "/saqws")
}
