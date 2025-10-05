package server

import (
	"flag"
	"os"
	"time"
	"github.com/gofiber/template/html/v2"
	"github.com/Atul-Kumar-Rana/Streaming-Chatroom/internal/handlers"
	w "github.com/Atul-Kumar-Rana/Streaming-Chatroom/pkg/webrtc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	// "github.com/gofiber/template/html"
	"github.com/gofiber/websocket/v2"
	"golang.org/x/net/websocket"
)
var(
	addr = flag.String("addr" ,":"+os.Getenv("PORT"),"")
	cert = flag.String("cert","","")
	key=flag.String("key","","")
)
func  Run() error  {
	flag.Parse()
	if *addr==":"{
		*addr=":8080"
	}

//  fiber similar to express .... here using html page
	engine := html.New("./views",".html")
	// using fiber engine 
	app    := fiber.New(fiber.Config{Views: engine})

	// middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// routes
	app.Get("/",handlers.Welcome)
	app.Get("/room/create",handlers.RoomCreate)
	app.Get("/room/:uuid",handlers.Room)
	app.Get("/room/:uuid/websocket",websocket.New(handlers.Roomwebsocket,websocket.Config{
		HandshakeTimeout:10*time.Second,
	}))
	app.Get("/room/:uuid/chat",handlers.RoomChat)
	app.Get("/room/:uuid/chat/websocket",websocket.New(handlers.RoomChatWebsocket))
	app.Get("/room/:uuid/viewer/websocket",websocket.New(handlers.RoomViewerWebsocket))
	app.Get("/stream/:ssuid",handlers.Stream)
	app.Get("/stream/:ssuid/websocket", websocket.New(handlers.StreamWebsocket,websocket.Config{
		HandshakeTimeout:10*time.Second,
	}))
	app.Get("/stream/:ssuid/chat/websocket",websocket.New(handlers.StreamChatWebsocket))
	app.Get("/stream/:ssuid/viewer/websocket",websocket.New(handlers.StreamViewerWebsocket))
	app.Static("/","./assets")


	// creating go rutine

	w.Room=make(map[string]*w.Room)
	w.Stream=make(map[string]*w.Room)
	go dispatchKeyFrames()
	// if dont have any certificate this this would work 
	if *cert !=""{
		return app.ListenTLS(*addr,*cert,*key)
	}
	return app.Listen(*addr)
	
	
}

func dispatchKeyFrames(){
		for range time.NewTicker(time.Second *3).c{
			for _,room :=range w.Rooms{
				room.Peers.DispatchKeyFrame()
			}
		}
	}