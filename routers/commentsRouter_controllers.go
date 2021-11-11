package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:ColasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:ColasController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:cola",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:ColasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:ColasController"],
        beego.ControllerComments{
            Method: "CrearCola",
            Router: "/crear/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:ColasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:ColasController"],
        beego.ControllerComments{
            Method: "RecibirMensajes",
            Router: "/mensajes",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:NotificacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:NotificacionController"],
        beego.ControllerComments{
            Method: "PostOneNotif",
            Router: "/enviar/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:NotificacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:NotificacionController"],
        beego.ControllerComments{
            Method: "Subscribe",
            Router: "/suscribir/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:NotificacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:NotificacionController"],
        beego.ControllerComments{
            Method: "GetTopics",
            Router: "/topics/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:NotificacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_api/controllers:NotificacionController"],
        beego.ControllerComments{
            Method: "CreateTopic",
            Router: "/topics/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
