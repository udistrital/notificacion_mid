package models

import (
	"time"
)

type NotificacionConfiguracionPerfil struct {
	Id                        int                        `orm:"column(id);pk"`
	NotificacionConfiguracion *NotificacionConfiguracion `orm:"column(notificacion_configuracion);rel(fk)"`
	Perfil                    *Perfil                    `orm:"column(perfil);rel(fk)"`
}
type Perfil struct {
	Id         int         `orm:"column(id);pk;auto"`
	Nombre     string      `orm:"column(nombre)"`
	Aplicacion *Aplicacion `orm:"column(aplicacion);rel(fk)"`
}
type NotificacionConfiguracion struct {
	Id                              int                                `orm:"column(id);pk;auto"`
	EndPoint                        string                             `orm:"column(end_point)"`
	MetodoHttp                      *MetodoHttp                        `orm:"column(metodo_http);rel(fk)"`
	Tipo                            *NotificacionTipo                  `orm:"column(tipo);rel(fk)"`
	CuerpoNotificacion              string                             `orm:"column(cuerpo_notificacion);type(json)"`
	Aplicacion                      *Aplicacion                        `orm:"column(aplicacion);rel(fk)"`
	NotificacionConfiguracionPerfil []*NotificacionConfiguracionPerfil `orm:"reverse(many)"`
}
type MetodoHttp struct {
	Id          int    `orm:"column(id);pk"`
	Nombre      string `orm:"column(nombre)"`
	Descripcion string `orm:"column(descripcion)"`
}
type NotificacionTipo struct {
	Id     int    `orm:"column(id);pk"`
	Nombre string `orm:"column(nombre);null"`
}

type NotificacionEstado struct {
	Id     int    `orm:"column(id);pk"`
	Nombre string `orm:"column(nombre);null"`
}

type Aplicacion struct {
	Id          int    `orm:"column(id);pk;auto"`
	Nombre      string `orm:"column(nombre)"`
	Descripcion string `orm:"column(descripcion)"`
	Dominio     string `orm:"column(dominio)"`
	Estado      bool   `orm:"column(estado)"`
}

type Notificacion struct {
	Id                        int                        `orm:"column(id);pk;auto"`
	FechaCreacion             time.Time                  `orm:"column(fecha_creacion);type(timestamp with time zone);auto_now_add"`
	EstadoNotificacion        *NotificacionEstado        `orm:"column(estado_notificacion);rel(fk)"`
	CuerpoNotificacion        string                     `orm:"column(cuerpo_notificacion);type(json);null"`
	NotificacionConfiguracion *NotificacionConfiguracion `orm:"column(notificacion_configuracion);rel(fk)"`
}

type Alert struct {
	Type string
	Code string
	Body interface{}
}

/*type Subscriber struct {
	Id       string
	Profiles []string
	Conn     *websocket.Conn // Only for WebSocket users; otherwise nil.
}*/
