package main

import (
	"fmt"
	"log"
	"net"
)

var (
	activeChatRooms = make(map[string]*ChatRoom)
	activeUsers     = make(map[string]User)
)

func main() {

	defaultConfig := getDefaults()
	loadDefaultUsers(defaultConfig.Users)
	loadDefaultChatrooms(defaultConfig.Chatrooms)
	//each connection/user has it's own goroutine and each chatroom has it's own goroutine.
	//remember clients != users != connections but they have a 1:1:1 relationship
	//go defaultChatroom.start()

	listener, err := net.Listen(defaultConfig.ConnectionType, defaultConfig.Port)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()

	fmt.Println("Telnet server started. Listening on port ", defaultConfig.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		//every user that logs in gets their own goroutine
		go handleConnection(&conn)
	}
}
