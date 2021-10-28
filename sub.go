package gosaqws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// SubStub connects to the publisher at the given address and call the
// OnMessageFn callback for every message received
type SubStub struct {
	OnMessageFn func(msg json.RawMessage)
	shutdown    bool
	conn        *websocket.Conn
}

func (s *SubStub) ConnectSub(scheme string, host string, port int, path string) {
	hostPort := fmt.Sprintf("%s:%d", host, port)
	pubUrl := url.URL{Scheme: scheme, Host: hostPort, Path: path}

	var conn *websocket.Conn
	connected := false
	for !connected {
		var err error
		conn, _, err = websocket.DefaultDialer.Dial(pubUrl.String(), nil)
		if err != nil {
			log.Println("Failed connecting to", pubUrl, "trying again in 1 s.:", err)
			time.Sleep(time.Second)
		} else {
			connected = true
		}
	}

	defer func() {
		err := conn.Close()
		log.Println("ERROR: Failed closing the ws to the connector but was going to stop listening anyway", err)
	}()

	for !s.shutdown {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("error with ws connection to conn:", err)
			return
		}

		var events []json.RawMessage
		err = json.Unmarshal(data, &events)

		if s.OnMessageFn != nil {
			for _, data := range events {
				s.OnMessageFn(data)
			}
		}
	}
}

func (s *SubStub) Disconnect() {
	s.shutdown = true
	s.conn.Close()
}
