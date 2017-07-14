package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/udistrital/notificacion_api/models"
)

// WebsocketController operations for Websocket
type WebsocketController struct {
	beego.Controller
}

// URLMapping ...
func (c *WebsocketController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

//funciones para las notificaciones por websockets

// Join method handles WebSocket requests for WebsocketController.
func (this *WebsocketController) Join() {
	UserId := this.GetString("id")
	if len(UserId) == 0 {
		fmt.Println("falta el parametro id")
		return
	}
	Profiles := strings.Split(this.GetString("profiles"), "_")

	if len(Profiles) == 0 {
		fmt.Println("falta el parametro profile")
		return
	}

	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	// Join chat room.
	Join(UserId, Profiles, ws)
	fmt.Println("Pasa Join")
	defer Leave(UserId)

	// Message receive loop.
	/*for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		publish <- newEvent(models.EVENT_MESSAGE, UserId, string(p))
	}*/
}

// broadcastWebSocket broadcasts messages to WebSocket users.
func broadcastWebSocket(event models.Event) {
	var response models.Notificacion
	models.SendJson("10.20.0.254/configuracion_api/v1/notificacion", "POST", &response, &event.Content)
	data, err := json.Marshal(&response)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}

	//for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
	// Immediately send event to WebSocket users.
	var ws *websocket.Conn
	if event.Content.PerfilDestino == 0 {
		ws = connectionsId[event.User]
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				unsubscribe <- event.User
			}
		}
	} else {
		for _, prof := range connectionsProfile {
			ws = prof[event.User]
			if ws != nil {
				if ws.WriteMessage(websocket.TextMessage, data) != nil {
					// User disconnected.
					unsubscribe <- event.User
				}
			}
		}

	}

	//}
}

//++++++++++++++++++++++++++++++++++++++++++++++++
// Post ...
// @Title Create
// @Description create Websocket
// @Param	body		body 	models.Websocket	true		"body for Websocket content"
// @Success 201 {object} models.Websocket
// @Failure 403 body is empty
// @router / [post]
func (c *WebsocketController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Websocket by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Websocket
// @Failure 403 :id is empty
// @router /:id [get]
func (c *WebsocketController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Websocket
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Websocket
// @Failure 403
// @router / [get]
func (c *WebsocketController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Websocket
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Websocket	true		"body for Websocket content"
// @Success 200 {object} models.Websocket
// @Failure 403 :id is not int
// @router /:id [put]
func (c *WebsocketController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Websocket
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *WebsocketController) Delete() {

}
