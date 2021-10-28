package gosaqws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// PubStub needs to be used by the publisher to store the messages
// and serve the messages to each subscriber
type PubStub struct {
	m         sync.RWMutex
	sessionId int
	q         []json.RawMessage
}

// Install the handler for the publisher at the specified path in the web-server
//
// Instead of calling Install, the `HandleSAQWS` can also be installed directly
func Install(path string, p *PubStub) {
	http.HandleFunc(path, p.HandleSAQWS)
}

// NewSession starts a new session for the publisher
//
// All messages that were previously appended are deleted.
func (p *PubStub) NewSession() {
	p.m.Lock()
	defer p.m.Unlock()
	p.sessionId++
	p.q = nil
}

// Append a range of bytes to the queue of messages to be published
func (p *PubStub) Append(data json.RawMessage) {
	p.m.Lock()
	defer p.m.Unlock()
	p.q = append(p.q, data)
}

var upgrader = websocket.Upgrader{}

func (p *PubStub) HandleSAQWS(w http.ResponseWriter, r *http.Request) {
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
	p.m.RLock()
	sessionId := p.sessionId
	cursor, err := p.sendCurrentSessionBacklog(clientConn, 0)
	p.m.RUnlock()

	// start listening to updates (using a spin-lock currently, will be improved)
	for true {
		time.Sleep(10 * time.Millisecond)

		p.m.RLock()

		// check if in the mean time no new session has started, if so reset the cursor
		if sessionId < p.sessionId {
			cursor = 0
			sessionId = p.sessionId
		}

		cursor, err = p.sendCurrentSessionBacklog(clientConn, cursor)
		p.m.RUnlock()
	}
}

func (p *PubStub) sendCurrentSessionBacklog(clientConn *websocket.Conn, cursorStart int) (cursor int, err error) {
	cursor = len(p.q)
	if cursorStart < cursor {
		raw, err := json.Marshal(p.q[cursorStart:cursor])
		if err != nil {
			log.Println("ERROR:marshalling slice of q:", err)
			return cursorStart, err
		}
		err = clientConn.WriteMessage(websocket.TextMessage, raw)
	}
	return cursor, err
}
