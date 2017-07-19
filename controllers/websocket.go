package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/udistrital/notificacion_api/models"
	"github.com/udistrital/notificacion_api/utilidades"
)

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	models.BaseController
}

// Join method handles WebSocket requests for WebSocketController.
func (this *WebSocketController) Join() {
	Id := this.GetString("id")
	if len(Id) == 0 {
		beego.Info("Cannot get User Id")
		return
	}
	Profiles := strings.Split(this.GetString("profiles"), ",")
	if len(Id) == 0 {
		beego.Info("Cannot get User Id")
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
	Join(Id, Profiles, ws)
	defer Leave(Id)

	// Message receive loop.
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		publish <- newEvent(models.EVENT_MESSAGE, Id, Profiles, string(p))
	}
}

// broadcastWebSocket broadcasts messages to WebSocket users.
func broadcastWebSocket(event models.Event) {
	data, err := json.Marshal(event)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}
	if connectionsId[event.User] != nil {
		ws := connectionsId[event.User]
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				unsubscribe <- event.User
			}
		}
	}
	for _, value := range event.Profiles {
		fmt.Println("message from ", event.User)
		if connectionsProfile[value] != nil {
			for _, con := range connectionsProfile[value] {
				ws := con
				if ws != nil {
					if ws.WriteMessage(websocket.TextMessage, data) != nil {
						// User disconnected.
						unsubscribe <- event.User
					}
				}
			}

		}
	}

}

// broadcastWebSocket broadcasts messages to WebSocket users
func (this *WebSocketController) PushNotificacion() {
	fmt.Println("entro")
	var v map[string]interface{}
	//UserId := c.GetString("id")
	//fmt.Println("Id ", UserId)
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &v); err == nil {
		//push notificacion-------
		var perfil int64
		var usuario int64
		var cuerpo string
		err = utilidades.FillStruct(v["PerfilDestino"], &perfil)
		err = utilidades.FillStruct(v["UsuarioDestino"], &usuario)
		err = utilidades.FillStruct(v["CuerpoNotificacion"], &cuerpo)
		publish <- newEvent(models.EVENT_MESSAGE, strconv.FormatInt(usuario, 10), strings.Split(strconv.FormatInt(perfil, 10), ","), cuerpo)
		this.Ctx.Output.SetStatus(201)
		alert := models.Alert{Type: "success", Code: "S_544", Body: v}
		this.Data["json"] = alert
	} else {
		alert := models.Alert{Type: "error", Code: "E_N001", Body: err.Error()}
		this.Data["json"] = alert
	}
	this.ServeJSON()
}
