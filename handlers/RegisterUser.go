package forum

import (
	"database/sql" // Package for SQL database operations
	"net/http"     // Package for HTTP client and server implementations
	"strings"      // Package for string manipulation

	"golang.org/x/crypto/bcrypt" // Package for password hashing
)

// Handler for registering a new user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1000000) // Parse the form data with a maximum memory of 1MB
	if err != nil {                      // If there is an error parsing the form data
		returnJson(w, "{\"error\": \"Failed to parse form data\"}", http.StatusBadRequest) // Return error
		return                                                                             // Exit function
	}
	username := r.FormValue("username")    // Get the username from the form data
	email := r.FormValue("email")          // Get the email from the form data
	password := r.FormValue("password")    // Get the password from the form data
	age := r.FormValue("age")              // Get the age from the form data
	gender := r.FormValue("gender")        // Get the gender from the form data
	firstName := r.FormValue("first_name") // Get the first name from the form data
	lastName := r.FormValue("last_name")   // Get the last name from the form data

	// Check if any required field is missing
	if username == "" || email == "" || password == "" || age == "" || gender == "" || firstName == "" || lastName == "" {
		returnJson(w, "{\"error\": \"Please provide all the details\"}", http.StatusBadRequest) // Return error
		return                                                                                  // Exit function
	} else if !checkEmail(email) { // Check if the email is valid
		returnJson(w, "{\"error\": \"Email not valid\"}", http.StatusBadRequest) // Return error
		return                                                                   // Exit function
	} else if len(strings.Trim(username, " ")) == 0 || len(strings.Trim(email, " ")) == 0 || len(strings.Trim(password, " ")) == 0 || len(strings.Trim(age, " ")) == 0 ||
		len(strings.Trim(gender, " ")) == 0 || len(strings.Trim(firstName, " ")) == 0 || len(strings.Trim(lastName, " ")) == 0 { // Check if any field contains only spaces
		returnJson(w, "{\"error\": \"Field can't contain only spaces\"}", http.StatusBadRequest) // Return error
		return                                                                                   // Exit function
	} else if len(strings.Trim(password, " ")) < 8 { // Check if the password is at least 8 characters long
		returnJson(w, "{\"error\": \"Password must be at least 8 characters long\"}", http.StatusBadRequest) // Return error
		return                                                                                               // Exit function
	}

	// Check if the email or username is already registered
	var existingEmail string
	var existingUsername string
	err = DB.QueryRow("SELECT email, username FROM users WHERE email = ? OR username = ?", email, username).Scan(&existingEmail, &existingUsername) // Execute query
	switch {
	case err == sql.ErrNoRows: // If no rows are returned, proceed with registration
	case err != nil: // If there is an error executing the query
		returnJson(w, "{\"error\": \"Failed to check existing email\"}", http.StatusInternalServerError) // Return error
		return                                                                                           // Exit function
	default: // If the email or username is already registered
		returnJson(w, "{\"error\": \"Email or username is already registered\"}", http.StatusBadRequest) // Return error
		return                                                                                           // Exit function
	}

	if existingEmail != "" || existingUsername != "" { // If the email or username is already registered
		returnJson(w, "{\"error\": \"Email or username is already registered\"}", http.StatusBadRequest) // Return error
		return                                                                                           // Exit function
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // Generate hashed password
	if err != nil {                                                                          // If there is an error hashing the password
		returnJson(w, "{\"error\": \"Internal server error\"}", http.StatusInternalServerError) // Return error
		return                                                                                  // Exit function
	}

	// Insert the new user into the database
	_, err = DB.Exec("INSERT INTO Users (username, email, password, age, gender, first_name, last_name) VALUES (?, ?, ?, ?, ?, ?, ?)", username, email, string(hashedPassword), age, gender, firstName, lastName) // Execute query
	if err != nil {                                                                                                                                                                                               // If there is an error inserting the user
		returnJson(w, "{\"error\": \"Failed to register user\"}", http.StatusInternalServerError) // Return error
		return                                                                                    // Exit function
	}

	// Log the user in and create a session
	sessionID, _, err := login(username, password, w) // Call the login function
	if err != nil {                                   // If there is an error logging in
		returnJson(w, "{\"error\": \""+err.Error()+"\"}", http.StatusInternalServerError) // Return error
		return                                                                            // Exit function
	}

	returnJson(w, "{\"status\":\"ok\", \"session\":\""+sessionID+"\"}", http.StatusOK) // Return success response
}

// List of valid email domain endings
var validEnds = []string{".com", ".org", ".co.bh"}

// checkEmail validates the email format
func checkEmail(email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") { // Check if email contains "@" and "."
		return false // Return false if not valid
	}

	for _, end := range validEnds { // Iterate over valid email endings
		if strings.HasSuffix(email, end) { // Check if email ends with a valid ending
			return true // Return true if valid
		}
	}
	return false // Return false if not valid
}
