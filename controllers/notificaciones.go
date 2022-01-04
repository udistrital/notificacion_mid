package controllers

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/notificacion_mid/helpers"
	"github.com/udistrital/notificacion_mid/models"
)

// NotificacionController operations for Notificacion
type NotificacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *NotificacionController) URLMapping() {
	c.Mapping("PostOneNotif", c.PostOneNotif)
	c.Mapping("Subscribe", c.Subscribe)
	c.Mapping("GetTopics", c.GetTopics)
	c.Mapping("CreateTopic", c.CreateTopic)
	c.Mapping("VerifSus", c.VerifSus)
}

// PostOneNotif ...
// @Title PostOneNotif
// @Description Envía una notificación a cualquier suscriptor. La propiedad "Atributos" es opcional
// @Param	body		body 	models.Notificacion	true		"Body de la notificación"
// @Success 201 {object} map[string]interface{Success string,Status boolean,Message string,Data map[string]interface{}}
// @Failure 400 Error en parametros ingresados
// @router /enviar/ [post]
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
	if notif.RemitenteId == "" || len(notif.DestinatarioId) == 0 || notif.Asunto == "" || notif.Mensaje == "" || notif.ArnTopic == "" {
		panic(map[string]interface{}{"funcion": "PostOneNotif", "err": "Error en parámetros de ingresos", "status": "400"})
	}

	if respuesta, err := helpers.PublicarNotificacion(notif); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": map[string]interface{}{"MessageId": respuesta}}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// Subscribe ...
// @Title Subscribe
// @Description Suscribe cualquier tipo de endpoint a un topic
// @Param	body		body 	models.Suscripcion	true		"Body de la suscripcion"
// @Param	atributos	query	string				false		"Atributos para filtrado de mensajes (atributo:valor, ...)"
// @Success 201 {object} map[string]interface{Success string,Status boolean,Message string,Data map[string]interface{}}
// @Failure 400 Error en parametros ingresados
// @router /suscribir/ [post]
func (c *NotificacionController) Subscribe() {
	var sub models.Suscripcion
	var atributos = make(map[string]string)

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/Subscribe/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	json.Unmarshal(c.Ctx.Input.RequestBody, &sub)
	for _, subscriptor := range sub.Suscritos {
		prot := subscriptor.Protocolo
		if prot != "kinesis" && prot != "lambda" && prot != "sqs" && prot != "email" && prot != "email-json" && prot != "http" && prot != "https" && prot != "application" && prot != "sms" && prot != "firehouse" {
			panic(map[string]interface{}{"funcion": "PostOneNotif", "err": "Protocolo invalido", "status": "400"})
		}
	}

	if v := c.GetString("atributos"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("error: atributos invalidos")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			atributos[k] = v
		}
	}

	if respuesta, err := helpers.Suscribir(sub, atributos); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": map[string]interface{}{"TopicARN": respuesta}}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// CreateTopic ...
// @Title CreateTopic
// @Description Crea un topic en sns
// @Param	body		body 	models.Topic	true		"Body para configuracion del topic"
// @Success 201 {object} map[string]interface{Success string,Status boolean,Message string,Data []string }
// @Failure 400 Error en parametros ingresados
// @router /topic/ [post]
func (c *NotificacionController) CreateTopic() {
	var topic models.Topic

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/CreateTopic/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()
	json.Unmarshal(c.Ctx.Input.RequestBody, &topic)
	if topic.Nombre == "" || topic.Display == "" {
		panic(map[string]interface{}{"funcion": "PostOneNotif", "err": "Error en parámetros de ingresos", "status": "400"})
	}

	if respuesta, err := helpers.CrearTopic(topic); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// GetTopics ...
// @Title GetTopics
// @Description Lista todos los ARN de los topics disponibles
// @Success 201 {object} map[string]interface{Success string,Status boolean,Message string,Data []string}
// @Failure 400 Error en parametros ingresados
// @router /topic/ [get]
func (c *NotificacionController) GetTopics() {
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

	if respuesta, err := helpers.ListaTopics(); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// VerifSus ...
// @Title VerifSus
// @Description Verifica la suscripcion
// @Param	suscripcion		body 	models.ConsultaSuscripcion	true		"Suscripcion a consultar"
// @Success 200 {string} Mensaje eliminado
// @Failure 404 not found resource
// @router /suscripcion/ [post]
func (c *NotificacionController) VerifSus() {
	var consulta models.ConsultaSuscripcion

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/VerifSus/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	json.Unmarshal(c.Ctx.Input.RequestBody, &consulta)
	if consulta.ArnTopic == "" || consulta.Endpoint == "" {
		panic(map[string]interface{}{"funcion": "VerifSus", "err": "Error en parámetros de ingresos", "status": "400"})
	}
	if respuesta, err := helpers.VerificarSuscripcion(consulta.ArnTopic, consulta.Endpoint); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// BorrarTopic ...
// @Title BorrarTopic
// @Description Borra el topic
// @Param	arnTopic		query 	string			true		"Arn del topic a eliminar"
// @Success 200 {string} Topic eliminado
// @Failure 404 not found resource
// @router /topic/ [delete]
func (c *NotificacionController) BorrarTopic() {
	arn := c.GetString("arnTopic")

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/BorrarTopic/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	if err := helpers.BorrarTopic(arn); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": "Topic eliminado"}
	} else {
		panic(err)
	}
	c.ServeJSON()
}
