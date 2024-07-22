package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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
func sendNotificationToClient(prefijo string, messageType int, message []byte) {
	for key, conn := range usuarios {
		// Separar el documento y el UUID
		parts := strings.SplitN(key, "-", 2)
		if len(parts) > 0 && parts[0] == prefijo {
			if err := conn.WriteMessage(messageType, message); err != nil {
				log.Printf("Error al enviar notificación a %s: %v", key, err)
			} else {
				log.Printf("\t- Enviada a: %s", prefijo)
			}
		}
	}
}

// Función para verificar si un usuario ya existe (está conectado)
func verifyClient(prefix string) bool {
	for key := range usuarios {
		if strings.HasPrefix(key, prefix+"-") {
			return true
		}
	}
	return false
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
		log.Printf("Error al actualizar la conexión: %v", err)
		return
	}
	defer conn.Close()

	// Leer documento del usuario desde el cliente
	var documento string
	err = conn.ReadJSON(&documento)
	if err != nil {
		log.Printf("Error al leer el documento del usuario: %v", err)
		return
	}

	if !verifyClient(documento) {
		// Generar un ID único para la conexión
		id := uuid.New().String()
		documento = documento + "-" + id

		// Registrar la conexión por el documento del usuario (identificador)
		usuarios[documento] = conn
		defer delete(usuarios, documento) // Limpiar la conexión al finalizar

		log.Printf("Usuario conectado: %s", strings.Split(documento, "-")[0])
	}

	for {
		// Obtener notificación o comprobar desconexión
		messageType, messageData, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error leyendo notificación: %v", err)
			} else {
				log.Printf("Usuario desconectado: %s", strings.Split(documento, "-")[0])
			}
			return
		}

		var notificaciones []map[string]interface{}
		err = json.Unmarshal(messageData, &notificaciones)
		if err != nil {
			log.Printf("Error al decodificar la notificación: %v", err)
			continue
		}

		log.Printf("Notificación recibida de: %s", strings.Split(documento, "-")[0])

		// Enviar notificación a los usuarios destino en tiempo real
		for i := 0; i < len(notificaciones); i++ {
			destinatario, ok := notificaciones[i]["destinatario"].(string)
			if ok {
				resNotificacion, err := json.Marshal(notificaciones[i])
				if err != nil {
					log.Printf("Error al codificar notificación: %v", err)
					continue
				}
				sendNotificationToClient(destinatario+"wc", messageType, resNotificacion)
			}
		}
	}
}
