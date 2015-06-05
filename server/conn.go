package server

import (
	"time"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/gorilla/websocket"
)

type conn struct {
	conn *websocket.Conn
	in   chan WSRequest
	out  chan []byte
}

func newConn(ws *websocket.Conn) *conn {
	return &conn{
		conn: ws,
		in:   make(chan WSRequest, 32),
		out:  make(chan []byte, 32),
	}
}

func (c *conn) send() {
	var err error
	ping := time.Tick(20 * time.Second)

	for {
		select {
		case msg, ok := <-c.out:
			if !ok {
				return
			}

			err = c.conn.WriteMessage(websocket.TextMessage, msg)

		case <-ping:
			err = c.conn.WriteJSON(WSResponse{Type: "ping"})
		}

		if err != nil {
			return
		}
	}
}

func (c *conn) recv() {
	var req WSRequest

	for {
		err := c.conn.ReadJSON(&req)
		if err != nil {
			close(c.in)
			return
		}

		c.in <- req
	}
}

func (c *conn) close() {
	close(c.out)
}