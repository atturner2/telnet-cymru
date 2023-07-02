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

	loadDefaultUsers()
	//loadDefaultChatrooms(defaultChatroomFilePath)
	//defaultChatroom := createChatRoom("a")
	//each connection/user has it's own goroutine and each chatroom has it's own goroutine.
	//remember clients != users != connections but they have a 1:1:1 relationship
	//go defaultChatroom.start()

	listener, err := net.Listen("tcp", ":23")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()

	fmt.Println("Telnet server started. Listening on port 23.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(&conn)
	}
}
