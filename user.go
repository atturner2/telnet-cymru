package main

import "fmt"

type User struct {
	Username string
	Password string
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
