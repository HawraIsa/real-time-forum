package forum

import (
	"encoding/json" // Package for encoding and decoding JSON
	// Package for formatted I/O
	"net/http" // Package for HTTP client and server implementations
	"slices"   // Package for working with slices
	"strconv"  // Package for converting strings to other types
	"time"     // Package for time-related functions
)

// GetPostsHandler handles retrieving posts and returns them as JSON
func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	userName := r.URL.Query().Get("username") // Get the username from query parameters

	if userName == "" { // If user is not logged in, return an error
		returnJson(w, "{\"error\":\"username not provided\"}", http.StatusUnauthorized) // Return error
		return                                                                          // Exit function
	}

	var userID = 0                                                                            // Variable to hold the user ID
	err := DB.QueryRow("SELECT userid FROM users WHERE username = ?", userName).Scan(&userID) // Query to get the user ID from the database
	if err != nil {                                                                           // If there is an error executing the query
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError) // Return error
		return                                                                    // Exit function
	}

	// Retrieve posts from the database
	posts, err := getAllPosts(userID) // Call function to get all posts
	if err != nil {                   // If there is an error retrieving posts
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError) // Return error
		return                                                                    // Exit function
	}

	// Set response header to application/json
	w.Header().Set("Content-Type", "application/json") // Set response content type to JSON

	// Marshal posts into JSON
	if err := json.NewEncoder(w).Encode(posts); err != nil { // Encode posts slice to JSON and write to response
		http.Error(w, "Failed to encode posts to JSON", http.StatusInternalServerError) // Return error
		return                                                                          // Exit function
	}
}

// getAllPosts retrieves all posts from the database
func getAllPosts(userID int) ([]Post, error) {
	var posts []Post // Slice to hold posts

	// Query the database to retrieve all posts
	rows, err := DB.Query("SELECT * FROM Posts ORDER BY timeposted DESC") // Execute query to get all posts ordered by time posted
	if err != nil {                                                       // If there is an error executing the query
		return nil, err // Return error
	}
	defer rows.Close() // Ensure rows are closed when function exits

	// Iterate over the rows and scan into Post structs
	for rows.Next() { // Loop through each row in the result set
		var post Post                                                                                 // Variable to hold a post
		var timeString string                                                                         // Variable to hold the time as a string
		err := rows.Scan(&post.PostID, &post.PostTitle, &post.PostContent, &timeString, &post.UserID) // Scan row into post struct
		if err != nil {                                                                               // If there is an error scanning the row
			return nil, err // Return error
		}
		unix, err := strconv.ParseInt(timeString, 10, 64) // Parse the timestamp into a number
		if err != nil {                                   // If there is an error parsing the timestamp
			return nil, err // Return error
		}
		post.TimePosted = time.Unix(unix, 0)                                             // Convert the timestamp to a time.Time object
		var usersLiked []int                                                             // Slice to hold user IDs who liked the post
		var usersDisliked []int                                                          // Slice to hold user IDs who disliked the post
		post.Likes, post.Dislikes, usersLiked, usersDisliked = getPostLikes(post.PostID) // Get the likes and dislikes for the post
		post.Categories, err = getPostCategories(post.PostID)                            // Get the categories for the post
		if err != nil {                                                                  // If there is an error getting the categories
			return nil, err // Return error
		}
		post.LikedByUser = slices.Contains(usersLiked, userID)             // Check if the user liked the post
		post.DislikedByUser = slices.Contains(usersDisliked, userID)       // Check if the user disliked the post
		post.Comments, err = GetCommentsForPost(strconv.Itoa(post.PostID)) // Get the comments for the post
		if err != nil {                                                    // If there is an error getting the comments
			return nil, err // Return error
		}
		post.CommentsCount = len(post.Comments) // Set the number of comments

		post.UserName, err = getUserName(post.UserID) // Get the username of the post author
		if err != nil {                               // If there is an error getting the username
			return nil, err // Return error
		}
		posts = append(posts, post) // Add the post to the slice
	}

	// Check for any errors encountered while iterating over rows
	if err := rows.Err(); err != nil { // If there is an error iterating over the rows
		return nil, err // Return error
	}

	return posts, nil // Return the slice of posts
}
