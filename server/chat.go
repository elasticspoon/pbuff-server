package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var chatRooms = make(map[int]*ChatRoom)

func Index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("templates/base.html", "templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
	}

	err = t.Execute(w, chatRooms)
	if err != nil {
		log.Printf("Error executing template: %s", err)
	}
}

func ShowChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("templates/base.html", "templates/chat.html")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
	}

	chatRoomId := chi.URLParam(r, "chatRoomId")
	id, err := strconv.Atoi(chatRoomId)
	if err != nil {
		log.Printf("Error converting chatRoomId to int: %s", err)
	}

	cr, ok := chatRooms[id]
	if !ok {
		log.Printf("Chat room does not exist: %d", id)
		http.Redirect(w, r, "/", http.StatusFound)
	}

	err = t.Execute(w, cr)
	if err != nil {
		log.Printf("Error executing template: %s", err)
	}
}

func CreateChatRoom(w http.ResponseWriter, r *http.Request) {
	cr := NewChatRoom()
	cr.Id = len(chatRooms)
	chatRooms[cr.Id] = cr

	go cr.Run()

	log.Printf("Created chat room: %d", cr.Id)
	http.Redirect(w, r, fmt.Sprintf("/chat/%d", cr.Id), http.StatusFound)
}

func JoinChatRoom(w http.ResponseWriter, r *http.Request) {
	chatRoomId := chi.URLParam(r, "chatRoomId")
	id, err := strconv.Atoi(chatRoomId)
	if err != nil {
		log.Printf("Error converting chatRoomId to int: %s", err)
	}

	cr, ok := chatRooms[id]
	if !ok {
		log.Printf("Chat room does not exist: %d", id)
		http.Redirect(w, r, "/", http.StatusFound)
	}

	log.Printf("Joining chat room: %d", id)
	ServeWs(cr, w, r)
}

func ServerAlive(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	_, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		log.Printf("Error: %s", err)
	}
}
