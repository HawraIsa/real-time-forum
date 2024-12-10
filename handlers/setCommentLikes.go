package forum

import (
	"database/sql"  // Package for SQL database operations
	"encoding/json" // Package for encoding and decoding JSON
	"net/http"      // Package for HTTP client and server implementations
	"strconv"       // Package for converting strings to other types
)

// HandleCommentLike handles the like or dislike action for a comment
func HandleCommentLike(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("username") // Get the username from query parameters

	if userName == "" { // If user is not logged in, return an error
		returnJson(w, "{\"error\":\"username not provided\"}", http.StatusUnauthorized) // Return error
		return                                                                          // Exit function
	}
	w.Header().Set("Content-Type", "application/json") // Set response content type to JSON
	if userName != "" {
		commentID, _ := strconv.Atoi(r.URL.Query().Get("id")) // Get the comment ID from query parameters and convert to integer
		isDislike := 0                                        // Initialize isDislike to 0 (like)
		if r.URL.Query().Get("isLike") == "false" {           // If isLike is false, set isDislike to 1 (dislike)
			isDislike = 1
		}
		var userID = 0                                                                            // Variable to hold the user ID
		err := DB.QueryRow("SELECT userid FROM users WHERE username = ?", userName).Scan(&userID) // Query to get the user ID from the database
		if err != nil {                                                                           // If there is an error executing the query
			http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError) // Return error
			return                                                                    // Exit function
		}
		likes, dislikes, actionTaken, err := setCommentLikes(commentID, userID, isDislike) // Call function to set comment likes or dislikes
		if err != nil {                                                                    // If there is an error setting comment likes or dislikes
			http.Error(w, err.Error(), http.StatusInternalServerError) // Return error
		}
		data, _ := json.Marshal(map[string]int{"likes": likes, "dislikes": dislikes, "actionTaken": actionTaken}) // Marshal the response data to JSON
		w.Write(data)                                                                                             // Write the response data
	} else {
		http.Error(w, "You must be logged in to like or dislike a comment", http.StatusUnauthorized) // Return error if user is not logged in
	}
}

// setCommentLikes sets the like or dislike for a comment
func setCommentLikes(commentID int, userId int, isDislike int) (int, int, int, error) {
	var id int                                                                                                                                                   // Variable to hold the postCommentLikeId
	var isDislikeDB int                                                                                                                                          // Variable to hold the current dislike status from the database
	err := DB.QueryRow("SELECT postCommentLikeId, likeType FROM PostCommentLikes WHERE commentId = ? and userId = ?", commentID, userId).Scan(&id, &isDislikeDB) // Check if there is a like or dislike from before
	if err == sql.ErrNoRows {                                                                                                                                    // If no likes or dislikes exist
		err = DB.QueryRow("Insert into PostCommentLikes (likeType, commentId, userId) values (?, ?, ?)", isDislike, commentID, userId).Scan() // Add a new like or dislike to the database
		if err != nil && err != sql.ErrNoRows {                                                                                               // If there is an error inserting the like or dislike
			return 0, 0, 0, err // Return error
		}
		likes, dislikes := getCommentLikes(commentID) // Get the updated likes and dislikes for the comment
		return likes, dislikes, isDislike, nil        // Return the updated likes, dislikes, and action taken
	} else if err != nil { // If there is a normal error
		return 0, 0, 0, err // Return error
	} else { // If a like or dislike already exists, update it
		deleted := false                                                                  // Variable to track if the like or dislike was deleted
		if (isDislikeDB == 1 && isDislike == 1) || (isDislikeDB == 0 && isDislike == 0) { // If the same action is being repeated, delete the like or dislike
			err = DB.QueryRow("delete from PostCommentLikes where postCommentLikeId = ?", id).Scan() // Delete the like or dislike
			if err != nil && err != sql.ErrNoRows {                                                  // If there is an error deleting the like or dislike
				return 0, 0, 0, err // Return error
			}
			deleted = true // Set deleted to true
		} else { // If a different action is being taken, update the like or dislike
			err = DB.QueryRow("update PostCommentLikes set likeType = ? where postCommentLikeId = ?", isDislike, id).Scan() // Update the like or dislike
			if err != nil && err != sql.ErrNoRows {                                                                         // If there is an error updating the like or dislike
				return 0, 0, 0, err // Return error
			}
		}
		likes, dislikes := getCommentLikes(commentID) // Get the updated likes and dislikes for the comment
		if deleted {                                  // If the like or dislike was deleted
			return likes, dislikes, -1, nil // Return the updated likes, dislikes, and action taken as -1 (deleted)
		} else {
			return likes, dislikes, isDislike, nil // Return the updated likes, dislikes, and action taken
		}
	}
}
