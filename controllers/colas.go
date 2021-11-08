package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/notificacion_api/helpers"
	"github.com/udistrital/notificacion_api/models"
)

// ColasController operations for Colas
type ColasController struct {
	beego.Controller
}

// URLMapping ...
func (c *ColasController) URLMapping() {
	c.Mapping("CrearCola", c.CrearCola)
	c.Mapping("RecibirMensajes", c.RecibirMensajes)
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
	logs.Debug(cola.Nombre)
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
// @Param	nombre	query 	string	true	"Nombre de la cola"
// @Success 201 {object} models.Cola
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

	url := "https://sqs.us-east-1.amazonaws.com/505909609706/" + c.GetString("nombre")

	if respuesta, err := helpers.RecibirMensajes(url); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err)
	}
	c.ServeJSON()
}
