package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

type ChatRoom struct {
	Name          string
	Messages      chan string
	Join          chan Client
	LeaveChatroom chan Client
	clients       map[string]Client
	Broadcast     chan string
	clientsMux    sync.RWMutex
}

func createChatRoom(roomName string) ChatRoom {
	mu.Lock()
	defer mu.Unlock()

	chatRoom := ChatRoom{
		Name:          roomName,
		Messages:      make(chan string),
		Join:          make(chan Client),
		LeaveChatroom: make(chan Client),
		clients:       make(map[string]Client),
		clientsMux:    sync.RWMutex{},
	}

	go chatRoom.start()

	activeChatRooms[roomName] = chatRoom
	return chatRoom
}

func (cr ChatRoom) start() {
	for {
		select {
		case client := <-cr.Join:
			fmt.Println("joining chatroom")
			cr.clientsMux.Lock()
			fmt.Println("here is the username being added: ", client.user.Username)
			cr.clients[client.user.Username] = client
			fmt.Println("Here are all the clients: ", clients)
			fmt.Println("Here are all the users in the chatroom: ", clients)
			cr.clientsMux.Unlock()
		case client := <-cr.LeaveChatroom:
			fmt.Println("leaving chatroom")
			cr.clientsMux.Lock()
			delete(cr.clients, client.user.Username)
			cr.clientsMux.Unlock()
		case message := <-cr.Messages:
			cr.clientsMux.RLock()
			user, parsedMessage := extractMessage(message)
			for username := range cr.clients {

				fmt.Println("Sending a message to user: ", cr.clients[username])
				client := cr.clients[username]
				if username != user {
					fmt.Println("Sending a message to user: ", cr.clients[username])
					client.send(parsedMessage)
				} else {
					fmt.Println("Skipping the user that sent the message")
				}
			}
			cr.clientsMux.RUnlock()

		}
	}
}

func handleJoinRoom(client *Client) {
	for {
		fmt.Fprint(client.writer, "Enter the name of the chat from the list below you want to join room you want to join: \n")
		for key := range activeChatRooms {
			fmt.Fprintf(client.writer, "%s\n", key)
		}
		client.writer.Flush()

		roomName, err := client.reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			return
		}

		// Check if the chat room exists
		roomName = strings.TrimSpace(roomName)
		chatRoom, exists := activeChatRooms[roomName]
		if !exists {
			fmt.Fprintf(client.writer, "The chat room '%s' does not exist.\n", roomName)
			client.writer.Flush()
			return
		}

		// Join the chat room
		fmt.Println("Here is the client name that SHOULD be joining: ", client.user.Username)
		chatRoom.Join <- *client
		//set the chatroom name on the user
		client.room = chatRoom
		fmt.Fprintf(client.writer, "You have joined the chat room '%s'.\n", roomName)
		client.writer.Flush()
		handleChatRoomInteraction(client)

		return
	}
}

func handleChatRoomInteraction(client *Client) {
	for {
		// Read input from the client
		input, err := client.Receive() // Assume Receive() reads user input from the client's connection

		if err != nil {
			// Handle error if unable to read input
			client.send("Error reading input. Please try again.")
			continue
		}

		trimmedInput := strings.TrimSpace(input)

		if trimmedInput == "exit" {
			client.send("Leaving the " + client.room.Name + " chat room.")
			client.room.LeaveChatroom <- *client // Leave the chat room
			client.room.Name = ""                // Set the current chat room to nil
			// Exit the loop and return to the main menu
			return
		} else {
			//username and timestamp are added onto the front of the message
			//so other users can see whos'd sending what, and the message handler
			//knows who sent what message so it can not send the message to the user that sent it
			trimmedInput = client.user.Username + ": " + trimmedInput
			fmt.Println("calling the message sending pipe")

			client.room.Messages <- trimmedInput // Leave the chat room

		}
	}
	fmt.Println("The handlechatroominteraction is returning!")
	return
}

func extractMessage(originalMessage string) (username, message string) {
	index := strings.Index(originalMessage, ":")
	if index >= 0 {
		username = originalMessage[:index]
	} else {
		username = originalMessage
	}
	//we want the original message returned so we can see who sent it in the chatroom
	return username, originalMessage
}
