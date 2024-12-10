package forum

import (
	"net/http" // Package for HTTP client and server implementations
	"strings"  // Package for string manipulation
	"time"     // Package for time-related functions
)

// AddPostResponse holds the data for the addpost HTML response
type AddPostResponse struct {
	Message    string     // Message to be displayed
	Categories []Category // List of categories
}

// Handler for creating a new post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	userName := r.URL.Query().Get("username") // Get the username from query parameters

	if userName == "" { // If username is not provided
		returnJson(w, "{\"error\":\"username not provided\"}", http.StatusUnauthorized) // Return error
		return                                                                          // Exit function
	}

	var userID = 0                                                                            // Variable to hold the user ID
	err := DB.QueryRow("SELECT userid FROM users WHERE username = ?", userName).Scan(&userID) // Query to get the user ID from the database
	if err != nil {                                                                           // If there is an error executing the query
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError) // Return error
		return                                                                    // Exit function
	}

	// Parse form data from request
	err = r.ParseMultipartForm(10000000) // Parse form data from request with a maximum memory of 10MB
	if err != nil {                      // If there is an error parsing the form data
		returnJson(w, "{\"error\":\"Failed to parse form data\"}", http.StatusBadRequest) // Return error
		return                                                                            // Exit function
	}

	// Extract post data from form
	postTitle := r.FormValue("title")                        // Get the post title from the form data
	postContent := r.FormValue("content")                    // Get the post content from the form data
	postCategories := strings.Split(r.FormValue("cat"), ",") // Get the post categories from the form data and split by comma

	// Validate input data
	if postTitle == "" || postContent == "" || len(strings.Trim(postTitle, " ")) == 0 || len(strings.Trim(postContent, " ")) == 0 {
		returnJson(w, "{\"error\":\"Please provide post title and post content\"}", http.StatusBadRequest) // Return error
		return                                                                                             // Exit function
	}

	// Insert the new post into the database
	result, err := DB.Exec("INSERT INTO posts (postTitle, postContent, timePosted, userID) VALUES (?, ?, ?, ?)", postTitle, postContent, time.Now().Unix(), userID) // Execute query to insert the post
	if err != nil {                                                                                                                                                 // If there is an error inserting the post
		returnJson(w, "{\"error\":\"Failed to create post:"+err.Error()+"\"}", http.StatusInternalServerError) // Return error
		return                                                                                                 // Exit function
	}
	postId, _ := result.LastInsertId() // Get the post ID that was just added

	for _, cat := range postCategories { // Loop over the categories chosen by the user and add them to categoriesposts table
		_, err = DB.Exec("INSERT INTO CategoriesPosts (postid, catid) VALUES (?, ?)", postId, cat) // Execute query to insert the category
		if err != nil {                                                                            // If there is an error inserting the category
			returnJson(w, "{\"error\":\"Failed to add categories:"+err.Error()+"\"}", http.StatusInternalServerError) // Return error
			return                                                                                                    // Exit function
		}
	}

	returnJson(w, "{\"status\":\"ok\"}", http.StatusOK) // Return success response
}
