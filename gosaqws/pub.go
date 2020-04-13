package gosaqws

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Pub struct {
	m sync.RWMutex
	q [][]byte
}

var pub Pub

func Install(path string) {
	http.HandleFunc(path, HandleWs)
}

func NewSession() {
	pub.m.Lock()
	defer pub.m.Unlock()
	pub.q = nil
}
func Append(data []byte) {
	pub.m.RLock()
	defer pub.m.RUnlock()
	pub.q = append(pub.q, data)
}

var upgrader = websocket.Upgrader{}

func HandleWs(w http.ResponseWriter, r *http.Request) {
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer clientConn.Close()

	pub.m.RLock()
	clientConn.WriteMessage(websocket.TextMessage, []byte("hallo"))
	pub.m.RUnlock()

	clientConn.Close()
}
