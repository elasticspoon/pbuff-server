package server

type ChatRoom struct {
	Users      map[*Client]bool
	Unregister chan *Client
	Broadcast  chan []byte
	Register   chan *Client
	Id         int
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		Users:      make(map[*Client]bool),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
	}
}

func (c *ChatRoom) Run() {
	for {
		select {
		case client := <-c.Register:
			c.Users[client] = true
		case client := <-c.Unregister:
			if _, ok := c.Users[client]; ok {
				delete(c.Users, client)
				close(client.Send)
			}
		case message := <-c.Broadcast:
			for client := range c.Users {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(c.Users, client)
				}
			}
		}
	}
}
