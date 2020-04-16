package gosaqws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	srv *http.Server
}

func (server *Server) Launch(port int) {
	addr := fmt.Sprintf(":%d", port)
	server.srv = &http.Server{Addr: addr}
	go server.srv.ListenAndServe()
}

func (server *Server) LaunchTLS(port int, crtFile string, keyFile string) {
	addr := fmt.Sprintf(":%d", port)
	server.srv = &http.Server{Addr: addr}
	go server.srv.ListenAndServeTLS(crtFile, keyFile)
}

func (server *Server) Shutdown() {
	if server.srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.srv.Shutdown(ctx)
		if err != nil {
			log.Printf("Error while shutting down webserver")
		}
	}
}
