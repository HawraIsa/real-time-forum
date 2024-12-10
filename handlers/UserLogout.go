package forum

import (
	"fmt"      // Package for formatted I/O
	"net/http" // Package for HTTP client and server implementations
)

// UserLogoutHandler handles the logout route
func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := logout(w, r) // Call the logout function
	if err != nil {     // If there is an error during logout
		returnJson(w, "{\"error\": \""+err.Error()+"\"}", http.StatusInternalServerError) // Return error response
	} else {
		returnJson(w, "{\"status\": \"ok\"}", http.StatusOK) // Return success response
	}
}

// logout performs the logout operation
func logout(w http.ResponseWriter, r *http.Request) error {
	userName := r.URL.Query().Get("username") // Get the username from query parameters
	if userName == "" {                       // If username is not provided
		return fmt.Errorf("Username not provided") // Return error
	}
	// Execute SQL query to delete the session for the user
	_, err := DB.Exec("DELETE FROM sessions where userID = (select userid from users where username = ?)", userName)
	if err != nil { // If there is an error executing the query
		return err // Return the error
	}
	return nil // Return nil if no error
}
