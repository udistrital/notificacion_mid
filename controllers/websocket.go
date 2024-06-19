package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/udistrital/notificacion_mid/helpers"
	"github.com/udistrital/notificacion_mid/models"
)

type WebSocketController struct {
	beego.Controller
}

// URLMapping ...
func (c *WebSocketController) URLMapping() {
	c.Mapping("WebSocket", c.WebSocket)
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Permitir todas las conexiones
		},
	}
	usuarios = make(map[string]*websocket.Conn) // Mapa para almacenar conexiones activas (usuarios)
)

// Función para enviar mensaje a un usuario específico
func sendMessageToClient(documento string, messageType int, message []byte) error {
	conn, ok := usuarios[documento]
	if !ok {
		return fmt.Errorf("el usuario %s no está conectado", documento)
	}
	return conn.WriteMessage(messageType, message)
}

// Función para broadcast a todos los usuarios
// func broadcastMessage(messageType int, message []byte) {
// 	for _, conn := range usuarios {
// 		if err := conn.WriteMessage(messageType, message); err != nil {
// 			fmt.Println("Error al enviar mensaje a cliente:", err)
// 		}
// 	}
// }

// WebSocket ...
// @Title WebSocket
// @Description Recibir mensaje por medio de webSocket
// @Success 200 {string} Mensaje recibido
// @Failure 502 Error de conexión
// @router / [get]
func (c *WebSocketController) WebSocket() {
	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		fmt.Println("Error al actualizar la conexión:", err)
		return
	}
	defer conn.Close()

	// Leer documento del usuario desde el cliente
	var documento string
	err = conn.ReadJSON(&documento)
	if err != nil {
		fmt.Println("Error al leer el documento del usuario:", err)
		return
	}

	// Registrar la conexión por el documento del usuario (identificador)
	usuarios[documento] = conn
	defer delete(usuarios, documento) // Limpiar la conexión al finalizar

	fmt.Printf("Usuario conectado: %s\n", documento)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error leyendo mensage:", err)
			return
		}
		fmt.Printf("Mensaje recibido de %s: %s\n", documento, message)

		// Convertir el mensaje en el modelo Notificacion
		var notificacion models.Notificacion
		err = json.Unmarshal(message, &notificacion)
		if err != nil {
			fmt.Println("Error al decodificar mensaje en modelo Notificacion:", err)
			continue
		}

		// Publicar la notificacion a los usuarios destino  de manera individuak y enviar el mensaje al cliente
		if usuarios, ok := notificacion.Atributos["UsuariosDestino"].([]interface{}); ok {
			delete(notificacion.Atributos, "UsuariosDestino")
			auxIdDeduplicacion := notificacion.IdDeduplicacion
			for _, usuario := range usuarios {
				if idUsuario, ok := usuario.(string); ok {
					mensajeBody := notificacion
					mensajeBody.IdDeduplicacion = auxIdDeduplicacion + idUsuario
					// mensajeBody.Atributos["UsuarioDestino"] = idUsuario
					mensajeBody.Atributos["UsuarioDestino"] = "7230282"
					mensajeBody.IdGrupoMensaje = idUsuario
					msg, _ := helpers.PublicarNotificacion(mensajeBody, true)

					modifiedMessage, err := json.Marshal(msg)
					if err != nil {
						fmt.Println("Error al codificar mensaje:", err)
						continue
					}

					// sendMessageToClient(idUsuario+"ws", messageType, modifiedMessage)
					sendMessageToClient("7230282ws", messageType, modifiedMessage)
				}
			}
		}
	}
}
