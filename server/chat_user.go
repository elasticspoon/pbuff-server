package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	chatRoom *ChatRoom
	conn     *websocket.Conn
	send     chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// client reading from the websocket connection
// and writing to the chat room broadcast channel
func (c *Client) read() {
	// cleanup
	defer func() {
		c.chatRoom.unregister <- c
		c.conn.Close()
	}()

	for {
		// loops reading messages until an error occurs
		// then defer closes and unregisters the client
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// i dont think this is necessary will happen in the defer
			// c.chatRoom.unregister <- c
			// c.conn.Close()
			break
		}
		msg := &Message{}
		err = json.Unmarshal(message, msg)
		if err != nil {
			log.Printf("Error unmarshalling message: %s, error: %s", message, err)
			continue
		}

		body := []byte(msg.Body)
		body = bytes.TrimSpace(bytes.ReplaceAll(body, newline, space))
		msg.Body = string(body)

		m, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling message: %s", err)
			continue
		}
		c.chatRoom.broadcast <- m
	}
}

// this function writes to the websocket connection from chatRoom send channel
func (c *Client) write() {
	// sender does not unregister the client on close
	defer func() {
		c.conn.Close()
	}()

	// read from the send channel until it is closed
	for message := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	// for {
	// 	select {
	// 	// blocks until a message is received or the channel is closed
	// 	// send is written to by the chat room broadcast channel
	// 	case message, ok := <-c.send:
	// 		// if the channel is closed, close the websocket connection
	// 		if !ok {
	// 			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	// 			return
	// 		}
	// 		// write the message to the websocket connection
	// 		c.conn.WriteMessage(websocket.TextMessage, message)
	// 	}
	// }
}

func ServeWs(chatRoom *ChatRoom, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{chatRoom: chatRoom, conn: conn, send: make(chan []byte, 256)}
	client.chatRoom.register <- client

	go client.write()
	go client.read()
}
