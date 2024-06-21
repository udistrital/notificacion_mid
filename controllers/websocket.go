package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/google/uuid"
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

// Función para enviar mensaje a un usuario
// Se tiene en cuenta si hay varias sesiones iniciadas (comparten el mismo prefijo = documento)
func sendMessageToClient(prefix string, messageType int, message []byte) {
	for key, conn := range usuarios {
		// Separar el documento y el UUID
		parts := strings.SplitN(key, "-", 2)
		if len(parts) > 0 && parts[0] == prefix {
			if err := conn.WriteMessage(messageType, message); err != nil {
				fmt.Println("Error al enviar mensaje a", key, ":", err)
			}
		}
	}
}

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

	// Generar un ID único para la conexión
	id := uuid.New().String()
	documento = documento + "-" + id

	// Registrar la conexión por el documento del usuario (identificador)
	usuarios[documento] = conn
	defer delete(usuarios, documento) // Limpiar la conexión al finalizar

	fmt.Printf("Usuario conectado: %s\n", documento)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("Error leyendo mensaje:", err)
			} else {
				fmt.Printf("Usuario desconectado: %s\n", documento)
			}
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

		// Publicar la notificacion a los usuarios destino de manera individual y enviar el mensaje al cliente
		if usuarios, ok := notificacion.Atributos["UsuariosDestino"].([]interface{}); ok {
			delete(notificacion.Atributos, "UsuariosDestino")
			auxIdDeduplicacion := notificacion.IdDeduplicacion
			for _, usuario := range usuarios {
				if idUsuario, ok := usuario.(string); ok {
					mensajeBody := notificacion
					mensajeBody.IdDeduplicacion = auxIdDeduplicacion + idUsuario
					mensajeBody.Atributos["UsuarioDestino"] = idUsuario
					mensajeBody.IdGrupoMensaje = idUsuario
					msg, _ := helpers.PublicarNotificacion(mensajeBody, true)

					modifiedMessage, err := json.Marshal(msg)
					if err != nil {
						fmt.Println("Error al codificar mensaje:", err)
						continue
					}

					sendMessageToClient(idUsuario+"wc", messageType, modifiedMessage)
				}
			}
		}
	}
}
