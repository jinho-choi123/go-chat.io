package safeWebsocket

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type SafeWebsocket struct {
	Conn *websocket.Conn
	Send chan []byte
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func NewSafeWebsocket(conn *websocket.Conn) *SafeWebsocket {
	return &SafeWebsocket{
		Conn: conn,
		Send: make(chan []byte),
	}
}

func (sws *SafeWebsocket) Read(publish_chan chan []byte, unregister_chan chan *SafeWebsocket) {
	defer func() {
		sws.Conn.Close()
		unregister_chan <- sws
	}()

	sws.Conn.SetReadLimit(maxMessageSize)
	sws.Conn.SetReadDeadline(time.Now().Add(pongWait))
	sws.Conn.SetPongHandler(func(string) error { sws.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := sws.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("%v", err.Error())
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// send the message to hub manager
		publish_chan <- message
	}
}

// hub manager로부터 broadcasting 된 message들을 websocket에 쓴다.
// 그리고, 주기적으로 client에게 ping message을 보낸다.
func (sws *SafeWebsocket) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		sws.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-sws.Send:
			sws.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				sws.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := sws.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(sws.Send)

			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-sws.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			sws.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := sws.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
