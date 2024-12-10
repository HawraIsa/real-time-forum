package forum

import (
	"encoding/json" // Package for encoding and decoding JSON
	"log"           // Package for logging messages
	"net/http"      // Package for HTTP client and server implementations
	"time"          // Package for time-related functions

	"github.com/gorilla/websocket" // Package for WebSocket support
)

// Upgrader to handle WebSocket connections with specified buffer sizes and origin check
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // Read buffer size
	WriteBufferSize: 1024, // Write buffer size
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

// Maps and channels to manage WebSocket connections and user statuses

var userConnections = make(map[string]map[*websocket.Conn]bool) // Map to track active WebSocket connections for each user
var broadcast = make(chan PrivateMessages)                      // Channel for broadcasting messages
var userBroadcast = make(chan string)                           // Channel for broadcasting new users
var userLogoutBroadcast = make(chan string)                     // Channel for broadcasting user logouts

var onlineUsers = make(map[string]bool) // Map of usernames to online status

// WebSocketHandler handles WebSocket requests from clients
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username") // Get username from query parameters
	if username == "" {                       // If username is not provided
		http.Error(w, "Username is required", http.StatusBadRequest) // Return error
		log.Println("WebSocketHandler: Username is required")        // Log error
		return                                                       // Exit function
	}

	log.Println("WebSocketHandler: Attempting to upgrade connection for user:", username)
	conn, err := upgrader.Upgrade(w, r, nil) // Upgrade HTTP connection to WebSocket
	if err != nil {                          // If there is an error during upgrade
		log.Println("WebSocketHandler: Error upgrading connection:", err) // Log error
		return                                                            // Exit function
	}

	defer conn.Close() // Ensure connection is closed when function exits
	log.Println("WebSocketHandler: Connection upgraded for user:", username)

	// Add connection to user's set of connections
	if userConnections[username] == nil {
		userConnections[username] = make(map[*websocket.Conn]bool)
	}

	userConnections[username][conn] = true
	userBroadcast <- username
	log.Println("WebSocketHandler: User added to online users:", username)
	log.Printf("WebSocketHandler: Current online users: %v", userConnections)

	defer func() {
		delete(userConnections[username], conn)
		if len(userConnections[username]) == 0 {
			delete(userConnections, username)
			userLogoutBroadcast <- username
			log.Println("WebSocketHandler: User disconnected and removed from online users:", username)
			log.Printf("WebSocketHandler: Current online users: %v", userConnections)
		}
	}()

	// Message reading loop
	for {
		var msg PrivateMessages
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocketHandler: Error reading JSON: %v", err)
			delete(userConnections[username], conn)
			if len(userConnections[username]) == 0 {
				delete(userConnections, username)
				log.Printf("WebSocketHandler: Current online users: %v", userConnections)
			}
			break
		}
		msg.TimeSent = time.Now()
		log.Println("WebSocketHandler: Message received:", msg)

		// Save to DB first, then broadcast if successful
		if err := saveMessageToDB(msg); err != nil {
			log.Printf("WebSocketHandler: Failed to save message to DB: %v", err)
			continue // Skip broadcasting if DB save fails
		}

		broadcast <- msg
	}
}

// handleMessages listens for broadcasts and sends them to all clients
func handleMessages() {
	for {
		msg := <-broadcast
		log.Println("handleMessages: Broadcasting message:", msg)
		// Message is already saved to DB, just broadcast
		for _, connections := range userConnections {
			for client := range connections {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("handleMessages: Error writing JSON: %v", err)
					client.Close()
					delete(connections, client)
				}
			}
		}
	}
}

// handleNewUser listens for new user notifications and broadcasts them to all clients
func handleNewUser() {
	for {
		user := <-userBroadcast
		log.Println("handleNewUser: Broadcasting new user:", user)
		for _, connections := range userConnections {
			for client := range connections {
				err := client.WriteJSON(map[string]string{"newUser": user})
				if err != nil {
					log.Printf("handleNewUser: Error writing JSON: %v", err)
					client.Close()
					delete(connections, client)
				}
			}
		}
	}
}

// handleUserLogout listens for user logout notifications and broadcasts them to all clients
func handleUserLogout() {
	for {
		user := <-userLogoutBroadcast
		log.Println("handleUserLogout: Broadcasting user logout:", user)
		for _, connections := range userConnections {
			for client := range connections {
				err := client.WriteJSON(map[string]string{"removeUser": user})
				if err != nil {
					log.Printf("handleUserLogout: Error writing JSON: %v", err)
					client.Close()
					delete(connections, client)
				}
			}
		}
	}
}

// saveMessageToDB saves a message to the database with error
func saveMessageToDB(msg PrivateMessages) error {
	log.Println("saveMessageToDB: Saving message to database:", msg)
	_, err := DB.Exec("INSERT INTO PrivateMessages (senderUsername, receiverUsername, messageText, timeSent) VALUES (?, ?, ?, ?)",
		msg.SenderUsername, msg.ReceiverUsername, msg.MessageText, msg.TimeSent)
	if err != nil {
		log.Printf("saveMessageToDB: Error saving message to database: %v", err)
		return err
	}
	log.Println("saveMessageToDB: Message successfully saved to database")
	return nil
}

// GetOnlineUsersHandler handles HTTP requests to fetch the list of online users
func GetOnlineUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetOnlineUsersHandler: Fetching online users")
	users := []string{} // Slice to hold online users
	for user := range userConnections {
		users = append(users, user)
	}
	log.Printf("GetOnlineUsersHandler: Online users list: %v", users)
	w.Header().Set("Content-Type", "application/json") // Set response content type to JSON
	

	
	err := json.NewEncoder(w).Encode(users)            // Encode users slice to JSON and write to response
	if err != nil {                                    // If there is an error encoding JSON
		log.Printf("GetOnlineUsersHandler: Error encoding JSON: %v", err) // Log error
	} else {
		log.Println("GetOnlineUsersHandler: Online users sent:", users)
	}
}


// init initializes the goroutines for handling messages, new users, and user logouts
func init() {
	go handleMessages()   // Start goroutine to handle messages
	go handleNewUser()    // Start goroutine to handle new users
	go handleUserLogout() // Start goroutine to handle user logouts
}
