package gosaqws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

type Sub struct {
	OnMessageFn func(msg []byte)
	shutdown    bool
	conn        *websocket.Conn
}

func (sub *Sub) ConnectSub(scheme string, host string, port int, path string) {
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

	for !sub.shutdown {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("error with ws connection to conn:", err)
			return
		}

		if sub.OnMessageFn != nil {
			sub.OnMessageFn(data)
		}
	}
}

func (sub *Sub) Disconnect() {
	sub.shutdown = true
	sub.conn.Close()
}
