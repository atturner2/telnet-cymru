package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func loadDefaultUsers() {
	// Open the CSV file
	file, err := os.Open("users.csv")
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

func loadDefaultChatrooms() {
	file, err := os.Open("data.csv")
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
			go chatroom.start()
		}

	}
}
