package chat


type Hub struct{
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}
// fxn to create new hub instance
// user can come and go out of room

func NewHub() *Hub{
	return &Hub{
		broadcast:make(chan []byte) ,
		register:make(chan *Client) ,
		unregister: make(chan *Client) ,
		client: make(map[*Client]bool),
	}

}

func (h *Hub) Run(){
	for{
		select{
		case client := <-h.register:
			h.clients[client]=true
		case client := <-h.unregister:
			if _,ok := h.clients[client];ok{
				delete(h.clients,client)
				close(client.Send)
			}
		case message := <-h.broadcast:
			for client := range h.clients{
				select{
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients,client)
				}
			}


		
		}
	
	}
}