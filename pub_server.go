package gosaqws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Server is a small wrapper around http.Server
//
// This wrapper provides just convenience functions to launch and shutdown (something that is forgotten
// in many applications) the webserver easily.
type Server struct {
	srv *http.Server
}

// Launch returns a handle to the server
//
// The handle is needed to Shutdown the server later
func Launch(port int) (server Server) {
	addr := fmt.Sprintf(":%d", port)
	server.srv = &http.Server{Addr: addr}
	go func() {
		err := server.srv.ListenAndServe()
		if err != nil {
			log.Println("ERROR in http.Server.ListenAndServe:", err)
		}
	}()
	return server
}

// Launch returns a handle to the server
//
// The handle is needed to Shutdown the server later
func LaunchTLS(port int, crtFile string, keyFile string) (server Server) {
	addr := fmt.Sprintf(":%d", port)
	server.srv = &http.Server{Addr: addr}
	go func() {
		err := server.srv.ListenAndServeTLS(crtFile, keyFile)
		if err != nil {
			log.Println("ERROR in http.Server.ListenAndServeTLS:", err)
		}
	}()
	return server
}

func (server *Server) Shutdown() {
	if server.srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.srv.Shutdown(ctx)
		if err != nil {
			log.Printf("Error while shutting down webserver")
		}
		server.srv = nil
	}
}
