package forum

import "database/sql" // Package for SQL database operations

// getUserName retrieves the username from the user ID
func getUserName(userID int) (string, error) {
	var userName string // Variable to hold the username
	// Query to get the username from the database using the user ID
	err := DB.QueryRow("SELECT username FROM Users where userId = ?", userID).Scan(&userName)
	if err == sql.ErrNoRows { // If no rows are returned, the user does not exist
		return "", nil // Return empty string and no error
	} else if err != nil { // If there is another error
		return "", err // Return the error
	}

	return userName, nil // Return the username
}
