package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type User struct {
	Username string
	Password string
}

type Client struct {
	conn   net.Conn
	writer *bufio.Writer
	reader *bufio.Reader
}

type ChatRoom struct {
	Name       string
	Messages   chan string
	Join       chan *Client
	Leave      chan *Client
	clients    map[*Client]bool
	clientsMux sync.RWMutex
}

var (
	activeChatRooms = make(map[string]ChatRoom)
	activeUsers     = make(map[string]User)
	clients         = make(map[*Client]bool)
	mu              sync.Mutex
)

func main() {
	defaultUser := createUser("default", "default")
	fmt.Print("Created default user: ", defaultUser)
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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	//defer conn.Close()

	client := Client{
		conn:   conn,
		writer: bufio.NewWriter(conn),
		reader: bufio.NewReader(conn),
	}

	fmt.Fprintf(conn, "You have selected the login option. Calling the login handler\n")
	//we only need the connection because the handler will ask for the login string
	go handleLogin(client)
	//let that goroutine run its course
	// Handle the received message

}

func handleLogin(client Client) {
	for {
		fmt.Fprintf(client.conn, "This is the login handler. should not be abl to enter chatrooms until the login is handled\n")

		// Process the message or perform any desired logic here
		fmt.Fprint(client.writer, "Please login or create an account (login/create): ")
		client.writer.Flush()

		command, err := client.reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			return
		}
		command = strings.TrimSpace(command)

		fmt.Print("Here is the command you selected: ", command)
		switch command {
		case "login":
			handleLoginCommand(client)
			return
		case "create":
			handleCreateUserCommand(client)
			return
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

// I forgot how to use mutexes, what caught my eye is the mutex has nothing to do with
// the actual objects it protects, the code is just written to only access the variables after it holds
// the mutex
// notice a User and a Client are NOT the same
func handleLoginCommand(client Client) {
	for {
		fmt.Fprintf(client.conn, "This is the login COMMAND handler. should not be abl to enter chatrooms until the login is handled\n")

		fmt.Fprint(client.writer, "Username: ")
		client.writer.Flush()

		username, err := client.reader.ReadString('\n')
		username = strings.TrimSpace(username)

		if err != nil {
			log.Println("Error reading username:", err)
			return
		}
		fmt.Fprintf(client.conn, "You have entered a user with username %s\n", username)

		fmt.Fprint(client.writer, "Password: ")
		client.writer.Flush()

		password, err := client.reader.ReadString('\n')
		password = strings.TrimSpace(password)

		if err != nil {
			log.Println("Error reading password:", err)
			return
		}
		if !authenticateUser(username, password) {
			fmt.Fprintln(client.writer, "Invalid username or password. Please try again.")
			client.writer.Flush()
			continue
		}

		fmt.Fprintf(client.writer, "Welcome, %s! You are now logged in.\n", username)
		client.writer.Flush()
		return
	}
}

func authenticateUser(username, password string) bool {
	user, exists := activeUsers[username]
	return exists && user.Password == password
}

func handleCreateUserCommand(client Client) {

}

func createUser(username, password string) User {
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
