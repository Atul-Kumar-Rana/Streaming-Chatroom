package handlers

import (
	"fmt"
	"os"
	"time"

	w "github.com/Atul-Kumar-Rana/Streaming-Chatroom/pkg/webrtc"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/websocket"
)

func Stream(c *fiber.Ctx) error {
	suuid := c.Params("suuid")
	if suuid == "" {
		c.Status(400)
		return nil
	}
	ws := "ws"
	if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
		ws = "wss"
	}
	w.RoomsLock.Lock()
	if _, ok := w.Stream[suuid]; ok {
		w.RoomsLock.Unlock()
		return c.Render("stream", fiber.Map{
			"StreamWebSocketAddr": fmt.Sprintf("%s://%s/stream/%s/websocket", ws, c.Hostname(), suuid),
			"ViewerWebSOcketAddr": fmt.Sprintf("%s://%s/stream/%s/chat/websocket", ws, c.Hostname(), suuid),
			"CHatWebSocketAddr":   fmt.Sprintf("%s://%s/stream/%s/viewer/websocket", ws, c.Hostname(), suuid),
			"Typr":                "stream",
		}, "layouts/main")
	}
	w.RoomsLock.Unlock()
	return c.Render("stream", fiber.Map{
		"NoStream": "true",
		"Leave":    "true",
	}, "layouts/main")
}

func StreamWebsocket(c *websocket.Conn) {
	suuid := c.Params("suuid")
	if suuid == "" {
		return

	}
	w.RoomsLock.Lock()
	if Stream, ok := w.Stream[uuid]; ok {
		w.RoomsLock.Unlock()
		w.StreamConn(c, stream.Peers)
		return

	}
	w.RoomsLock.Unlock()
}

func StreamViewerWebsocket(c *websocket.conn) {
	suuid := c.Params("suuid")
	if suuid == "" {
		return

	}
	w.RoomsLock.Lock()
	if Stream, ok := w.Stream[uuid]; ok {
		w.RoomsLock.Unlock()
		w.ViewerConn(c, stream.Peers)
		return

	}
	w.RoomsLock.Unlock()

}
func ViewerConn(c *websocket.Conn, p *w.Peers) {
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
			w.Write([]byte(fmt.Sprintf("%d",len(p.Connections))))
		}

	}

}
