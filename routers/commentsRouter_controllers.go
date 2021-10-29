package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:NotificacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:NotificacionController"],
        beego.ControllerComments{
            Method: "PostOneNotif",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
