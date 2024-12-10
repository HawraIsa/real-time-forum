package forum

import (
	"net/http" // Package for HTTP client and server implementations
	"strings"  // Package for string manipulation
)

// Handler for creating a new comment
func CreateComment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10000000) // Parse form data from request with a maximum memory of 10MB
	if err != nil {                       // If there is an error parsing the form data
		returnJson(w, "{\"error\":\"Failed to parse form data\"}", http.StatusBadRequest) // Return error
		return                                                                            // Exit function
	}

	// Extract comment data from the form
	commentContent := r.FormValue("comment") // Get the comment content from the form data
	postID := r.FormValue("id")              // Get the post ID from the form data

	// Check if user is authenticated
	userName := r.URL.Query().Get("username") // Get the username from query parameters

	if userName == "" { // If username is not provided
		returnJson(w, "{\"error\":\"username not provided\"}", http.StatusUnauthorized) // Return error
		return                                                                          // Exit function
	}

	var userID = 0                                                                           // Variable to hold the user ID
	err = DB.QueryRow("SELECT userid FROM users WHERE username = ?", userName).Scan(&userID) // Query to get the user ID from the database
	if err != nil {                                                                          // If there is an error executing the query
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError) // Return error
		return                                                                    // Exit function
	}

	// Validate input data
	if commentContent == "" || postID == "" || len(strings.Trim(commentContent, " ")) == 0 { // Check if comment content or post ID is empty or contains only spaces
		returnJson(w, "{\"error\":\"Please provide comment text\"}", http.StatusBadRequest) // Return error
		return                                                                              // Exit function
	}

	// Insert the new comment into the database
	_, err = DB.Exec("INSERT INTO comments (commentContent, userID, postID) VALUES (?, ?, ?)", commentContent, userID, postID) // Execute query to insert the comment
	if err != nil {                                                                                                            // If there is an error inserting the comment
		returnJson(w, "{\"error\":\"Failed to create comment: "+err.Error()+"\"}", http.StatusInternalServerError) // Return error
		return                                                                                                     // Exit function
	}

	returnJson(w, "{\"status\":\"ok\"}", http.StatusOK) // Return success response
}
