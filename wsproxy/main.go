package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"time"
)

var serv = "localhost:8080"
var downstreamScheme = "ws"
var downstreamHost = "localhost:9876"
var mapPaths = map[string]string{
	"/saqws": "/saqws",
}
var upgrader = websocket.Upgrader{}

func proxyToClient(clientConn *websocket.Conn, serverConn *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		_, message, err := serverConn.ReadMessage()
		if err != nil {
			log.Println("read error", err)
			return
		}
		log.Printf("recv: %s", message)
		if err := clientConn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func proxyToServer(done chan struct{}) {
	for i := 0; i < 1000; i++ {
		time.Sleep(1 * time.Second)
		log.Println("proxyToServer slept more")
		select {
		case <-done:
			log.Println("Shutting down proxyToServer")
			return
		default:
			// just continue
		}
	}
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer clientConn.Close()

	// make a wsconnection to the downstream
	serverURL := url.URL{Scheme: downstreamScheme, Host: downstreamHost, Path: "/saqws"}
	serverConn, _, err := websocket.DefaultDialer.Dial(serverURL.String(), nil)
	if err != nil {
		log.Fatal("Dial", err)
	}
	defer serverConn.Close()

	chanToClient := make(chan struct{})
	chanToServer := make(chan struct{})
	go proxyToClient(clientConn, serverConn, chanToClient)
	go proxyToServer(chanToServer)

	select {
	case <-chanToClient:
		log.Println("server disconnected")
	}
	log.Println("about done")
	close(chanToServer)
	time.Sleep(1 * time.Second)
}

func main() {
	for k, _ := range mapPaths {
		http.HandleFunc(k, handleWS)
	}
	log.Fatal(http.ListenAndServe(serv, nil))
}
