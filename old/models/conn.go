package models

import (
	"github.com/gorilla/websocket"
	// "strings"
)

type Connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func NewConnection(ws *websocket.Conn) *Connection {
	return &Connection{send: make(chan []byte, 256), ws: ws}
}

// Read loop for each connection reads a message and broadcasts it
func (c *Connection) Reader(h Hub) {
	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		// values := strings.Split(string(message), "::")
		// profiles := strings.Split(values[2], ",")
		// h.SendPersonalMessage(SendingMessage{M:[]byte(values[0]), ConnValues: ConnValues{C:nil, Id:values[1], Profile:profiles}})
	}
	c.ws.Close()
}

// Write loop for each connection writes whatever comes across the send channel
func (c *Connection) Writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}
