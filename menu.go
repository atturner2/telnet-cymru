package main

import (
	"fmt"
	"log"
	"strings"
)

func handleMainMenu(client *Client) {
	for {
		fmt.Println("At the main menu with user: ", client.user.Username)
		fmt.Fprint(client.writer, "Please select 'join' to join a chatroom, 'create' to create one, or 'logout to log out.' : ")
		client.writer.Flush()

		command, err := client.reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			return
		}

		command = strings.TrimSpace(command)
		switch command {
		case "join":
			fmt.Fprintln(client.conn, "You have selected to join a chat room!")
			fmt.Println("CLIENT JOINING ROOM: ", client.user.Username)
			handleJoinRoom(client)
			fmt.Println("After the join room call")
		case "create":
			fmt.Fprintln(client.conn, "You have selected to create your own chat room!")
			//handleCreateRoom(client)

		case "logout":
			client.Logout()
			client.LoggedOut = true
			fmt.Println("Should be logged out, ", client.LoggedOut)
			return
		default:
			fmt.Fprintln(client.writer, "Invalid command. Please try again.")
			client.writer.Flush()
		}

	}
	return
}
