package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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
	defer Leave(Id)
	Join(Id, Profiles, ws)

	// Message receive loop.
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		var m map[string]interface{}
		err = json.Unmarshal(p, &m)
		//publish <- newEvent(models.EVENT_MESSAGE, Id, nil, Profiles, m, time.Now().Local(),)
	}
}

// broadcastWebSocket broadcasts messages to WebSocket users.
func broadcastWebSocket(event models.Event) {
	data, err := json.Marshal(event)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}
	for _, user := range event.UserDestination {
		if connectionsId[user] != nil {
			ws := connectionsId[user]
			if ws != nil {
				if ws.WriteMessage(websocket.TextMessage, data) != nil {
					// User disconnected.
					unsubscribe <- event.User
				}
			}
		}
	}

	for _, value := range event.Profiles {
		fmt.Println("message from ", event.User)
		if connectionsProfile[value] != nil {
			for user, con := range connectionsProfile[value] {
				ws := con
				var m []models.Notificacion
				utilidades.GetJson(beego.AppConfig.String("configuracionUrl")+"notificacion_estado_usuario/getOldNotification/"+value+"/"+user, &m)
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
	var res map[string]interface{}
	//UserId := c.GetString("id")
	//fmt.Println("Id ", UserId)
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &v); err == nil {
		//push notificacion-------
		beego.Info("Data ", v)
		var perfil []string
		var usuario string
		var usuarioDestino []string
		var cuerpo map[string]interface{}
		var alias string
		var estiloIcono string
		var estado string

		err = utilidades.FillStruct(v["DestinationProfiles"], &perfil)
		err = utilidades.FillStruct(v["Application"], &usuario)
		err = utilidades.FillStruct(v["NotificationBody"], &cuerpo)
		err = utilidades.FillStruct(v["UserDestination"], &usuarioDestino)
		err = utilidades.FillStruct(v["Alias"], &alias)
		err = utilidades.FillStruct(v["EstiloIcono"], &estiloIcono)
		err = utilidades.FillStruct(v["Estado"], &estado)

		publish <- newEvent(models.EVENT_MESSAGE, usuario, usuarioDestino, perfil, cuerpo, time.Now().Local(), alias, estiloIcono, estado)
		j, error := json.Marshal(cuerpo)
		if error == nil {
			if v["UserDestination"] == "" {
				data := map[string]interface{}{
					"CuerpoNotificacion":        string(j),
					"NotificacionConfiguracion": map[string]interface{}{"Id": v["ConfiguracionNotificacion"]}}
				utilidades.SendJson(beego.AppConfig.String("configuracionUrl")+"notificacion/", "POST", &res, data)
				beego.Info("respuesta servicio", res)
				beego.Info(beego.AppConfig.String("configuracionUrl") + "notificacion")
				this.Ctx.Output.SetStatus(201)
				alert := models.Alert{Type: "success", Code: "S_544", Body: v}
				this.Data["json"] = alert
			} else {

				data := map[string]interface{}{
					"CuerpoNotificacion":        string(j),
					"NotificacionConfiguracion": map[string]interface{}{"Id": v["ConfiguracionNotificacion"]}}

				notificacion := map[string]interface{}{
					"Notificacion": data,
					"Usuarios":     v["UserDestination"]}

				utilidades.SendJson(beego.AppConfig.String("configuracionUrl")+"notificacion_estado_usuario/pushNotificationUser", "POST", &res, notificacion)
				beego.Info(beego.AppConfig.String("configuracionUrl") + "notificacion_estado_usuario/pushNotificationUser")
				this.Ctx.Output.SetStatus(201)
				alert := models.Alert{Type: "success", Code: "S_544", Body: notificacion}
				this.Data["json"] = alert
			}
		} else {
			beego.Info(error)
		}

	} else {
		alert := models.Alert{Type: "error", Code: "E_N001", Body: err.Error()}
		this.Data["json"] = alert
	}
}

// broadcastWebSocket broadcasts messages to WebSocket users from db
func (this *WebSocketController) PushNotificacionDb() {
	fmt.Println("entro")
	var m []models.Notificacion
	query := this.GetString("query")
	//fmt.Println("Id ", UserId)
	if err := utilidades.GetJson(beego.AppConfig.String("configuracionUrl")+"/notificacion/?query="+query, &m); err == nil {
		for _, v := range m {
			//push notificacion-------
			beego.Info("Data ", v)
			var perfil []string
			var usuario string
			var cuerpo map[string]interface{}
			var alias string
			var estiloicono string
			var estado string

			for _, profiledata := range v.NotificacionConfiguracion.NotificacionConfiguracionPerfil {
				perfil = append(perfil, profiledata.Perfil.Nombre)
			}
			usuario = v.NotificacionConfiguracion.Aplicacion.Nombre
			err = json.Unmarshal([]byte(v.CuerpoNotificacion), &cuerpo)
			publish <- newEvent(models.EVENT_MESSAGE, usuario, nil, perfil, cuerpo, v.FechaCreacion, alias, estiloicono, estado)
		}
		this.Ctx.Output.SetStatus(201)
		alert := models.Alert{Type: "success", Code: "S_544", Body: m}
		this.Data["json"] = alert
	} else {
		alert := models.Alert{Type: "error", Code: "E_N001", Body: err.Error()}
		this.Data["json"] = alert
	}
	this.ServeJSON()
}
