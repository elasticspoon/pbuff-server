package server

type ChatRoom struct {
	Users      map[*Client]bool
	unregister chan *Client
	broadcast  chan []byte
	register   chan *Client
	Id         int
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		Users:      make(map[*Client]bool),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
	}
}

func (cr *ChatRoom) Run() {
	for {
		select {
		// if a client registers, add them to the chat room
		case client := <-cr.register:
			cr.Users[client] = true
			// if a client unregisters:
			// - remove them from the chat room
			// - close their send channel
			// - if there are no more users in the chat room, set the chat room to nil
		case client := <-cr.unregister:
			if _, ok := cr.Users[client]; ok {
				delete(cr.Users, client)
				close(client.send)
				// if len(cr.Users) == 0 {
				// 	chatRooms[cr.Id] = nil
				// }
			}
			// if a message is broadcasted, loop through all the users in the chat room
		case message := <-cr.broadcast:
			for client := range cr.Users {
				select {
				// if the user's send channel is available, send the message
				case client.send <- message:
					// if the user's send channel is not available
					// - remove them from the chat room.
					// - close their send channel
				default:
					close(client.send)
					delete(cr.Users, client)
				}
			}
		}
	}
}
