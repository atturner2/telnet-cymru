package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Client struct {
	conn      net.Conn
	writer    *bufio.Writer
	reader    *bufio.Reader
	room      ChatRoom
	user      User
	LoggedOut bool
}

var (
	clients = make(map[*Client]bool)
	mu      sync.Mutex
)

func NewClient(conn *net.Conn) Client {
	//room and user have not been set yet
	return Client{
		conn:      *conn,
		writer:    bufio.NewWriter(*conn),
		reader:    bufio.NewReader(*conn),
		LoggedOut: false,
	}
}

func handleConnection(conn *net.Conn) {
	client := NewClient(conn)

	fmt.Fprintf(*conn, "You have selected the login option. Calling the login handler\n")
	//we only need the connection because the handler will ask for the login string
	for {
		if handleLogin(&client) {
			for {
				fmt.Println("\ncalling handlemainmenu with ", client.user.Username, "\n")
				handleMainMenu(&client)
				fmt.Println("after main menu execution")
				fmt.Println("Here is the client LoggedOut:", client.LoggedOut)
				if client.LoggedOut {
					fmt.Println("Client logged out")
					break // Exit the loop if the client chooses to log out
				}
				fmt.Println("Ran once")
			}
		}
	}
}

func handleLoginCommand(client *Client) bool {
	for {
		fmt.Fprint(client.writer, "Enter username and password as prompted or 'exit' to go back to the entry menu: ")

		fmt.Fprint(client.writer, "Username: ")
		client.writer.Flush()

		username, err := client.reader.ReadString('\n')
		username = strings.TrimSpace(username)
		if username == "exit" {
			break
		}

		if err != nil {
			log.Println("Error reading username:", err)
			return false
		}
		fmt.Fprintf(client.conn, "You have entered a user with username %s\n", username)
		fmt.Fprint(client.writer, "Password: ")
		client.writer.Flush()

		password, err := client.reader.ReadString('\n')
		password = strings.TrimSpace(password)

		if err != nil {
			log.Println("Error reading password:", err)
			return false
		}
		if !authenticateUser(username, password) {
			fmt.Fprintln(client.writer, "Invalid username or password. Please try again.")
			client.writer.Flush()
			continue
		}
		client.user.Username = username
		client.user.Password = password
		fmt.Fprintf(client.writer, "Welcome, %s! You are now logged in.\n", username)
		client.writer.Flush()
		return true
	}
	return false
}

func handleLogin(client *Client) bool {
	for {
		// Process the message or perform any desired logic here
		fmt.Fprint(client.writer, "Please login or create an account (login/create): ")
		client.writer.Flush()

		command, err := client.reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			return false
		}
		command = strings.TrimSpace(command)

		fmt.Print("Here is the command you selected: ", command)
		switch command {
		case "login":
			if handleLoginCommand(client) {
				return true
			}

		case "create":
			handleCreateUserCommand(*client)
			return false
		default:
			fmt.Fprintf(client.conn, "You have entered an invalid command: %s\n", command)

			fmt.Fprintln(client.writer, "Invalid command. Please try again, here is what you entered: ")
			client.writer.Flush()
		}

		command = strings.TrimSpace(command)

		response := fmt.Sprintf("You have selected the login message, this is the response object\n")
		client.conn.Write([]byte(response))
	}
}
func (c *Client) Logout() {
	if c.room.Name != "" {
		c.room.LeaveChatroom <- *c
		c.room.Name = ""
	}
	c.LoggedOut = true
	fmt.Println("Should be logging out")

	return
}

func (c Client) send(message string) {
	c.writer.WriteString(message + "\n")
	c.writer.Flush()
}

func (c Client) Receive() (string, error) {
	message, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(message), nil
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
