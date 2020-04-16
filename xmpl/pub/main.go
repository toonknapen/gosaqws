package main

import (
	"gosaqws/gosaqws"
	"log"
	"time"
)

func main() {
	gosaqws.Install("/saqws")

	var srv gosaqws.Server
	srv.Launch(9876)
	log.Println("launched")

	numSessions := 50
	for sessionId := 0; sessionId < numSessions; sessionId++ {
		gosaqws.NewSession()

		events := []string{"one", "two", "three", "four", "five", "six"}
		numEvents := len(events)
		for eventId := 0; eventId < numEvents; eventId++ {
			gosaqws.Append([]byte(events[eventId]))
			time.Sleep(time.Second)
		}
	}

	log.Println("shutting down")
	srv.Shutdown()
}
