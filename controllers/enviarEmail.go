package controllers

import (
	"encoding/json"

	//"errors"
	//"strings"

	"github.com/astaxie/beego"

	"github.com/udistrital/notificacion_mid/helpers"
	"github.com/udistrital/notificacion_mid/models"
)

// EnviarEmailController operations for inputicacion
type EnviarEmailController struct {
	beego.Controller
}

// URLMapping ...
func (c *EnviarEmailController) URLMapping() {
	c.Mapping("PostSendEmail", c.PostSendEmail)
}

// PostSendEmail ...
// @Title PostSendEmail
// @Description Envia un correo a los destinatarios elegidos
// @Param	body		body 	models.inputicacion	true		"Body de la inputicación"
// @Success 201 {object} map[string]interface{Success string,Status boolean,Message string,Data map[string]interface{}}
// @Failure 400 Error en parametros ingresados
// @router /enviar_email/ [post]
func (c *EnviarEmailController) PostSendEmail() {
	var input models.SendEmailInput

	defer helpers.ErrorController(c.Controller, "PostSendEmail")

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &input); err != nil {
		panic(map[string]interface{}{"funcion": "PostSendEmail", "err": err, "status": "400"})
	}

	if respuesta, err := helpers.SendEmail(input); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": map[string]interface{}{"Result": respuesta}}
	} else {
		panic(err)
	}
	c.ServeJSON()
}
