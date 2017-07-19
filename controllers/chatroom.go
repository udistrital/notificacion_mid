package controllers

import (
	"container/list"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/udistrital/notificacion_api/models"
)

type Subscription struct {
	Archive []models.Event      // All the events from the archive.
	New     <-chan models.Event // New events coming in.
}

func newEvent(ep models.EventType, user string, profiles []string, msg string) models.Event {
	return models.Event{ep, user, profiles, int(time.Now().Unix()), msg}
}

func Join(user string, profiles []string, ws *websocket.Conn) {
	subscribe <- Subscriber{Name: user, Profiles: profiles, Conn: ws}
}

func Leave(user string) {
	unsubscribe <- user
}

type Subscriber struct {
	Name     string
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
			if !isUserExist(subscribers, sub.Name) {
				subscribers.PushBack(sub) // Add user to the end of list.
				connectionsId[sub.Name] = sub.Conn
				fmt.Println(connectionsId)
				for _, profile := range sub.Profiles {
					if _, ok := connectionsProfile[profile]; ok {
						(connectionsProfile[profile])[sub.Name] = sub.Conn
					} else {
						connectionsProfile[profile] = make(map[string]*websocket.Conn)
						(connectionsProfile[profile])[sub.Name] = sub.Conn
					}
					beego.Info("Register profile:", profile)
				}
				// Publish a JOIN event.
				publish <- newEvent(models.EVENT_MESSAGE, sub.Name, sub.Profiles, "Se unio al ws")
				beego.Info("New user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			} else {
				publish <- newEvent(models.EVENT_MESSAGE, sub.Name, sub.Profiles, "reload") // Publish a LEAVE event.
				beego.Info("Old user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			}
		case event := <-publish:

			broadcastWebSocket(event)
			models.NewArchive(event)

			if event.Type == models.EVENT_MESSAGE {
				beego.Info("Message from", event.User, ";Content:", event.Content)
			}
		case unsub := <-unsubscribe:
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).Name == unsub {
					subscribers.Remove(sub)
					delete(connectionsId, unsub)
					// Clone connection.
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
						beego.Error("WebSocket closed:", unsub)
					}
					publish <- newEvent(models.EVENT_LEAVE, unsub, nil, "logout") // Publish a LEAVE event.
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
		if sub.Value.(Subscriber).Name == user {
			return true
		}
	}
	return false
}
