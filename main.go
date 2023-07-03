package main

import (
	"fmt"
	"log"
	"net"
)

var (
	activeChatRooms = make(map[string]*ChatRoom)
	activeUsers     = make(map[string]*User)
)

func main() {

	defaultConfig := getDefaults()                //getdefualts returns the path to the files of default users
	loadDefaultUsers(defaultConfig.Users)         //loads all the users in the csv file
	loadDefaultChatrooms(defaultConfig.Chatrooms) //loads all the chatrooms in the csv file
	//each connection/user has it's own goroutine and each chatroom has it's own goroutine.
	//remember clients != users != connections but they have a 1:1:1 relationship
	//fmt.Println("Here is the config object: ", defaultConfig.ConnectionType, ",", defaultConfig.Port)
	//port needs the colon on it

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
		//every connection gets it's own goroutine
		go handleConnection(&conn)
	}
}
