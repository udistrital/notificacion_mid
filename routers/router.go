package routers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/notificacion_api/controllers"
)

func init() {
	// Register routers.

	// WebSocket.
	beego.Router("/ws", &controllers.WebsocketController{})
	beego.Router("/ws/join", &controllers.WebsocketController{}, "get:Join")

}
