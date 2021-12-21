package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:ColasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:ColasController"],
        beego.ControllerComments{
            Method: "BorrarCola",
            Router: "/cola/:cola",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:ColasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:ColasController"],
        beego.ControllerComments{
            Method: "CrearCola",
            Router: "/crear/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:ColasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:ColasController"],
        beego.ControllerComments{
            Method: "BorrarMensaje",
            Router: "/mensaje/:cola",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:ColasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:ColasController"],
        beego.ControllerComments{
            Method: "BorrarMensajeId",
            Router: "/mensaje/:cola/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:ColasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:ColasController"],
        beego.ControllerComments{
            Method: "RecibirMensajes",
            Router: "/mensajes",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:NotificacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:NotificacionController"],
        beego.ControllerComments{
            Method: "PostOneNotif",
            Router: "/enviar/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:NotificacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:NotificacionController"],
        beego.ControllerComments{
            Method: "Subscribe",
            Router: "/suscribir/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:NotificacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:NotificacionController"],
        beego.ControllerComments{
            Method: "VerifSus",
            Router: "/suscripcion/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:NotificacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:NotificacionController"],
        beego.ControllerComments{
            Method: "BorrarTopic",
            Router: "/topic/",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:NotificacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:NotificacionController"],
        beego.ControllerComments{
            Method: "CreateTopic",
            Router: "/topics/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:NotificacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/notificacion_mid/controllers:NotificacionController"],
        beego.ControllerComments{
            Method: "GetTopics",
            Router: "/topics/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
