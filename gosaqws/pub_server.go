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

func (server *Server) Launch(port int) {
	addr := fmt.Sprintf(":%d", port)
	server.srv = &http.Server{Addr: addr}
	go func() {
		err := server.srv.ListenAndServe()
		if err != nil {
			log.Println("ERROR in http.Server.ListenAndServe:", err)
		}
	}()
}

func (server *Server) LaunchTLS(port int, crtFile string, keyFile string) {
	addr := fmt.Sprintf(":%d", port)
	server.srv = &http.Server{Addr: addr}
	go func() {
		err := server.srv.ListenAndServeTLS(crtFile, keyFile)
		if err != nil {
			log.Println("ERROR in http.Server.ListenAndServeTLS:", err)
		}
	}()

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
