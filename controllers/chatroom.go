package controllers

import (
	"container/list"
	"strings"
	"time"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/udistrital/notificacion_api/models"
	"github.com/udistrital/notificacion_api/utilidades"
)

type Subscription struct {
	Archive []models.Event      // All the events from the archive.
	New     <-chan models.Event // New events coming in.
}

func newEvent(ep models.EventType, user string, userDestination []string, profiles []string, msg map[string]interface{}, date time.Time, alias string, estiloicono string, estado string) models.Event {
	return models.Event{ep, user, profiles, int(time.Now().Unix()), msg, date, userDestination, alias, estiloicono, estado }
}

// Join ...
// @Title Join
// @Description Join to webSocket notification
// @Param	user		path 	string	true		"user from autentication"
// @Param	profiles		path 	string	true		"users profiles"
// @Success 200 {object} :user profiles conected
// @Failure 403 :profiles is empty
// @Failure 403 :user is empty
// @router /:user/:profiles [get]
func Join(user string, user_time string, profiles []string, ws *websocket.Conn) {
	var m []models.Notificacion
	utilidades.GetJson(beego.AppConfig.String("configuracionUrl")+"notificacion_estado_usuario/getOldNotification/"+strings.Join(profiles, ",")+"/"+user, &m)
	beego.Info(m)
	subscribe <- Subscriber{Name: user, UserTime: user_time,Profiles: profiles, Conn: ws}
}

func Leave(user string) {
	unsubscribe <- user
}

type Subscriber struct {
	Name     string
	UserTime  string
	Profiles []string
	Conn     *websocket.Conn // Only for WebSocket users; otherwise nil.
}

var (
	// Channel for new join users.
	subscribe = make(chan Subscriber, 10)
	// Channel for exit users.
	unsubscribe = make(chan string, 10)
	// Send events here to publish them.
	publish = make(chan models.Event, 10)

	subscribers = list.New()
	//definicion del pull de usuarios a notificar.
	connectionsId      = make(map[string]*websocket.Conn)
	connectionsProfile = make(map[string]map[string]*websocket.Conn)
)

// This function handles all incoming chan messages.
func chatroom() {
	for {
		select {
		case sub := <-subscribe:
			// var cuerpo = map[string]interface{}{"Message": "Conectado"}
			// var usuarioDestino []string
			if !isUserExist(subscribers, sub.Name) {
				subscribers.PushBack(sub) // Add user to the end of list.
				connectionsId[sub.Name+"_"+sub.UserTime] = sub.Conn
				// beego.Info(connectionsId)
				// beego.Info(users(connectionsId, sub.Name))
				for _, profile := range sub.Profiles {
					if _, ok := connectionsProfile[profile]; ok {
						(connectionsProfile[profile])[sub.Name+"_"+sub.UserTime] = sub.Conn
					} else {
						connectionsProfile[profile] = make(map[string]*websocket.Conn)
						(connectionsProfile[profile])[sub.Name+"_"+sub.UserTime] = sub.Conn
					}
					beego.Info("Register profile:", profile)
				}
				// Publish a JOIN event.
				// publish <- newEvent(models.EVENT_MESSAGE, sub.Name, usuarioDestino, sub.Profiles, cuerpo, time.Now().Local(), "", "", "conected") // Publish a LEAVE event. remove this message for prodct.
				beego.Info("New user:", sub.Name+"_"+sub.UserTime, ";WebSocket:", sub.Conn != nil)
			} else {
				// publish <- newEvent(models.EVENT_MESSAGE, sub.Name, usuarioDestino, sub.Profiles, cuerpo, time.Now().Local(),"", "", "conected") // Publish a LEAVE event. remove this message for prodct.
				beego.Info("Old user:", sub.Name+"_"+sub.UserTime, ";WebSocket:", sub.Conn != nil)
			}

		case event := <-publish:

			broadcastWebSocket(event)
			models.NewArchive(event)

			// if event.TypeEvent == models.EVENT_MESSAGE {
			// 	beego.Info("Message from", event.User, ";Content:", event.FechaCreacion)
			// }
		case unsub := <-unsubscribe:
			// beego.Info(unsub)
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).Name + "_" + sub.Value.(Subscriber).UserTime == unsub {
					subscribers.Remove(sub)
					delete(connectionsId, unsub)
					// Clone connection.
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
						beego.Error("WebSocket unsubscribe:", unsub)
					}
					//publish <- newEvent(models.EVENT_LEAVE, unsub, nil, "logout") // Publish a LEAVE event.
					break
				}
			}
		}
	}
}

func init() {
	go chatroom()
}

func isUserExist(subscribers *list.List, user string) bool {
	// beego.Info(user)
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).Name == user {
			return false
		}
	}
	return false
}
// 
func users(userMap map[string]*websocket.Conn, userName string) []string{
	keys := make([]string, 0, len(userMap))
    for k := range userMap {
		beego.Info(k)
		i := strings.Index(k, userName)
		if(i != -1) {
			keys = append(keys, k)
		}
	}
	return keys
}