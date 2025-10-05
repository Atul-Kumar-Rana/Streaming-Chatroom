package handlers

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/Atul-Kumar-Rana/Streaming-Chatroom/pkg/chat"
	"github.com/Atul-Kumar-Rana/Streaming-Chatroom/pkg/webrtc"
	w "github.com/Atul-Kumar-Rana/Streaming-Chatroom/pkg/webrtc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	guuid "github.com/google/uuid"
	"golang.org/x/net/websocket"
)

func RoomCreate(c *fiber.Ctx) error {
	return c.Redirect(fmt.Sprintf("/room/%s", guuid.New().String()))

}
func Room(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		c.Status(400)
		return nil
	}
	ws := "ws"
	if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
		ws = "wss"
	}
	uuid, suuid, _ := createOrGetRoom(uuid)
	return c.Render("peer", fiber.Map{
		"RoomWebSocketAddr":   fmt.Sprintf("%s://%s/room%s/websocket", ws, c, c.Hostname(), uuid),
		"RoomLink":            fmt.Printf("%s://%s/room/%s", c.Protocol(), c.Hostname(), uuid),
		"ChatWebSocketAddr":   fmt.Printf("%s://%s/room/%s/chat/websocket", ws, c.Hostname(), uuid),
		"ViewerWebSocketAddr": fmt.Sprintf("%s://%s/room/%s/viewer/websocket", ws, c.Hostname(), uuid),
		"StreamLink":          fmt.Printf("%s://%s/stream/%s", c.Protocol(), c.Hostname(), suuid),
		"Type":                "room",
		}"layouts/main")
}

func Roomwebsocket(c *websocket.conn) {
	uuid := c.Parse("uuid")
	if uuid == "" {
		return
	}
	//  this createOrGEtRoom fx will create room if that room id doesnt exist
	_, _, room := createOrGetRoom(uuid)
	w.RoomConn(c, room.Peers)
}

func createOrGetRoom(uuid string) (string, string, *w.Room) {
	w.RoomsLock.Lock()
	defer w.RoomsLock.Unlock()
	h := sha256.New()
	h.Write([]byte(uuid))
	suuid := fmt.Sprintf("%x", h.Sum(nil))

	if room := w.Rooms[uuid]; room != nil {
		if _, ok := w.Streams[suuid]; !ok {
			w.Streams[uuid] = room

		}
		return uuid, suuid, room
	}
	hub := chat.NewHub()
	p := &w.Peers{}
	p.TreackLocals = make(map[string]*webrtc.TreackLocalStaticRTP)
	room := &w.Room{
		Peers: p,
		Hub: hub,}
	go hub.Run()
	return uuid,suuid , room
	
}

func RoomViewerWebsocket(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}

	w.RoomsLock.Lock()
	if peer, ok := w.Rooms[uuid]; ok {
		w.RoomsLock.Unlock()
		roomViewerConn(c, peer.Peers)
		return
	}
}

// sees next connection fxn from fiber
func roomViewerConn(c *websocket.Conn, p *w.Peers) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer c.Close()

	for {
		select {
		case <-ticker.c:
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Writer([]byte(fmt.Sprintf("%d", len(p.Connection))))
		}
	}

}

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}
