package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
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
	c.Mapping("RecibirMensajesPorUsuario", c.RecibirMensajesPorUsuario)
	c.Mapping("BorrarMensaje", c.BorrarMensaje)
	c.Mapping("BorrarMensajeFiltro", c.BorrarMensajeFiltro)
	c.Mapping("BorrarCola", c.BorrarCola)
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

	defer helpers.ErrorController(c.Controller, "CrearCola")

	json.Unmarshal(c.Ctx.Input.RequestBody, &cola)
	if cola.NombreCola == "" {
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
// @Param	numMax			query 	int		false	"Numero máximo de mensajes que se pueden recibir (1-10) Por defecto, su valor es 1"
// @Success 201 {object} models.Mensaje
// @Failure 400 Error en parametros ingresados
// @router /mensajes [get]
func (c *ColasController) RecibirMensajes() {
	defer helpers.ErrorController(c.Controller, "RecibirMensajes")

	tiempoOcultoStr := c.GetString("tiempoOculto")
	if tiempoOcultoStr == "" {
		tiempoOcultoStr = "3"
	}

	tiempoOculto, err := strconv.Atoi(tiempoOcultoStr)
	if err != nil {
		panic(map[string]interface{}{"funcion": "RecibirMensaje", "err": "Error en parámetros de ingresos", "status": "400"})
	}

	numMaxStr := c.GetString("numMax")
	if numMaxStr == "" {
		numMaxStr = "1"
	}

	numMax, err := strconv.Atoi(numMaxStr)
	if err != nil {
		panic(map[string]interface{}{"funcion": "RecibirMensaje", "err": "Error en parámetros de ingresos", "status": "400"})
	}

	if respuesta, err := helpers.RecibirMensajes(c.GetString("nombre"), tiempoOculto, numMax); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// RecibirMensajesPorUsuario ...
// @Title RecibirMensajesPorUsuario
// @Description Lista todos los mensajes de una cola por el documento de un usuario
// @Param	nombre			query 	string	true	"Nombre de la cola"
// @Param	documento		query 	int		true	"Documento del usuario"
// @Param	numRevisados	query 	int		false	"Cantidad de mensajes revisados a recibir, seguidos de los mensajes pendientes. Por defecto se recibiran 5 si la cantidad de mensajes revisados es igual o mayor a este valor, de lo contrario se recibiran la cantidad de mensajes revisados disponibles. Para obtener todos, asignar el valor de -1"
// @Success 201 {object} models.Mensaje
// @Failure 400 Error en parametros ingresados
// @router /mensajes/usuario [get]
func (c *ColasController) RecibirMensajesPorUsuario() {
	defer helpers.ErrorController(c.Controller, "RecibirMensajesPorUsuario")

	nombreCola := c.GetString("nombre")
	documento := c.GetString("documento")

	numMaxRevisadosStr := c.GetString("numRevisados")
	if numMaxRevisadosStr == "" {
		numMaxRevisadosStr = "5"
	}

	numRevisados, err := strconv.Atoi(numMaxRevisadosStr)
	if err != nil {
		panic(map[string]interface{}{"funcion": "RecibirMensajesPorUsuario", "err": "Error en parámetros de ingresos", "status": "400"})
	}

	if respuesta, err := helpers.RecibirMensajesPorUsuario(nombreCola, documento, numRevisados); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// EsperarMensajes ...
// @Title EsperarMensajes
// @Description Espera por un tiempo determinado a que los mensajes estén disponibles y devuelve los recibidos en ese intervalo de tiempo
// @Param	nombre			query 	string	true	"Nombre de la cola"
// @Param	tiempoEspera	query 	int		true	"Tiempo de espera del api por mensajes"
// @Param	cantidad		query 	int		false	"Cantidad máxima de mensajes a recibir. Esta cantidad debe ser menor a diez veces el tiempo de espera, ya que se pueden obtener máximo 10 mensajes por segundo. Por defecto, se recibirán todos"
// @Param	filtro			query 	string	false	"Recepción de mensajes filtrados por metadata. Tiene el funcionamiento de un and, por lo tanto sólo devuelve los valores que cumplan con todo el filtro"
// @Success 201 {object} models.Mensaje
// @Failure 400 Error en parametros ingresados
// @router /mensajes/espera [get]
func (c *ColasController) EsperarMensajes() {
	filtro := make(map[string]string)

	defer helpers.ErrorController(c.Controller, "EsperarMensajes")

	tiempoEsperaStr := c.GetString("tiempoEspera")
	cantidadStr := c.GetString("cantidad")

	if v := c.GetString("filtro"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("error: atributos invalidos")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			filtro[k] = v
		}
	}

	tiempoEspera, err := strconv.Atoi(tiempoEsperaStr)
	if err != nil {
		panic(map[string]interface{}{"funcion": "EsperarMensajes", "err": "Error en parámetros de ingresos", "status": "400"})
	}

	cantidad, err := strconv.Atoi(cantidadStr)
	if (err != nil && cantidadStr != "") || cantidad > tiempoEspera*10 {
		panic(map[string]interface{}{"funcion": "EsperarMensajes", "err": "Error en parámetros de ingresos", "status": "400"})
	}

	if respuesta, err := helpers.EsperarMensajes(c.GetString("nombre"), tiempoEspera, cantidad, filtro); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// BorrarMensaje ...
// @Title BorrarMensaje
// @Description Borra la notificación de la cola
// @Param	cola		path 	string			true		"Nombre de la cola en donde está el mensaje"
// @Param	mensaje		body 	models.Mensaje	true		"Mensaje a borrar"
// @Success 200 {string} Mensaje eliminado
// @Failure 404 not found resource
// @router /mensajes/:cola [post]
func (c *ColasController) BorrarMensaje() {
	colaStr := c.Ctx.Input.Param(":cola")
	var mensaje models.Mensaje

	defer helpers.ErrorController(c.Controller, "BorrarMensaje")

	json.Unmarshal(c.Ctx.Input.RequestBody, &mensaje)
	if mensaje.ReceiptHandle == "" {
		panic(map[string]interface{}{"funcion": "BorrarMensaje", "err": "Error en parámetros de ingresos", "status": "400"})
	}

	if conteo, err := helpers.BorrarMensaje(colaStr, mensaje); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": map[string]interface{}{"MensajesEliminados": conteo}}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// BorrarMensajeFiltro ...
// @Title BorrarMensajeFiltro
// @Description Borra la notificación de la cola según el key y el valor ingresados en el body
// @Param	filtro		body 	models.Filtro		true		"Filtro de los mensajes a borrar"
// @Success 200 {string} Mensaje eliminado
// @Failure 404 not found resource
// @router /mensajes [post]
func (c *ColasController) BorrarMensajeFiltro() {
	var filtro models.Filtro

	defer helpers.ErrorController(c.Controller, "BorrarMensajeFiltro")

	json.Unmarshal(c.Ctx.Input.RequestBody, &filtro)
	if filtro.NombreCola == "" {
		panic(map[string]interface{}{"funcion": "BorrarMensajeFiltro", "err": "Error en parámetros de ingresos", "status": "400"})
	}

	if conteo, err := helpers.BorrarMensajeFiltro(filtro); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": map[string]interface{}{"MensajesEliminados": conteo}}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// BorrarCola ...
// @Title BorrarCola
// @Description Borra la cola
// @Param	cola		path 	string			true		"Nombre de la cola a borrar"
// @Success 200 {string} Cola eliminada
// @Failure 502 Error en borrado de cola
// @router /cola/:cola [delete]
func (c *ColasController) BorrarCola() {
	colaStr := c.Ctx.Input.Param(":cola")

	defer helpers.ErrorController(c.Controller, "BorrarCola")

	if err := helpers.BorrarCola(colaStr); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": "Cola eliminada"}
	} else {
		panic(err)
	}
	c.ServeJSON()
}
