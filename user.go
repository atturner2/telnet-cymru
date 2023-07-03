package main

import (
	"fmt"
	"log"
	"strings"
)

type User struct {
	Username string
	Password string
}

func createUser(username, password string) User {
	//we do need a mutex here because we are editing the activeUsers object
	mu.Lock()
	fmt.Println("creating a user with username ", username, " and password: ", password)
	defer mu.Unlock()

	user := User{
		Username: username,
		Password: password,
	}
	activeUsers[username] = user
	fmt.Println("Here are all the active users: ", activeUsers)
	return user
}

func authenticateUser(username, password string) bool {
	//we dont need a mutex here because it's not ever editing the active users, just checking the list
	user, exists := activeUsers[username]
	return exists && user.Password == password
}

func userExists(username string) bool {
	//we also dont need a mutex here because it's not ever editing the active users, just checking the list
	_, exists := activeUsers[username]
	return exists
}

func handleCreateUserCommand(client Client) {
	for {
		fmt.Fprint(client.writer, "Username: ")
		client.writer.Flush()

		username, err := client.reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading username:", err)
			return
		}

		username = strings.TrimSpace(username)
		if username == "" {
			fmt.Fprintln(client.writer, "Username cannot be empty. Please try again.")
			client.writer.Flush()
			continue
		}

		if userExists(username) {
			fmt.Fprintln(client.writer, "Username already exists. Please choose a different username.")
			client.writer.Flush()
			continue
		}

		fmt.Fprint(client.writer, "Password: ")
		client.writer.Flush()

		password, err := client.reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading password:", err)
			return
		}

		password = strings.TrimSpace(password)
		if password == "" {
			fmt.Fprintln(client.writer, "Password cannot be empty. Please try again.")
			client.writer.Flush()
			continue
		}

		createUser(username, password)

		fmt.Fprintf(client.writer, "Account created. Welcome, %s!\n", username)
		client.writer.Flush()

		return

	}

}
