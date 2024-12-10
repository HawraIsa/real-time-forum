package forum

import (
	"database/sql" // Package for SQL database operations
	"fmt"          // Package for formatted I/O
	"net/http"     // Package for HTTP client and server implementations
	"time"         // Package for time-related functions

	"github.com/google/uuid"     // Package for generating UUIDs
	"golang.org/x/crypto/bcrypt" // Package for password hashing
)

// login handles user login and session creation
func login(username string, password string, w http.ResponseWriter) (string, string, error) {
	// Retrieve user from the database based on the provided identifier (email or username)
	var storedPassword string // Variable to hold the stored password hash
	var userID int
	var name string
	// Variable to hold the user ID
	err := DB.QueryRow("SELECT password, username, userID FROM users WHERE username = ? or email = ?", username, username).Scan(&storedPassword, &name, &userID) // Query to get the stored password and user ID
	switch {
	case err == sql.ErrNoRows: // If no rows are returned, the username/password is incorrect
		return "", "", fmt.Errorf("username/password not correct") // Return error
	case err != nil: // If there is another error
		return "", "", err // Return the error
	}

	// Compare the stored password hash with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)) // Compare the hashed password with the provided password
	if err != nil {                                                               // If the password does not match
		return "", "", err // Return the error
	}

	// Delete old sessions for the user so that only the most recent one is kept
	_, err = DB.Exec("DELETE FROM sessions where userID = ?", userID) // Execute query to delete old sessions
	if err != nil {                                                   // If there is an error executing the query
		return "", "", err // Return the error
	}

	// Create a new session
	sessionID := uuid.New().String()                                                                                                                                                                 // Generate a new UUID for the session ID
	_, err = DB.Exec("INSERT INTO sessions (sessionData, userID, createdAt, lastAccessedAt, expiry) VALUES (?, ?, ?, ?, ?)", sessionID, userID, time.Now(), time.Now(), time.Now().AddDate(0, 0, 1)) // Execute query to insert the new session
	if err != nil {                                                                                                                                                                                  // If there is an error inserting the session
		return "", "", err // Return the error
	}

	return sessionID, name, nil // Return the session ID
}
