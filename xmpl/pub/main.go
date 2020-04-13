package main

import (
	"gosaqws/gosaqws"
	"time"
)

func main() {
	gosaqws.Install("/saqws")

	var srv gosaqws.Server
	srv.Launch(9876)

	time.Sleep(time.Minute)
	srv.Shutdown()
}
