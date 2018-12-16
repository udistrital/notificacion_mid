package routers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/notificacion_api/controllers"
)

func init() {
	// Register routers.
	//beego.Router("/", &controllers.AppController{})
	// Indicate AppController.Join method to handle POST requests.
	//beego.Router("/join", &controllers.AppController{}, "post:Join")

	// Long polling.
	/*beego.Router("/lp", &controllers.LongPollingController{}, "get:Join")
	beego.Router("/lp/post", &controllers.LongPollingController{})
	beego.Router("/lp/fetch", &controllers.LongPollingController{}, "get:Fetch")*/

	// WebSocket.
	//beego.Router("/ws", &controllers.WebSocketController{})
	beego.Router("/v1/ws/join", &controllers.WebSocketController{}, "get:Join")
	beego.Router("/v1/ws", &controllers.WebSocketController{}, "post:PushNotificacion")
	//send notification via api
	//beego.Router("/notify", &controllers.WebSocketController{}, "post:PushNotificacion")
	beego.Router("/v1/api/notify", &controllers.WebSocketController{}, "post:PushNotificacion")
	beego.Router("/v1/api/notify", &controllers.WebSocketController{}, "get:PushNotificacionDb")
}
