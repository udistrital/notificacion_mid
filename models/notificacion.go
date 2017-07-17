package models

import (
	"time"

	"github.com/gorilla/websocket"
)

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
	Id                 int                 `orm:"column(id);pk"`
	UsuarioDestino     int64               `orm:"column(usuario_destino);null"`
	PerfilDestino      int64               `orm:"column(perfil_destino);null"`
	AplicacionOrigen   *Aplicacion         `orm:"column(aplicacion_origen);rel(fk)"`
	FechaCreacion      time.Time           `orm:"column(fecha_creacion);type(timestamp without time zone)"`
	EstadoNotificacion *NotificacionEstado `orm:"column(estado_notificacion);rel(fk)"`
	CuerpoNotificacion string              `orm:"column(cuerpo_notificacion);type(json);null"`
	TipoNotificacion   *NotificacionTipo   `orm:"column(tipo_notificacion);rel(fk)"`
}

type Alert struct {
	Type string
	Code string
	Body interface{}
}

type Subscriber struct {
	Id       string
	Profiles []string
	Conn     *websocket.Conn // Only for WebSocket users; otherwise nil.
}
