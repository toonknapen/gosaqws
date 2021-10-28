package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/toonknapen/gosaqws"
)

func main() {
	p := gosaqws.PubStub{}
	gosaqws.Install("/saqws", &p)

	srv := gosaqws.Launch(9876)
	log.Println("launched")

	numSessions := 50
	for sessionId := 0; sessionId < numSessions; sessionId++ {
		p.NewSession()

		events := []string{"one", "two", "three", "four", "five", "six"}
		numEvents := len(events)
		for eventId := 0; eventId < numEvents; eventId++ {
			data, _ := json.Marshal(events[eventId])
			log.Println("Appending", events[eventId])
			p.Append(data)
			time.Sleep(time.Second)
		}
	}

	log.Println("shutting down")
	srv.Shutdown()
}
