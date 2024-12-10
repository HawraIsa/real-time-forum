package forum

import (
	"database/sql" // Package for SQL database operations
	"net/http"     // Package for HTTP client and server implementations
)

// Handler for user login
func UserLogin(w http.ResponseWriter, r *http.Request) {
	// Parse form data from request
	err := r.ParseMultipartForm(1000000) // Parse the form data with a maximum memory of 1MB
	if err != nil {                      // If there is an error parsing the form data
		returnJson(w, "{\"error\": \"Failed to parse form data\"}", http.StatusUnauthorized) // Return error
		return                                                                               // Exit function
	}

	// Extract login data from form
	identifier := r.FormValue("username")                      // This can be either email or username
	password := r.FormValue("password")                        // Get the password from the form data
	sessionID, username, err := login(identifier, password, w) // Attempt to log in the user
	if err != nil {                                            // If there is an error during login
		returnJson(w, err.Error(), http.StatusBadRequest) // Return error
		return                                            // Exit function
	}

	// Add user to onlineUsers map only if login is successful
	if identifier != "" { // If identifier is provided
		onlineUsers[username] = true // Mark user as online
	}

	// Respond with success message or token
	returnJson(w, "{\"status\":\"ok\", \"session\":\""+sessionID+"\", \"username\":\""+username+"\"}", http.StatusOK) // Return success response
}

// Handler for authenticating a user session
func AuthLogin(w http.ResponseWriter, r *http.Request) {
	// Extract session data from query parameters
	session := r.URL.Query().Get("session")                                                                                                      // Get session data from query parameters
	var userName string                                                                                                                          // Variable to hold the username
	err := DB.QueryRow("select username from users where userid = (SELECT userid FROM sessions WHERE sessionData = ?)", session).Scan(&userName) // Get the username from the database using the session data
	if err == sql.ErrNoRows {                                                                                                                    // If no rows are returned
		returnJson(w, "{\"error\": \"User data not found\"}", http.StatusUnauthorized) // Return error
		return                                                                         // Exit function
	}

	// Respond with success message and username
	returnJson(w, "{\"status\":\"ok\", \"username\":\""+userName+"\"}", http.StatusOK) // Return success response
}
