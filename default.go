package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	ConnectionType string `json:"connectionType"`
	Port           string `json:"port"`
	Users          string `json:"users"`
	Chatrooms      string `json:"chatrooms"`
	Test           string `json:"test"`
}

// default loaders load all of the default users and chatrooms from the .csv files
// its easier to call it with the path so we only access the json file once from main
func loadDefaultUsers(path string) {
	// Open the CSV file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all rows from the CSV file
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// Process each row
	for _, row := range rows {
		// Ensure the row contains two fields (username and password)
		if len(row) >= 2 {
			username := row[0]
			password := row[1]
			createUser(username, password)
		}
	}
}

// same thing here, easier to only access the json file once
func loadDefaultChatrooms(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all rows from the CSV file, each row is a chatroom
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}
	//
	for _, row := range rows {
		if len(row) == 1 {
			chatroomName := row[0]
			chatroom := createChatRoom(chatroomName)
			//we start the chatroom/create it's goroutine here
			go chatroom.start()
		}

	}
}

func getDefaults() Config {
	fileContent, err := os.Open("config/config.json")

	if err != nil {
		log.Fatal(err)
	}

	defer fileContent.Close()

	byteResult, _ := ioutil.ReadAll(fileContent)

	var config Config
	err = json.Unmarshal(byteResult, &config)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Access the configuration values
	return config

}
