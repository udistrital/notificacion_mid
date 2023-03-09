package helpers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func ContainsJson(slice []string, element map[string]interface{}) (contains bool) {
	for _, v := range slice {
		if v == element["Value"] {
			return true
		}
	}
	return
}

func ContainsString(slice []string, element string) (contains bool) {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return
}

// Manejo único de errores para controladores sin repetir código
func ErrorController(c beego.Controller, controller string) {
	if err := recover(); err != nil {
		logs.Error(err)
		localError := err.(map[string]interface{})
		c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + controller + "/" + (localError["funcion"]).(string))
		c.Data["data"] = (localError["err"])
		if status, ok := localError["status"]; ok {
			c.Abort(status.(string))
		} else {
			c.Abort("500")
		}
	}
}
