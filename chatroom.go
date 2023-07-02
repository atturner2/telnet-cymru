package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
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
	//the mutex here is for when we access the activechatrooms object, it just stores
	//pointers to all of the existing chatrooms
	chatRoom := ChatRoom{
		Name:          roomName,
		Messages:      make(chan string),
		Join:          make(chan Client),
		LeaveChatroom: make(chan Client),
		clients:       make(map[string]Client),
		clientsMux:    sync.RWMutex{},
	}

	go chatRoom.start()
	activeChatRooms[roomName] = &chatRoom
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
			fmt.Println("Here are all the clients: ", cr.clients)
			fmt.Println("Here are all the users in the chatroom: ", cr.clients)
			cr.clientsMux.Unlock()
		case client := <-cr.LeaveChatroom:
			fmt.Println("leaving chatroom")
			cr.clientsMux.Lock()
			delete(cr.clients, client.user.Username)
			cr.clientsMux.Unlock()
		case message := <-cr.Messages:
			//this is where we send all the messages, skip the user
			//that sent the message, and write the message to the log files
			//we could also do this in the send function of the users,
			//but its better to do it here so you don't have to work to avoid duplicates
			//and don't have to worry as much about race conditions
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
			handleWriteToLogFile(cr.Name, parsedMessage)
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
		client.room = *chatRoom
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
		//this should never happen
		username = originalMessage
	}
	originalMessage = addTimeStamp(originalMessage)
	//we want the original message returned so we can see who sent it in the chatroom
	return username, originalMessage
}

// this is just to add the timestamp to the messages, for both the user receiving them
// and for the log files.
func addTimeStamp(message string) (finalMessage string) {
	currentTime := time.Now()
	timestamp := currentTime.Format("2006.01.02 15:04:05")
	return fmt.Sprintf("[%s] %s", timestamp, message)
}

// note that this handles the case where the user wants to create the chatroom,
// the chatroom from the default .csv files are just created upon startup without any user input
func handleCreateChatRoom(client *Client) {
	for {
		fmt.Fprint(client.writer, "Please enter the name of the Chatroom you would like to create: ")
		client.writer.Flush()

		chatRoomName, err := client.reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading chatroom name:", err)
			return
		}

		chatRoomName = strings.TrimSpace(chatRoomName)
		if chatRoomName == "" {
			fmt.Fprintln(client.writer, "Username cannot be empty. Please try again.")
			client.writer.Flush()
			continue
		}

		if userExists(chatRoomName) {
			fmt.Fprintln(client.writer, "Username already exists. Please choose a different username.")
			client.writer.Flush()
			continue
		}

		createChatRoom(chatRoomName)

		fmt.Fprintf(client.writer, "Chatroom Created, %s!\n", chatRoomName)
		client.writer.Flush()

		return

	}
}

func handleWriteToLogFile(chatRoomName, message string) {
	fileName := chatRoomName + ".log"
	//this will open the log file if it already exists and create it if it doesnt
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Failed to create or open log file:", err)
	}
	defer file.Close()
	//this is a cool convention, you just set the output of 'log' to the file instead of stdout
	//and you can print to the file like this.
	log.SetOutput(file)

	//write the log message to the log file, remember it has already been parsed for timestamp and user
	log.Println(message)
}
