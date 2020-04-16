package gosaqws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type Pub struct {
	m         sync.RWMutex
	sessionId int
	q         [][]byte
}

var pub Pub

// Install the handler for the publisher at the specified path in the web-server
//
// Instead of calling Install, the `HandleSAQWS` can also be installed directly
func Install(path string) {
	http.HandleFunc(path, HandleSAQWS)
}

// NewSession starts a new session for the publisher
//
// All messages that were previously appended are deleted.
func NewSession() {
	pub.m.Lock()
	defer pub.m.Unlock()
	pub.sessionId++
	pub.q = nil
}

// Append a range of bytes to the queue of messages to be published
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
	defer func() {
		err := clientConn.Close()
		if err != nil {
			log.Println("ERROR while closing websocket:", err)
		}
	}()

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

		cursor, err = sendCurrentSessionBacklog(clientConn, cursor)
		pub.m.RUnlock()
	}
}

func sendCurrentSessionBacklog(clientConn *websocket.Conn, cursorStart int) (cursor int, err error) {
	cursor = len(pub.q)
	if cursorStart < cursor {
		raw, err := json.Marshal(pub.q[cursorStart:cursor])
		if err != nil {
			log.Println("ERROR:marshalling slice of q:", err)
			return cursorStart, err
		}
		err = clientConn.WriteMessage(websocket.TextMessage, raw)
	}
	return cursor, err
}
