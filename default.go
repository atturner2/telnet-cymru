package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type DefaultConfig struct {
	defaultFiles        defaultUserConfig `yaml:"filepaths"`
	defaultServerConfig ServerConfig      `yaml:"server"`
}

type ServerConfig struct {
	connection string `json:"type"`
	address    string `json:"address"`
}
type TestConfig struct {
	testString string `json:"test"`
}

// remember this isn't storing actual chatrooms and users ,
// just the filepaths to the config files
type defaultUserConfig struct {
	defaultUsers     string `json:"users"`
	defaultChatrooms string `json:"chatrooms"`
}

// default loaders load all of the default users and chatrooms from the .csv files
// its easier to call it with the path so we only access the yaml file once from main
func loadDefaultUsers() {
	// Open the CSV file
	file, err := os.Open("config/users.csv")
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

// same thing here, easier to only access the yaml file once
func loadDefaultChatrooms(path string) {
	file, err := os.Open("config/chatrooms.csv")
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

func getDefaultsFilepath() (string, string) {
	//filePath := "config/config.json" // Replace with the actual file path
	fileContent, err := os.Open("config/config.json")
	if err != nil {
		log.Fatal(err)
		return "fuck", "fuck"
	}
	defer fileContent.Close()
	// Read the YAML file
	jsonFile, err := ioutil.ReadAll(fileContent)
	if err != nil {
		log.Fatalf("Failed to read json file: %v", err)
	}
	fmt.Println("Hereis my json file: ", jsonFile)

	// Parse YAML into the configuration struct
	var config DefaultConfig
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	var testConfig TestConfig
	err = json.Unmarshal(jsonFile, &testConfig)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Access the configuration values
	fmt.Println("user path test:", testConfig.testString)

	fmt.Println("user path:", config)
	fmt.Println("user path:", config.defaultServerConfig)
	fmt.Println("user path:", config.defaultFiles.defaultUsers)
	fmt.Println("user path:", config.defaultFiles.defaultChatrooms)

	return "config.Users", "config.Chatrooms"

}
