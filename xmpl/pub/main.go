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

	time.Sleep(10 * time.Second)
	log.Println("shutting down")
	srv.Shutdown()
}
