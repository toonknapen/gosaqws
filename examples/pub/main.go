package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/toonknapen/gosaqws"
)

func main() {
	gosaqws.Install("/saqws")

	srv := gosaqws.Launch(9876)
	log.Println("launched")

	numSessions := 50
	for sessionId := 0; sessionId < numSessions; sessionId++ {
		gosaqws.NewSession()

		events := []string{"one", "two", "three", "four", "five", "six"}
		numEvents := len(events)
		for eventId := 0; eventId < numEvents; eventId++ {
			data, _ := json.Marshal(events[eventId])
			log.Println("Appending", events[eventId])
			gosaqws.Append(data)
			time.Sleep(time.Second)
		}
	}

	log.Println("shutting down")
	srv.Shutdown()
}
