package handlers

import (
	w "github.com/Atul-Kumar-Rana/Streaming-Chatroom/pkg/webrtc"
	"github.com/Atul-Kumar-Rana/Streaming-Chatroom/pkg/chat"
	"github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber"
	"github.com/gofiber/websocket/v2"
	// "github.com/google/uuid"
	"golang.org/x/net/websocket"
)

func RoomChat(c *fiber.Ctx) error {
	return c.Render("chat", fiber.Map{}, "layouts/main")

}
func RoomChatWebsocket(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}
	w.RoomsLock.Lock()
	room := w.Rooms[uuid]
	w.RoomsLock.Unlock()
	if room == nil {
		return
	}
	if room.Hub == nil {
		return
	}
	chat.PeerChatConn(c.Conn, room.Hub)
}

func StreamChatWebsocket(c *websocket.Conn) {
	suuid := c.Params("suuid")
	if ssuid == "" {
		return
	}
	w.RoomsLock.Lock()
	if stream, ok := w.Streams[suuid]; ok {
		w.RoomsLock.Unlock()
		if stream.Hub == nil {
			hub := chat.NewHub()
			stream.Hub = hub
			go hub.Run()
		}
		chat.PeerChatConn(c.Conn, stream.Hub)
		return
	}
	w.RoomsLock.Unlock()
}
