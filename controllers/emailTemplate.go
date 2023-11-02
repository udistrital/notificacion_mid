package controllers

import (
	"encoding/json"
	"strconv"

	//"errors"
	//"strings"

	"github.com/astaxie/beego"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/udistrital/notificacion_mid/helpers"
	//"github.com/udistrital/notificacion_mid/models"
)

// EnviarEmailController operations for inputicacion
type EmailTemplateController struct {
	beego.Controller
}

// URLMapping ...
func (c *EmailTemplateController) URLMapping() {
	c.Mapping("CreateEmailTemplate", c.CreateEmailTemplate)
	c.Mapping("GetEmailTemplate", c.GetEmailTemplate)
	c.Mapping("ListEmailTemplate", c.ListEmailTemplate)
	c.Mapping("UpdateEmailTemplate", c.UpdateEmailTemplate)
	c.Mapping("DeleteEmailTemplate", c.DeleteEmailTemplate)
}

// CreateEmailTemplate ...
// @Title CreateEmailTemplate
// @Description Crea un template para un correo personalizado
// @Param	body		body 	ses.CreateTemplateInput	true		"Body ses.CreateTemplateInput struct"
// @Success 201 {object} map[string]interface{Success string,Status boolean,Message string,Data map[string]interface{}}
// @Failure 400 Error en parametros ingresados
// @router / [post]
func (c *EmailTemplateController) CreateEmailTemplate() {
	var input ses.CreateTemplateInput

	defer helpers.ErrorController(c.Controller, "EmailTemplateController")

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &input); err != nil {
		panic(map[string]interface{}{"funcion": "CreateEmailTemplate", "err": err, "status": "400"})
	}

	if respuesta, err := helpers.CreateEmailTemplate(input); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": map[string]interface{}{"Result": respuesta}}
	} else {
		panic(map[string]interface{}{"funcion": "CreateEmailTemplate/", "err": err, "status": "502"})
	}
	c.ServeJSON()
}

// GetEmailTemplate ...
// @Title GetEmailTemplate
// @Description traer todos los template de aws ses
// @Param	templateName	path 	string		true	"Nombre del template"
// @Success 201 {object} map[string]interface{Success string,Status boolean,Message string,Data map[string]interface{}}
// @Failure 400 Error en parametros ingresados
// @router /:templateName [get]
func (c *EmailTemplateController) GetEmailTemplate() {
	defer helpers.ErrorController(c.Controller, "EmailTemplateController")

	var input ses.GetTemplateInput

	if templateName := c.Ctx.Input.Param(":templateName"); templateName != "" {
		input.TemplateName = &templateName
	} else {
		panic(map[string]interface{}{"funcion": "GetEmailTemplate/", "err": "nombre del template vacio", "status": "400"})
	}

	if respuesta, err := helpers.GetEmailTemplate(input); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": map[string]interface{}{"Result": respuesta}}
	} else {
		panic(map[string]interface{}{"funcion": "GetEmailTemplate/", "err": err, "status": "502"})
	}
	c.ServeJSON()
}

// UpdateEmailTemplate ...
// @Title UpdateEmailTemplate
// @Description traer todos los template de aws ses
// @Param	body		body 	ses.UpdateTemplateInput	true		"Body ses.UpdateTemplateInput struct"
// @Success 200 {object} map[string]interface{Success string,Status boolean,Message string,Data map[string]interface{}}
// @Failure 400 Error en parametros ingresados
// @router / [put]
func (c *EmailTemplateController) UpdateEmailTemplate() {
	var input ses.UpdateTemplateInput

	defer helpers.ErrorController(c.Controller, "EmailTemplateController")

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &input); err != nil {
		panic(map[string]interface{}{"funcion": "UpdateEmailTemplate", "err": err, "status": "400"})
	}

	if _, err := helpers.UpdateEmailTemplate(input); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": map[string]interface{}{"Result": "Template actualizado con exito"}}
	} else {
		panic(map[string]interface{}{"funcion": "UpdateEmailTemplate", "err": err, "status": "502"})
	}
	c.ServeJSON()
}

// ListEmailTemplate ...
// @Title ListEmailTemplate
// @Description traer todos los template de aws ses
// @Param	nextToken	query 	string		false	"posicion en la lista de una llamada previa"
// @Param	maxItems	query 	int		false	"Numero máximo de mensajes que se pueden recibir (1-100) Por defecto, su valor es 10"
// @Success 200 {object} map[string]interface{Success string,Status boolean,Message string,Data map[string]interface{}}
// @Failure 400 Error en parametros ingresados
// @router / [get]
func (c *EmailTemplateController) ListEmailTemplate() {
	defer helpers.ErrorController(c.Controller, "EmailTemplateController")

	var input ses.ListTemplatesInput

	maxItemsStr := c.GetString("maxItems")

	if maxItemsStr != "" {
		if maxItems, err := strconv.Atoi(maxItemsStr); err != nil {
			maxItemsInt64 := int64(maxItems)
			input.MaxItems = &maxItemsInt64
		} else {
			panic(map[string]interface{}{"funcion": "PostSendEmail/ErrorMaxItems", "err": err, "status": "400"})
		}
	}

	nextToken := c.GetString("nextToken")
	if nextToken != "" {
		input.NextToken = &nextToken
	}

	if respuesta, err := helpers.ListEmailTemplates(input); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": map[string]interface{}{"Result": respuesta}}
	} else {
		panic(map[string]interface{}{"funcion": "ListEmailTemplate/", "err": err, "status": "502"})
	}
	c.ServeJSON()
}

// DeleteEmailTemplate ...
// @Title DeleteEmailTemplate
// @Description Borra el topic
// @Param    templateName    path    string    true    "Nombre del template a elimminar"
// @Success 200 {object} map[string]interface{Success string,Status boolean,Message string,Data map[string]interface{}}
// @Failure 404 not found resource
// @router /:templateName [delete]
func (c *EmailTemplateController) DeleteEmailTemplate() {
	defer helpers.ErrorController(c.Controller, "EmailTemplateController")

	var input ses.DeleteTemplateInput

	if templateName := c.Ctx.Input.Param(":templateName"); templateName != "" {
		input.TemplateName = &templateName
	} else {
		panic(map[string]interface{}{"funcion": "DeleteEmailTemplate", "err": "nombre del template vacio", "status": "400"})
	}

	if _, err := helpers.DeleteEmailTemplate(input); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": map[string]interface{}{"Result": "template eliminado con exito"}}
	} else {
		panic(map[string]interface{}{"funcion": "DeleteEmailTemplate/", "err": err, "status": "502"})
	}
	c.ServeJSON()
}
