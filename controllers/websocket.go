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

// Función para enviar la notificación a un usuario
// Se tiene en cuenta si hay varias sesiones iniciadas (comparten el mismo prefijo = documento)
func sendNotificationToClient(prefix string, messageType int, message []byte) {
	for key, conn := range usuarios {
		// Separar el documento y el UUID
		parts := strings.SplitN(key, "-", 2)
		if len(parts) > 0 && parts[0] == prefix {
			if err := conn.WriteMessage(messageType, message); err != nil {
				fmt.Println("Error al enviar notificación a", key, ":", err)
			}
		}
	}
}

// WebSocket ...
// @Title WebSocket
// @Description Recibir notificación por medio de webSocket
// @Success 200 {string} Notificación recibida
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
		// Obtener notificación o comprobar desconexión
		messageType, messageData, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("Error leyendo notificación:", err)
			} else {
				fmt.Printf("Usuario desconectado: %s\n", documento)
			}
			return
		}

		var notificacion map[string]interface{}
		err = json.Unmarshal(messageData, &notificacion)
		if err != nil {
			fmt.Println("Error al decodificar la notificación:", err)
			continue
		}

		fmt.Printf("Notificación recibida de %s: %s\n", documento, notificacion)

		// Registrar la notificación a los usuarios destino de manera individual y enviarla al cliente
		if usuarios, ok := notificacion["destinatarios"].([]interface{}); ok {
			delete(notificacion, "destinatarios")
			for _, usuario := range usuarios {
				if _, ok := usuario.(string); ok {
					// notificacion["destinatario"] = idUsuario
					notificacion["destinatario"] = "7230282"

					res, err := helpers.PublicarNotificacionCrud(notificacion)
					if err == nil {
						resNotificacion, err := json.Marshal(res["Data"])
						if err != nil {
							fmt.Println("Error al codificar notificación:", err)
							continue
						}
						// sendNotificationToClient(idUsuario+"wc", messageType, resNotificacion)
						sendNotificationToClient("7230282wc", messageType, resNotificacion)
					}
				}
			}
		}
	}
}
