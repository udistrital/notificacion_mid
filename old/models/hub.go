package models

import (
	"fmt"
)

type ConnValues struct{

	C *Connection
	Id string
	Profile []string

}

type SendingMessage struct{

	M []byte
	ConnValues ConnValues

}

type Hub struct {
	// Registered connections by id.
	connectionsId map[string]*Connection

	// Registered connections by profile.
	connectionsProfile map[string]map[*Connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Inbound messages to the connections by id.
	sendPersonalMessage chan SendingMessage

	// Inbound messages to the connections by the profile.
	sendProfileMessage chan SendingMessage

	// Register requests from the connections.
	register chan ConnValues

	// Unregister requests from connections.
	unregister chan ConnValues
}

func NewHub() Hub {
	return Hub{
		broadcast:   					make(chan []byte),
		sendPersonalMessage: 	make(chan SendingMessage),
		sendProfileMessage: 	make(chan SendingMessage),
		register:    					make(chan ConnValues),
		unregister:  					make(chan ConnValues),
		connectionsId: 				make(map[string]*Connection),
		connectionsProfile:	  make(map[string]map[*Connection]bool),
	}
}

// Adds a connection to the connection map
func (h *Hub) Register(connValues ConnValues) {
	h.register <- connValues
}

// Removes a connection from the connection map
func (h *Hub) Unregister(connValues ConnValues) {
	h.unregister <- connValues
}

func (h *Hub) SendPersonalMessage(sendingMessage SendingMessage){
	h.sendPersonalMessage <- sendingMessage
}

func (h *Hub) SendProfileMessage(sendingMessage SendingMessage){
	h.sendProfileMessage <- sendingMessage
}

// Hub's main loop handles commands for the connection map
func (h *Hub) Run() {
	for {
		select {
		// Adds a connection
		case connValues := <-h.register:
			fmt.Println("Connect")
			h.connectionsId[connValues.Id] = connValues.C
			for _,profile := range connValues.Profile {
				if _, ok := h.connectionsProfile[profile]; ok {
					(h.connectionsProfile[profile])[connValues.C] = true
				}else{
					h.connectionsProfile[profile] = make(map[*Connection]bool)
					(h.connectionsProfile[profile])[connValues.C] = true
				}
				fmt.Printf("Register user: %s\n",connValues.Id)
			}

		// Removes a connection
		case connValues := <-h.unregister:
			fmt.Println("Disconnect")
			c := h.connectionsId[connValues.Id]
			delete(h.connectionsId, connValues.Id)
			for _,profile := range connValues.Profile {
				delete(h.connectionsProfile[profile],c)
			}
			close(c.send)
		// Sends a mesage to each connected client
		// case m := <-h.broadcast:
		// 	fmt.Printf("Broadcasting: %s\n", m)
		// 	for c := range h.connections {
		// 		select {
		// 		case c.send <- m:
		// 		default:
		// 			//delete(h.connections, c)
		// 			close(c.send)
		// 			go c.ws.Close()
		// 		}
		// 	}
		case sendingMessage := <-h.sendPersonalMessage:
			fmt.Printf("Sending message: %s\n",sendingMessage.M)
			if c, ok := h.connectionsId[sendingMessage.ConnValues.Id]; ok{
				select {
					case c.send <- sendingMessage.M:
					default:
						//delete(h.connections, c)
						// close(c.send)
						// go c.ws.Close()
				}
			}

		case sendingMessage := <-h.sendProfileMessage:
			fmt.Printf("Sending message: %s\n",sendingMessage.M)
			for _,profile := range sendingMessage.ConnValues.Profile{
				for c,_ := range h.connectionsProfile[profile]{
					select {
						case c.send <- sendingMessage.M:
						default:
							//delete(h.connections, c)
							// close(c.send)
							// go c.ws.Close()
					}
				}
			}

		}
	}
}
