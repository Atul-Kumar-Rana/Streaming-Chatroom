package chat

import (
	"bytes"
	"log"

	"github.com/fasthttp/websocket"
	"golang.org/x/text/message"

	"time"
)

const(
	writeWait = 10*time.Second
	pongWait  = 60*time.Second
	pingPeriod  = (pongWait *9)/10
	maxMessageSize  =  512

)

var(
	newLine = [] byte{'\n'}
	space = []byte{' '}
)


type Client struct{
	Hub *Hub
	Conn *websocket.Conn
	Send chan []byte

}

var upgrader = websocket.fastHTTPUpggrader{
	ReadBufferSize : 1024,
	WriteBufferSize: 1024,
}

func ( c *Client) readPump(){
	defer func ()  {
		c.Hub.unregister<-c
		c.Conn.close()
		
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReaddDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func (string)error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait));
		return nil
	})
	for{
		_,message,err :=c.Conn.ReadMessage()
		if err!=nil{
			if websocket.IsUnexpectedCloseError(err,websocket.CloseGoingAway, websocket.CloseAbnormalClosure){
				log.Printf("error: %v",err)

			}
			break;
		}
		message= bytes.TrimSpace(bytes.Replace(message,newLine,space,-1))
		c.Hub.broadcast<-message
	}

}

func ( c *Client) writePump(){
	ticker := time.NewTicker(pingPeriod)
	defer func ()  {
		ticker.Stop()
		c.Conn.Close()
	}()
	for{
		select{
		case message, ok := <-c.Send:
			c.Conn.SetWriterDeadline(time.Now().Add(writeWait))
			if !ok{
				return
			}
			w,err :=c.Conn.NextWriter(websocket.TextMessage)
			if err!=nil{
				return
			}
			// for multiple message would iterate over it and send those messages 
			w.Write(message)
			n:=len(c.Send)
			for i:=0; i<n ; i++{
				w.Write(newLine)
				w.Write(<-c.Send)
			}
			if err:= w.Close(); err!=nil{
				return
			}
		case <- ticker.C:
			c.Conn.SetWriterDeadline(time.Now().Add(writeWait))
			if err:= c.Conn.WriteMessae(websocket.PingMessage,nil);err!=nil{
				return
			}

		}
	}


}

func PeerChatConn( c *websocket.conn , hub *Hub)  {
	client := &Client{Hub:hub , Conn: c , Send:make(chan []byte,256)}
	client.Hub.register <- client
	go client.writePump()
	client.readPump()
}