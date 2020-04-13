package gosaqws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	shutdown bool
	srv      *http.Server
}

func (server Server) Launch(port int) {
	addr := fmt.Sprintf(":%d", port)
	server.srv = &http.Server{Addr: addr}
	go server.srv.ListenAndServe()

	for !server.shutdown {
		time.Sleep(time.Second)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.srv.Shutdown(ctx)
	if err != nil {
		log.Printf("Error while shutting down webserver")
	}

	server.shutdown = false
}

func (server Server) LaunchTLS(port int, crtFile string, keyFile string) {
	addr := fmt.Sprintf(":%d", port)
	server.srv = &http.Server{Addr: addr}
	go server.srv.ListenAndServeTLS(crtFile, keyFile)

	for !server.shutdown {
		time.Sleep(time.Second)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.srv.Shutdown(ctx)
	if err != nil {
		log.Printf("Error while shutting down webserver")
	}

	server.shutdown = false
}

func (server Server) Shutdown() {
	server.shutdown = true

	for server.shutdown {
		time.Sleep(time.Second)
	}
}
