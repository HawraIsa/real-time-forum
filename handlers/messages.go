package forum

import (
	"encoding/json" // Package for encoding and decoding JSON
	"net/http"      // Package for HTTP client and server implementations
	"strconv"       // Package for converting strings to other types
)

// GetMessagesHandler handles HTTP requests to fetch messages between two users
func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	senderUsername := r.URL.Query().Get("senderUsername")     // Get senderUsername from query parameters
	receiverUsername := r.URL.Query().Get("receiverUsername") // Get receiverUsername from query parameters
	page := r.URL.Query().Get("page")                         // Get page number from query parameters
	if senderUsername == "" || receiverUsername == "" {       // If either username is not provided
		returnJson(w, "{\"error\": \"SenderUsername and ReceiverUsername are required\"}", http.StatusBadRequest) // Return error
		return                                                                                                    // Exit function
	}

	pageNumber, err := strconv.Atoi(page) // Convert page number to integer
	if err != nil || pageNumber < 1 {     // If there is an error or page number is less than 1
		pageNumber = 1 // Default to page 1
	}

	limit := 10                        // Number of messages per page
	offset := (pageNumber - 1) * limit // Calculate offset for SQL query

	// Fetch messages from the database
	messages := []PrivateMessages{} // Slice to hold messages
	query := `
		SELECT messageId, senderUsername, receiverUsername, messageText, timeSent
		FROM PrivateMessages
		WHERE (senderUsername = ? AND receiverUsername = ?) OR (senderUsername = ? AND receiverUsername = ?)
		ORDER BY timeSent DESC LIMIT ? OFFSET ?
	` // SQL query to fetch messages
	rows, err := DB.Query(query, senderUsername, receiverUsername, receiverUsername, senderUsername, limit, offset) // Execute query
	if err != nil {                                                                                                 // If there is an error executing query
		returnJson(w, "{\"error\": \""+err.Error()+"\"}", http.StatusInternalServerError) // Return error
		return                                                                            // Exit function
	}
	defer rows.Close() // Ensure rows are closed when function exits

	for rows.Next() { // Iterate over result set
		var msg PrivateMessages                                                                                                        // Variable to hold a message
		if err := rows.Scan(&msg.MessageID, &msg.SenderUsername, &msg.ReceiverUsername, &msg.MessageText, &msg.TimeSent); err != nil { // Scan row into message
			returnJson(w, "{\"error\": \""+err.Error()+"\"}", http.StatusInternalServerError) // Return error
			return                                                                            // Exit function
		}
		messages = append(messages, msg) // Add message to slice
	}

	json.NewEncoder(w).Encode(messages) // Encode messages slice to JSON and write to response
}

// GetMessagesUsersHandler handles HTTP requests to fetch users for messaging
func GetMessagesUsersHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username") // Get username from query parameters
	if username == "" {                       // If username is not provided
		returnJson(w, "{\"error\": \"username is required\"}", http.StatusBadRequest) // Return error
		return                                                                        // Exit function
	}

	// Fetch users from the database
	users := []string{} // Slice to hold users
	query := `
		SELECT username
		FROM users
		WHERE username != ?
	` // SQL query to fetch users
	rows, err := DB.Query(query, username) // Execute query
	if err != nil {                        // If there is an error executing query
		returnJson(w, "{\"error\": \""+err.Error()+"\"}", http.StatusInternalServerError) // Return error
		return                                                                            // Exit function
	}
	defer rows.Close() // Ensure rows are closed when function exits

	for rows.Next() { // Iterate over result set
		var user string                          // Variable to hold a username
		if err := rows.Scan(&user); err != nil { // Scan row into user
			returnJson(w, "{\"error\": \""+err.Error()+"\"}", http.StatusInternalServerError) // Return error
			return                                                                            // Exit function
		}
		users = append(users, user) // Add user to slice
	}

	json.NewEncoder(w).Encode(users) // Encode users slice to JSON and write to response
}

// GetAllUserMessagesHandler handles HTTP requests to fetch all messages for a user
func GetAllUserMessagesHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    username := r.URL.Query().Get("username")
    if username == "" {
        http.Error(w, `{"error": "username is required"}`, http.StatusBadRequest)
        return
    }

    // Fetch all messages where the user is either sender or receiver
    messages := []PrivateMessages{}
    query := `
        SELECT messageId, senderUsername, receiverUsername, messageText, timeSent
        FROM PrivateMessages
        WHERE senderUsername = ? OR receiverUsername = ?
        ORDER BY timeSent DESC
    `
    rows, err := DB.Query(query, username, username)
    if err != nil {
        http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var msg PrivateMessages
        if err := rows.Scan(&msg.MessageID, &msg.SenderUsername, &msg.ReceiverUsername, &msg.MessageText, &msg.TimeSent); err != nil {
            http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
            return
        }
        messages = append(messages, msg)
    }

    if err := json.NewEncoder(w).Encode(messages); err != nil {
        http.Error(w, `{"error": "Failed to encode messages"}`, http.StatusInternalServerError)
        return
    }
}