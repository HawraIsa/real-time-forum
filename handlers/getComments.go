package forum

// GetCommentsForPost retrieves all comments for a post from the database
func GetCommentsForPost(postID string) ([]Comment, error) {
	var comments []Comment // Slice to hold the comments

	// Query the database to retrieve all comments for the specified post
	rows, err := DB.Query("SELECT commentId, commentContent, userId, postId FROM Comments WHERE PostID = ?", postID) // Execute query to get comments for the post
	if err != nil {                                                                                                  // If there is an error executing the query
		return nil, err // Return error
	}
	defer rows.Close() // Ensure rows are closed when function exits

	// Iterate over the rows and scan into Comment structs
	for rows.Next() { // Loop through each row in the result set
		var comment Comment                                                                             // Variable to hold a comment
		err := rows.Scan(&comment.CommentID, &comment.CommentContent, &comment.UserID, &comment.PostID) // Scan row into comment struct
		if err != nil {                                                                                 // If there is an error scanning the row
			return nil, err // Return error
		}
		comment.UserName, err = getUserName(comment.UserID) // Get the username of the comment author
		if err != nil {                                     // If there is an error getting the username
			return nil, err // Return error
		}
		comment.Likes, comment.Dislikes = getCommentLikes(comment.CommentID) // Get the likes and dislikes for the comment
		comments = append(comments, comment)                                 // Add the comment to the slice
	}

	// Check for any errors encountered while iterating over rows
	if err := rows.Err(); err != nil { // If there is an error iterating over the rows
		return nil, err // Return error
	}

	return comments, nil // Return the slice of comments
}
