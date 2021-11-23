package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/notificacion_mid/helpers"
	"github.com/udistrital/notificacion_mid/models"
)

// ColasController operations for Colas
type ColasController struct {
	beego.Controller
}

// URLMapping ...
func (c *ColasController) URLMapping() {
	c.Mapping("CrearCola", c.CrearCola)
	c.Mapping("RecibirMensajes", c.RecibirMensajes)
	c.Mapping("DeleteMessages", c.Delete)
}

// CrearCola ...
// @Title CrearCola
// @Description Crea colas sqs
// @Param	body		body 	models.Cola		true		"Configuración para la creación de la cola"
// @Success 201 {object} models.Colas
// @Failure 403 body is empty
// @router /crear/ [post]
func (c *ColasController) CrearCola() {
	var cola models.Cola

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/CrearCola/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	json.Unmarshal(c.Ctx.Input.RequestBody, &cola)
	if cola.Nombre == "" {
		panic(map[string]interface{}{"funcion": "CrearCola", "err": "Error en parámetros de ingresos", "status": "400"})
	}

	if respuesta, err := helpers.CrearCola(cola); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": map[string]interface{}{"QueueARN": respuesta}}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// RecibirMensajes ...
// @Title RecibirMensajes
// @Description Lista hasta 10 mensajes en cola
// @Param	nombre			query 	string	true	"Nombre de la cola"
// @Param	tiempoOculto	query 	int		false	"El tiempo en segundos que un mensaje recibido se ocultará en la cola"
// @Success 201 {object} models.Mensaje
// @Failure 400 Error en parametros ingresados
// @router /mensajes [get]
func (c *ColasController) RecibirMensajes() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/GetTopics/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	tiempoOcultoStr := c.GetString("tiempoOculto")
	if tiempoOcultoStr == "" {
		tiempoOcultoStr = "3"
	}

	tiempoOculto, err := strconv.Atoi(tiempoOcultoStr)
	if err != nil {
		panic(map[string]interface{}{"funcion": "RecibirMensaje", "err": "Error en parámetros de ingresos", "status": "400"})
	}

	if respuesta, err := helpers.RecibirMensajes(c.GetString("nombre"), tiempoOculto); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// DeleteMessage ...
// @Title DeleteMessage
// @Description Borra la notificación de la cola
// @Param	cola		path 	string			true		"Nombre de la cola en donde está el mensaje"
// @Param	mensaje		body 	models.Mensaje	true		"Mensaje a borrar"
// @Success 200 {string} Mensaje eliminado
// @Failure 404 not found resource
// @router /:cola [delete]
func (c *ColasController) Delete() {
	colaStr := c.Ctx.Input.Param(":cola")
	var mensaje models.Mensaje

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/BorrarMensaje/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	json.Unmarshal(c.Ctx.Input.RequestBody, &mensaje)
	if mensaje.ReceiptHandle == "" {
		panic(map[string]interface{}{"funcion": "BorrarMensaje", "err": "Error en parámetros de ingresos", "status": "400"})
	}

	if err := helpers.BorrarMensaje(colaStr, mensaje); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": "Mensaje eliminado"}
	} else {
		panic(err)
	}
	c.ServeJSON()
}
