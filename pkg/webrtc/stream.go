package webrtc

import (
	"fmt"

	"github.com/pion/webrtc/v3"
)

// Stream represents a WebRTC stream for video calling and chatting.
type Stream struct {
	PeerConnection *webrtc.PeerConnection
	DataChannel    *webrtc.DataChannel
}

// NewStream initializes a new WebRTC stream.
func NewStream() (*Stream, error) {
	// WebRTC configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create peer connection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %w", err)
	}

	// Create data channel for chat
	dataChannel, err := peerConnection.CreateDataChannel("chat", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create data channel: %w", err)
	}

	stream := &Stream{
		PeerConnection: peerConnection,
		DataChannel:    dataChannel,
	}

	// Set up handlers
	stream.setupHandlers()

	return stream, nil
}

// setupHandlers sets up event handlers for the stream.
func (s *Stream) setupHandlers() {
	s.DataChannel.OnOpen(func() {
		fmt.Println("Data channel opened")
	})

	s.DataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Received chat message: %s\n", string(msg.Data))
		// Handle chat message here
	})

	s.PeerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		fmt.Printf("Received track: %s\n", track.Kind().String())
		// Handle incoming video/audio track here
	})
}

// CreateOffer generates an SDP offer for signaling.
func (s *Stream) CreateOffer() (webrtc.SessionDescription, error) {
	offer, err := s.PeerConnection.CreateOffer(nil)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}
	err = s.PeerConnection.SetLocalDescription(offer)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}
	return offer, nil
}

// SetRemoteDescription sets the remote SDP.
func (s *Stream) SetRemoteDescription(desc webrtc.SessionDescription) error {
	return s.PeerConnection.SetRemoteDescription(desc)
}

// AddTrack adds a local media track (video/audio) to the peer connection.
func (s *Stream) AddTrack(track *webrtc.TrackLocalStaticSample) (*webrtc.RTPSender, error) {
	return s.PeerConnection.AddTrack(track)
}

// SendChatMessage sends a chat message over the data channel.
func (s *Stream) SendChatMessage(message string) error {
	if s.DataChannel != nil && s.DataChannel.ReadyState() == webrtc.DataChannelStateOpen {
		return s.DataChannel.SendText(message)
	}
	return fmt.Errorf("data channel not open")
}

// Close closes the stream and cleans up resources.
func (s *Stream) Close() error {
	if s.DataChannel != nil {
		s.DataChannel.Close()
	}
	return s.PeerConnection.Close()
}
