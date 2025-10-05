package handlers

import (
	"fmt"
	"os"
	"time"
	w "github.com/Atul-Kumar-Rana/Streaming-Chatroom/pkg/webrtc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	guuid "github.com/google/uuid"
	"golang.org/x/net/websocket"
)

func RoomCreate(c *fiber.Ctx) error{
	return c.Redirect(fmt.Sprintf("/room/%s",guuid.New().String()))

}
func Room(c *fiber.Ctx)error{
	uuid :=c.Params("uuid")
if uuid ==""{
	c.Status(400)
	return nil
}
uuid,suuid,_ :=createOrGetRoom(uuid)
}

func Roomwebsocket(c *websocket.conn){
	uuid := c.Parse("uuid")
	if uuid==""{
		return
	}
	//  this createOrGEtRoom fx will create room if that room id doesnt exist 
	_,_,room := createOrGetRoom(uuid)  
	w.RoomConn(c,room.Peers)
}

func createOrGetRoom(uuid string)(string,string,*w Room){
   
}

func RoomViewerWebsocket(c *websocket.Conn){

}
func romViewerConn(c *websocket.Conn , p *w.Peers){

}

typr websocketMessage struct{
	Event string `json:"event"`
	Data string `json:"data"`
}
