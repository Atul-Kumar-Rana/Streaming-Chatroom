package webrtc

import (
	"sync"

	"github.com/Atul-Kumar-Rana/Streaming-Chatroom/pkg/webrtc"
)

type Room struct {
	Peers *Peers
	hub   *chat.Hub
}
type Peers struct {
	ListLock     sync.RWMutex
	connections  []PeerConnectionState
	TreackLocals map[string]*webrtc.TrackLocalStaticRTP
}

func (p *Peers) DispatchKeyFrame() {

}
