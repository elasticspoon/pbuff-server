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

var (
	chatRooms  = make(map[int]*ChatRoom)
	chatRoomId = 0
)

type Message struct {
	Body string `json:"body"`
	Time int    `json:"time"`
}

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

	if len(chatRooms) == 0 || id > len(chatRooms) || id < 0 || chatRooms[id] == nil {
		log.Printf("Chat room does not exist: %d", id)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	cr := chatRooms[id]

	err = t.Execute(w, cr)
	if err != nil {
		log.Printf("Error executing template: %s", err)
	}
}

func DeleteChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	chatRoomId := chi.URLParam(r, "chatRoomId")
	id, err := strconv.Atoi(chatRoomId)
	if err != nil {
		log.Printf("Error converting chatRoomId to int: %s", err)
	}

	if _, ok := chatRooms[id]; !ok {
		log.Printf("Chat room does not exist: %d", id)
		http.Redirect(w, r, "/", http.StatusFound)
	}

	delete(chatRooms, id)
	// http.Redirect(w, r, "/", http.StatusFound)
}

func CreateChatRoom(w http.ResponseWriter, r *http.Request) {
	cr := NewChatRoom()

	cr.Id = chatRoomId
	chatRooms[chatRoomId] = cr
	chatRoomId++

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

	if len(chatRooms) == 0 || id > len(chatRooms) || id < 0 || chatRooms[id] == nil {
		log.Printf("Chat room does not exist: %d", id)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	cr := chatRooms[id]

	// log.Printf("Joining chat room: %d", id)
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
