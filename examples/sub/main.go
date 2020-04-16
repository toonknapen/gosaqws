package main

import (
	gosaqws2 "gosaqws"
	"log"
)

func OnMessage(data []byte) {
	log.Println(string(data))
}

func main() {
	var sub gosaqws2.Sub
	sub.OnMessageFn = OnMessage
	sub.ConnectSub("ws", "localhost", 9876, "/saqws")
}
