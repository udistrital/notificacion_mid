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
	"bytes"
	"encoding/binary"
)

var addr = flag.String("addr", ":"+os.Getenv("APP_PORT"), "http service address")
var indexFile = "index.html"
var h models.Hub
var url_crud = os.Getenv("URL_CRUD")

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
	var msg models.Notificacion
	var response interface{} //models.Notificacion
	var tipo_notificacion string

	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if msg.UsuarioDestino == 0{
			tipo_notificacion = "profile"
	}else{
			tipo_notificacion = "personal"
	}

	models.SendJson(url_crud+"/v1/notificacion","POST",&response, &msg)
	// fmt.Println("err: ", err5)
	// fmt.Println("Error para verlo: ", response)

	buf := &bytes.Buffer{}
	b, err = json.Marshal(response)
    if err != nil {
        fmt.Printf("Error: %s", err)
        return
    }

	err = binary.Write(buf, binary.BigEndian, &b)
	if err != nil {
			panic(err)
	}
	mensaje := buf.Bytes()
	switch tipo_notificacion {
	case "personal":
		h.SendPersonalMessage(models.SendingMessage{M:[]byte(mensaje), ConnValues: models.ConnValues{C:nil, Id:fmt.Sprint(msg.UsuarioDestino), Profile:nil}})
	case "profile":
		var profiles [] string
		profiles[0] = fmt.Sprint(msg.PerfilDestino)
		h.SendProfileMessage(models.SendingMessage{M:[]byte(mensaje), ConnValues: models.ConnValues{C:nil, Id:"", Profile:profiles}})
	}

}

func get_notification(w http.ResponseWriter, r *http.Request) {
	// values := r.URL.Query()
}

func main() {
	flag.Parse()
	h = models.NewHub()

	go h.Run()

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/register", wsHandler)
	http.HandleFunc("/unregister", unregister)
	http.HandleFunc("/push_notification", push_notification)
	// http.HandleFunc("/get_notification", get_notification)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
