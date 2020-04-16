package main

import (
	gosaqws2 "gosaqws"
	"log"
	"time"
)

func main() {
	gosaqws2.Install("/saqws")

	var srv gosaqws2.Server
	srv.Launch(9876)
	log.Println("launched")

	numSessions := 50
	for sessionId := 0; sessionId < numSessions; sessionId++ {
		gosaqws2.NewSession()

		events := []string{"one", "two", "three", "four", "five", "six"}
		numEvents := len(events)
		for eventId := 0; eventId < numEvents; eventId++ {
			gosaqws2.Append([]byte(events[eventId]))
			time.Sleep(time.Second)
		}
	}

	log.Println("shutting down")
	srv.Shutdown()
}
