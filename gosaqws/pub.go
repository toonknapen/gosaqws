package gosaqws

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type Pub struct {
	m         sync.RWMutex
	sessionId int32
	q         [][]byte
}

var pub Pub

func Install(path string) {
	http.HandleFunc(path, HandleSAQWS)
}

func NewSession() {
	pub.m.Lock()
	defer pub.m.Unlock()
	pub.sessionId++
	pub.q = nil
}

func Append(data []byte) {
	pub.m.RLock()
	defer pub.m.RUnlock()
	pub.q = append(pub.q, data)
}

var upgrader = websocket.Upgrader{}

func HandleSAQWS(w http.ResponseWriter, r *http.Request) {
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer clientConn.Close()

	// catch up on all that is already available in the session
	pub.m.RLock()
	sessionId := pub.sessionId
	cursor, err := sendCurrentSessionBacklog(clientConn, 0)
	pub.m.RUnlock()

	// start listening to updates (using a spin-lock currently, will be improved)
	for true {
		time.Sleep(10 * time.Millisecond)

		pub.m.RLock()

		// check if in the mean time no new session has started, if so reset the cursor
		if sessionId < pub.sessionId {
			cursor = 0
			sessionId = pub.sessionId
		}

		numMsgAvailable := len(pub.q)
		if cursor < numMsgAvailable {
			for ; cursor < numMsgAvailable; cursor++ {
				err = clientConn.WriteMessage(websocket.TextMessage, pub.q[cursor])
			}
		}

		pub.m.RUnlock()
	}
}

func sendCurrentSessionBacklog(clientConn *websocket.Conn, cursorStart int) (cursor int, err error) {
	numMsgTotal := len(pub.q)
	for ; cursor < numMsgTotal; cursor++ {
		log.Println(string(pub.q[cursor]))
		err = clientConn.WriteMessage(websocket.TextMessage, pub.q[cursor])
	}
	return cursor, err
}
