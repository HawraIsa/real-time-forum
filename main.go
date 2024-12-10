package main

import (
	"fmt"                  // Package for formatted I/O
	forum "forum/handlers" // Importing the forum handlers package
	"net/http"             // Package for HTTP client and server implementations
)

func main() {
	db, err := getDatabase() // Initialize the database
	if err != nil {          // If there is an error initializing the database
		fmt.Println(err) // Print the error
		return           // Exit the function
	}
	forum.DB = db // Assign the database to the forum package's DB variable

	// Define HTTP routes and their corresponding handler functions

	http.HandleFunc("/getposts", forum.GetPostsHandler) // Route for getting posts
	http.HandleFunc("/register", forum.RegisterUser)    // Route for user registration
	http.HandleFunc("/login", forum.UserLogin)          // Route for user login

	http.HandleFunc("/authSession", forum.AuthLogin)    // Route for session authentication
	http.HandleFunc("/logout", forum.UserLogoutHandler) // Route for user logout
	http.HandleFunc("/addpost", forum.CreatePost)       // Route for creating a post
	http.HandleFunc("/addComment", forum.CreateComment) // Route for adding a comment

	// Routes that return JSON responses

	http.HandleFunc("/getcategories", forum.GetCategoriesHandler) // Route for getting categories
	http.HandleFunc("/likePost/", forum.HandleLike)               // Route for liking a post
	http.HandleFunc("/likeComment/", forum.HandleCommentLike)     // Route for liking a comment

	// WebSocket route
	http.HandleFunc("/ws", forum.WebSocketHandler)                     // Route for WebSocket connections
	http.HandleFunc("/getPrivateMessages", forum.GetMessagesHandler)   // Route for getting private messages
	http.HandleFunc("/getOnlineUsers", forum.GetOnlineUsersHandler)    // Route for getting online users
	http.HandleFunc("/getMessageUsers", forum.GetMessagesUsersHandler) // Route for getting message users
	http.HandleFunc("/getAllUserMessages", forum.GetAllUserMessagesHandler)

	http.HandleFunc("/", forum.Index) // Route for the index page

	// Serve static files
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("template")))) // Serve template files
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))       // Serve asset files

	// Start the HTTP server
	fmt.Println("Server started on http://localhost:8080") // Print server start message
	http.ListenAndServe(":8080", nil)                      // Start the server on port 8081
}
