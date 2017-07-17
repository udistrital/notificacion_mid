package controllers

import (
	"container/list"
	"encoding/json"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/udistrital/notificacion_api/models"
)

// RoomController operations for Room
// RoomController operations for Room
type RoomController struct {
	beego.Controller
}

type Subscription struct {
	Archive []models.Event      // All the events from the archive.
	New     <-chan models.Event // New events coming in.
}

// URLMapping ...
func (c *RoomController) URLMapping() {

}

//funciones para la sala de notificaciones

func newEvent(ep models.EventType, user string, msg *models.Notificacion) models.Event {
	return models.Event{ep, user, int(time.Now().Unix()), msg}
}

func Join(id string, profiles []string, ws *websocket.Conn) {
	subscribe <- models.Subscriber{Id: id, Profiles: profiles, Conn: ws}
}

func Leave(user string) {
	unsubscribe <- models.Subscriber{Id: user}
}

var (
	// Channel for new join users.
	subscribe = make(chan models.Subscriber)
	// Channel for exit users.
	unsubscribe = make(chan models.Subscriber)
	// Send events here to publish them.
	publish = make(chan models.Event)
	// Long polling waiting list.
	waitingList = list.New()
	subscribers = list.New()

	connectionsId      = make(map[string]*websocket.Conn)
	connectionsProfile = make(map[string]map[string]*websocket.Conn)
)

// This function handles all incoming chan messages.
func chatroom() {
	for {

		select {
		case sub := <-subscribe:
			if !isUserExist(subscribers, sub.Id) {
				subscribers.PushBack(sub) // Add user to the end of list.
				connectionsId[sub.Id] = sub.Conn
				for _, profile := range sub.Profiles {
					if _, ok := connectionsProfile[profile]; ok {
						(connectionsProfile[profile])[sub.Id] = sub.Conn
					} else {
						connectionsProfile[profile] = make(map[string]*websocket.Conn)
						(connectionsProfile[profile])[sub.Id] = sub.Conn
					}
					beego.Info("Register profile:", profile)
				}
				// Publish a JOIN event.
				//publish <- newEvent(models.EVENT_JOIN, sub.Name, "")
				beego.Info("New user:", sub.Id, ";WebSocket:", sub.Conn != nil)
			} else {
				beego.Info("Old user:", sub.Id, ";WebSocket:", sub.Conn != nil)
			}

		case event := <-publish:
			// Notify waiting list.
			for ch := waitingList.Back(); ch != nil; ch = ch.Prev() {
				ch.Value.(chan bool) <- true
				waitingList.Remove(ch)
			}

			broadcastWebSocket(event)
			models.NewArchive(event)

			if event.Type == models.EVENT_MESSAGE {
				beego.Info("Message from", event.User, ";Content:", event.Content)
			}

		case unsub := <-unsubscribe:
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(models.Subscriber).Id == unsub.Id {
					subscribers.Remove(sub)
					delete(connectionsId, unsub.Id)
					for _, prof := range sub.Value.(models.Subscriber).Profiles {
						delete(connectionsProfile[prof], unsub.Id)
					}

					// Clone connection.
					ws := sub.Value.(models.Subscriber).Conn
					if ws != nil {
						ws.Close()
						publish <- newEvent(models.EVENT_LEAVE, unsub.Id, &models.Notificacion{})
						beego.Error("WebSocket closed:", unsub.Id)
					}
					publish <- newEvent(models.EVENT_LEAVE, unsub.Id, &models.Notificacion{}) // Publish a LEAVE event.
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
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(models.Subscriber).Id == user {
			return true
		}
	}
	return false
}

//++++++++++++++++++++++++++++++++++++++++
// Post ...
// @Title Create
// @Description create Room
// @Param	body		body 	models.Room	true		"body for Room content"
// @Success 201 {object} models.Room
// @Failure 403 body is empty
// @router / [post]
func (c *RoomController) Post() {
	var v map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		//push notificacion-------
		c.Ctx.Output.SetStatus(201)
		alert := models.Alert{Type: "success", Code: "S_544", Body: nil}
		c.Data["json"] = alert
	} else {
		alert := models.Alert{Type: "success", Code: "E_N001", Body: err.Error()}
		c.Data["json"] = alert
	}
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Room by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Room
// @Failure 403 :id is empty
// @router /:id [get]
func (c *RoomController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Room
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Room
// @Failure 403
// @router / [get]
func (c *RoomController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Room
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Room	true		"body for Room content"
// @Success 200 {object} models.Room
// @Failure 403 :id is not int
// @router /:id [put]
func (c *RoomController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Room
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *RoomController) Delete() {

}
