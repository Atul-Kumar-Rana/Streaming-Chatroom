package webrtc

import (
	"log"
	"sync"
	"github.com/pion/rtcp"
	"github.com/Atul-Kumar-Rana/Streaming-Chatroom/pkg/chat"
	"github.com/pion/webrtc/v3"

	// "github.com/Atul-Kumar-Rana/Streaming-Chatroom/pkg/webrtc"
	"github.com/gofiber/websocket/v2"
)

type Room struct {
	Peers *Peers
	hub   *chat.Hub
}
type Peers struct {
	ListLock    sync.RWMutex
	connections []PeerConnectionState
	TrackLocals map[string]*webrtc.TrackLocalStaticRTP
}

type PeerConnectionState struct {
	PeerConnection *webrtc.PeerConnection
	websocket      *ThreadSafeWriter
}

type ThreadSafeWriter struct {
	conn  *websocket.Conn
	Mutex sync.Mutex
}

//	func (t *ThreadSafeWriter) WriteJSON(v interface())error{
//		t.Mutex.Lock()
//		defer t.Mutex.Unlock()
//		return t.Conn.WriteJSON(v)
//	}
func (t *ThreadSafeWriter) WriteJSON(v interface{}) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	return t.conn.WriteJSON(v)
}

func (p *Peers) AddTrack(t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	p.ListLock.Lock()
	defer func() {
		p.ListLock.Unlock()
		p.SignalPeerConnections()

	}()
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	p.TrackLocals[t.ID()] = trackLocal
	return trackLocal

}

func (p *Peers) RemoveTrack(t *webrtc.TrackLocalStaticRTP) {
	p.ListLock.Lock()
	defer func() {
		p.ListLock.Unlock()
		p.SignalPeerConnections()

	}()
	delete(p.TrackLocals, t.ID())
}
func (p *Peers) SignalPeerConnections() {
	p.ListLock.Lock()
	defer func() {
		p.ListLock.Unlock()
		p.DispatchKeyFrame()
	}()

	for i := 0; i < len(p.connections); i++ {
		if p.connections[i].PeerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
			p.connections = append(p.connections[:i], p.connections[i+1:]...)
			log.Println("Removed closed connection")
			i-- // Adjust index after removal
			continue
		}
		existingSenders := map[string]bool{}
		for _, sender := range p.connections[i].PeerConnection.GetSenders() {
			if sender.Track() == nil {
				continue
			}
			existingSenders[sender.Track().ID()] = true
			if _, ok := p.TrackLocals[sender.Track().ID()]; !ok {
				if err := p.connections[i].PeerConnection.RemoveTrack(sender); err != nil {
					log.Println("Error removing track:", err)
					continue
				}
			}
		}
		for _, receiver := range p.connections[i].PeerConnection.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}
			existingSenders[receiver.Track().ID()] = true
		}
		for trackID := range p.TrackLocals {
			if _, ok := existingSenders[trackID]; !ok {
				if _, err := p.connections[i].PeerConnection.AddTrack(p.TrackLocals[trackID]); err != nil {
					log.Println("Error adding track:", err)
					continue
				}
			}
		}
	}
}

// func (p *Peers) SignalPeerConnections(){
// p.ListLock.Lock()
// defer func ()  {
// 	p.ListLock.Unlock()
// 	p.DispatchKeyFrame()
// }()
//  attemptSync := func (tryAgain bool)  {
// 	for i := range p.Connections{
// 		if p.Connections[i].PeerConnection.ConnectionState()== webrtc.PeerConnectionStateClosed{
// 			p.Connections = append(p.Connections[:i],p.Connections[i+1:] ...)
// 			log.Println("a".p.Connection)
// 			return true
// 		}
// 		existingSenders := map[string]bool{}
// 		for _,sender := range p.Connections[i].PeerConnection.GetSenders(){
// 			if sender.Track()==nil{
// 				continue
// 			}
// 			existingSenders[sender.Track().ID()]==true
// 			if _,ok :=p.TrackLocals[sender.Track().ID()];!ok{
// 			if err := p.Connections[i].PeerConnection.RemoveTrack(sender);err!=nil{
// 				return  true
// 			}
// 		}
// 	}
// 	for _,receiver := range p.Connections[i].PeerConnection.GetReceivers(){
// 		if receiver.Track()==nil{
// 			continue
// 		}
// 		existingSenders[receiver.Track().ID()]=true

// 	for trackID := range p.TrackLocals{
// 		if _,ok := existingSenders[trackID]; !ok{
// 			if _,err := p.Connections[i].PeerConnection.AddTrack(p.TrackLocals[trackID]); err!=nil{
// 				return true
// 			}
// 		}
// 	}
// 	}
// 	}
//  }
// }

func (p *Peers) DispatchKeyFrame() {
	for _, peer := range p.connections {
		for _, sender := range peer.PeerConnection.GetSenders() {
			if track := sender.Track(); track != nil {
				// Only video tracks support keyframe requests
				if track.Kind() == webrtc.RTPCodecTypeVideo {
					_ = peer.PeerConnection.WriteRTCP([]rtcp.Packet{
   						 &rtcp.PictureLossIndication{MediaSSRC: sender.SSRC()},
						})
				}
			}
		}
	}

}

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}
