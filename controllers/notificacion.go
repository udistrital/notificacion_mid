package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/notificacion_api/helpers"
	"github.com/udistrital/notificacion_api/models"
)

// NotificacionController operations for Notificacion
type NotificacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *NotificacionController) URLMapping() {
	c.Mapping("PostOneNotif", c.PostOneNotif)
}

// PostOneNotif ...
// @Title PostOneNotif
// @Description Envía una notificación a cualquier suscriptor. La propiedad "Atributos" es opcional
// @Param	body		body 	models.Notificacion	true		"Body de la notificación"
// @Success 201 {object} map[string]interface{"Success","Status","Message","Data"}
// @Failure 400 Error en parametros ingresados
// @router / [post]
func (c *NotificacionController) PostOneNotif() {
	var notif models.Notificacion

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/PostOneNotif/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	json.Unmarshal(c.Ctx.Input.RequestBody, &notif)
	if notif.RemitenteId == "" || len(notif.DestinatarioId) == 0 || notif.Asunto == "" || notif.Mensaje == "" {
		panic(map[string]interface{}{"funcion": "PostOneNotif", "err": "Error en parámetros de ingresos", "status": "400"})
	}
	if respuesta, err := helpers.PublicarNotificacion(notif); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err)
	}
	c.ServeJSON()
}
