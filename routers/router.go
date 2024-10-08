// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/udistrital/notificacion_mid/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/notificaciones",
			beego.NSInclude(
				&controllers.NotificacionController{},
			),
		),
		beego.NSNamespace("/colas",
			beego.NSInclude(
				&controllers.ColasController{},
			),
		),
		beego.NSNamespace("/email",
			beego.NSInclude(
				&controllers.EnviarEmailController{},
			),
		),
		beego.NSNamespace("/template_email",
			beego.NSInclude(
				&controllers.EmailTemplateController{},
			),
		),
		beego.NSNamespace("/ws",
			beego.NSInclude(
				&controllers.WebSocketController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
