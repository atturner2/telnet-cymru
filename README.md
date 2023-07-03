# telnet-cymru
This project was written as part of my application for employment at Cymru. 
# Description
This project is a very simple chatroom. Users can log in to an account configured in one of the default files or create their own account. They can log out and log in as someone else. Once they are logged in, 
they can select a chat room to join or create one of their own. There are default chatrooms created in the chatrooms.csv file and the users.csv file. The application runs on the configurations from the congig.json
file. The only hardcoded file location in the project is that config file, but if you wanted to you could reconfigure it with env vars or something. 
# Installation
Pull down or unzip the project. Run "Go Build" to generate the executable and then run ./telnet-cymru to start the server. To connect a client to the chatroom, first run 'brew install telnet' and then run 
'telnet localhost PORT_NUMBER' where PORT_NUMBER is the "port" field in config/config.json, currently is 23. Each server you start is a single client connection that gets it's own goroutine. They all connect on the same 
port but they each get a goroutine created for them.
# Usage
The server will send you prompts to log in or create a user manually. You can also create users by adding them to the config/users.csv file. The format must be username,password with no trailing spaces. Note that creating a user does NOT log you in as that user and currently there is no support for forgotten passwords. Once you are logged in you can either join a chatroom from a list of chatrooms or create one of your own. 
The precompiled list of chatrooms is located in the chatrooms.csv file and you are free to edit that file. 
# Features
* Users can log in and create accounts. Users can then join and create their own chatrooms for themselves and all of the other users.
* Users in a chatroom will recieve all of the messages in the chatroom but will not see previous messages from before they joined.
* All messages will be logged in the logs/name_of_chatroom.csv file. Note the logs are attached to the chatroom, NOT the user.
* Users can log in and out of their accounts and chatrooms as much as they want, a single connection only has one user at a time but can cycle through as many users as you want.
# How It Works
* Upon startup, the program will call getDefaults which will grab all of the default config. It will then call the default functions to load users and load chatrooms that are in the CSV files.
* The server then listens for connections and each connection will get it's own goroutine. That goroutine will create a Client object for it's connection. Those have a 1:1 relationship and will stay together for the entire execution.
* The login/create client for loop structure can be found in client.go, basically you can either create a user or log in as an existing one, but there is no support beyond that. Note that because of the create
functionality you can have the users.csv and chatrooms.csv files empty and it will work fine, you just have to create a user.
* You can log in and out of a chatroom or an account or multiple chatrooms and multiple accounts as many times as you want. This is why I had to have Users seperate from Clients, though like I said above that relationship
could and should have been better managed.
* Chatrooms have channels for communication among the users in the chatroom. The channel communicates with the user sending the message and then sends the message to all of the users in the chatroom. Remember each user has their own goroutine and each chatroom has its own goroutine. 
# Challenges/limitations/future Improvements (Why this is a bad design)
* In 20/20 hindsight I should have just reached out to Ryan right away and asked if he had Docker set up and just written a set of instructions with this project to run some databases in a Docker container so this could actually have real authentication and not just be running everything in memory. This is honestly a terrible design because there is no database and alot of the challenge of the project came from that more than anything else, all of the users and clients etc are stored in various objects/pointers which made this harder and more fragile than necessary. In reality this all could have been done by storing objects and memory addresses in a database and just using mutexes to regulate access.  This also would have eliminated the need for all the .CSV files because the defaults could have been stored in a database. 
* There is no reason for the chatroom objects to be storing all of the clients, I should have just had it store the names of the clients and look them up a the communal struct of pointers of active clients (Note clients != users)
instead of the way I did it where the chatroom stores all of it's clients. This is just a waste of memory in my opinion. Overall the relationship of connections:clients:users is not managed very well and should be refactored.
* The login functionality doesn't throw a specific error if you type a username that doesnt exist vs. an existing username with a wrong password
* There is a list of active users and each user is storing it's own login status in that list, but the clients are also tracking login status (Look at Logout() and handleLoginCommand functions in client.go, as well as userIsAlreadyLoggedIn in user.go
* This is a hacky patch to the "logging into the same user twice" issue that I basically realized had a bug at the last minute and threw it together. Overall I should have set up all of the user log in and log out before even touching the chatroom functionality.
# Required
* Clients interact with service via CLI telnet
* Support multiple connected clients
* Messages are send to all connected clients -> I assumed if I implemented multiple chatrooms it should not send messages to users outside of the chatroom, which seemed to fit the Optional Features requirement. 
* Messages are written to log file
* Basic config stored in config gile
# Optional Features
* Allow users to join channels (I called them 'chatrooms')
I decided at the beginning to definitely do this feature and get that all working before I added more. Unfortunately it took a little more work than I thought. I would have liked to add more, including the help functionality and ignore but ran myself out of time. help would have been easy to add, ignore would have taken some work/design decisions.
# Third Party Code/documentation resources
* The code in this project is my own but I looked at alot of different examples to figure out how to do it. Some of those examples can be seen below but it is not an exhaustive list. Alot of documentation was looked at.
* 
* https://golangdocs.com/golang-read-json-file
* https://stackoverflow.com/questions/36417199/how-to-broadcast-message-using-channel
* https://go.dev/tour/concurrency/2
* https://www.geeksforgeeks.org/fmt-fprintf-function-in-golang-with-examples/
* https://nqbao.medium.com/writing-a-chat-server-in-go-3b61ccc2a8ed
* https://gobyexample.com/mutexes


