package webrtc

import (
	"sync"

	"github.com/Atul-Kumar-Rana/Streaming-Chatroom/pkg/webrtc"
	"github.com/gofiber/websocket/v2"
)


func RoomConn(c *websocket.Conn , p *Peers){
	var config webrtc.Configuration


	PeerConnection,err:=webrtc.NewPeerConnection(config)
	if(err!=nil){
		logPrintln(err)
		return
	}

	newPeer:=PeerConnectionSTate(
		PeerConnection: PeerConnection,
		WebSocket: &ThreadSafeWriter{},
		Conn: c,
		Mutex: sync.Mutex{},
	)
}