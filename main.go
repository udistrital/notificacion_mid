package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"github.com/udistrital/notificacion_api/models"
	"log"
	"net/http"
	"strings"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
)

var addr = flag.String("addr", ":"+os.Getenv("APP_PORT"), "http service address")
var indexFile = "index.html"
var h models.Hub

type RequestMessage struct {
	Type  		string 		`json:"type"`
	Id 				string 		`json:"id"`
	Profile		[]string	`json:"profile"`
	Message 	string 		`json:"message"`
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	// fmt.Println(r.URL.Path)
	// fmt.Printf("%#v\n", decoder)
	// data := strings.Split(r.URL.Path,"/")
	// data := strings.Split("/1/admin/ws","/")
	fmt.Println("Values %s",values["id"][0])
	id := values["id"][0]

	profiles := strings.Split(values["profile"][0],"_")
	// Upgrade the HTTP connection to a websocket. TODO: check origin
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}
	c := models.NewConnection(ws)
	h.Register(models.ConnValues{C: c, Id: id, Profile: profiles})
	defer func() { h.Unregister(models.ConnValues{C: c, Id: id, Profile: profiles})}()
	go c.Writer()
	c.Reader(h)
}

func unregister(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values["id"][0]
	profiles := strings.Split(values["profile"][0],"_")
	h.Unregister(models.ConnValues{C: nil, Id: id, Profile: profiles})
}

func push_notification(w http.ResponseWriter, r *http.Request) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msg RequestMessage
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	switch msg.Type {
	case "personal":
		h.SendPersonalMessage(models.SendingMessage{M:[]byte(msg.Message), ConnValues: models.ConnValues{C:nil, Id:msg.Id, Profile:nil}})
	case "profile":
		h.SendProfileMessage(models.SendingMessage{M:[]byte(msg.Message), ConnValues: models.ConnValues{C:nil, Id:msg.Id, Profile:msg.Profile}})
	}

}

func get_notification(w http.ResponseWriter, r *http.Request) {
	// values := r.URL.Query()
}

func main() {
	flag.Parse()
	h = models.NewHub()

	go h.Run()

	// http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/register", wsHandler)
	http.HandleFunc("/unregister", unregister)
	http.HandleFunc("/push_notification", push_notification)
	http.HandleFunc("/get_notification", get_notification)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
