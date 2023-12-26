package main

import (
	"net/http"

	"github.com/elasticspoon/pbuff-server/server"
	"github.com/go-chi/chi"
	// "github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	// r.Use(middleware.Logger)

	// r.Get("/livereload.js", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./assets/livereload.js")
	// })
	// r.Get("/alive", server.ServerAlive)

	r.Get("/", server.Index)
	r.Get("/new", server.CreateChatRoom)
	r.Get("/chat/{chatRoomId}", server.ShowChat)

	r.Get("/ws/{chatRoomId}", server.JoinChatRoom)
	r.Get("/chat.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/chat.js")
	})

	http.ListenAndServe(":3000", r)
}
