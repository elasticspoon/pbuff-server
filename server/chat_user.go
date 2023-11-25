package server

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	ChatRoom *ChatRoom
	Conn     *websocket.Conn
	Send     chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *Client) Read() {
	defer func() {
		c.ChatRoom.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			c.ChatRoom.Unregister <- c
			c.Conn.Close()
			break
		}
		message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))
		c.ChatRoom.Broadcast <- message
	}
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func ServeWs(chatRoom *ChatRoom, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{ChatRoom: chatRoom, Conn: conn, Send: make(chan []byte, 256)}
	client.ChatRoom.Register <- client

	go client.Write()
	go client.Read()
}
