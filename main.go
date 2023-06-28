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
	activeChatRooms = make(map[string]*ChatRoom)
	activeUsers     = make(map[string]string)
	clients         = make(map[*Client]bool)
	mu              sync.Mutex
)

func main() {
	defaultChatRoom := createChatRoom("default")
	defaultUser := createUser("default", "default")
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
	defer conn.Close()

	fmt.Fprintf(conn, "Welcome to the Telnet server!\n")
	fmt.Fprintf(conn, "Please select an option from the following list: login to an exising account, or create an existing user\n")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		if strings.TrimSpace(message) == "exit" {
			break
		}
		if strings.TrimSpace(message) == "login" {
			fmt.Fprintf(conn, "You have selected the login option. Calling the login handler\n")
			//we only need the connection because the handler will ask for the login string
			go handleLogin(conn)
		}

		// Handle the received message
	}
}

func handleLogin(conn net.Conn) {
	fmt.Printf("You have selected the login option, this is in the handler\n")

	// Process the message or perform any desired logic here

	response := fmt.Sprintf("You have selected the login message, this is the response object\n")
	conn.Write([]byte(response))
}
