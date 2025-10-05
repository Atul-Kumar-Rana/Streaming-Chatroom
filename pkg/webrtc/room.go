package webrtc

import (
	"log"
    "sync"

    "github.com/pion/webrtc/v3"

	// "github.com/Atul-Kumar-Rana/Streaming-Chatroom/pkg/webrtc"
	"github.com/gofiber/websocket/v2"
)


func RoomConn(c *websocket.Conn, p *Peers) {
    var config webrtc.Configuration

    peerConnection, err := webrtc.NewPeerConnection(config)
    if err != nil {
        log.Println(err)
        return
    }

    newPeer := PeerConnectionState{
        PeerConnection: peerConnection,
        websocket:      &ThreadSafeWriter{conn: c, Mutex: sync.Mutex{}},
    }
    // You may want to add newPeer to p.Connections here
}