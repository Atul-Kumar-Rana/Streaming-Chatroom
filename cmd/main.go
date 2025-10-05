package main

import (
	"log"

	"github.com/Atul-Kumar-Rana/Streaming-Chatroom/internal/server"
)
func main(){
	if err := server.Run();err!=nil{
	log.Fatal(err.Error())
}

}