package forum

import (
	"fmt" // Package for formatted I/O
)

// getPostLikes retrieves the number of likes and dislikes for a post, along with the user IDs who liked or disliked the post
func getPostLikes(postId int) (int, int, []int, []int) {
	var likes int           // Variable to hold the number of likes
	var dislikes int        // Variable to hold the number of dislikes
	var likedUsers []int    // Slice to hold user IDs who liked the post
	var dislikedUsers []int // Slice to hold user IDs who disliked the post

	// Query to get the number of likes and dislikes for the post
	err := DB.QueryRow("select (SELECT count(postId) FROM postLikes WHERE postIsDisliked != 1 and postId = ?), (SELECT count(postId) FROM postLikes WHERE postIsDisliked = 1 and postId = ?)", postId, postId).Scan(&likes, &dislikes)
	if err != nil { // If there is an error executing the query
		fmt.Println(err)                       // Print the error
		return 0, 0, likedUsers, dislikedUsers // Return default values
	}

	// Query to get the user IDs who liked the post
	rows, err := DB.Query("SELECT userID FROM postLikes where postIsDisliked != 1 and postid = ?", postId)
	if err != nil { // If there is an error executing the query
		fmt.Println(err)                       // Print the error
		return 0, 0, likedUsers, dislikedUsers // Return default values
	}
	defer rows.Close() // Ensure rows are closed when function exits

	// Iterate over the rows and scan into likedUsers slice
	for rows.Next() {
		var userId int            // Variable to hold a user ID
		err := rows.Scan(&userId) // Scan row into userId
		if err != nil {           // If there is an error scanning the row
			fmt.Println(err)                       // Print the error
			return 0, 0, likedUsers, dislikedUsers // Return default values
		}
		likedUsers = append(likedUsers, userId) // Add the user ID to the likedUsers slice
	}

	// Query to get the user IDs who disliked the post
	rows, err = DB.Query("SELECT userID FROM postLikes where postIsDisliked == 1 and postid = ?", postId)
	if err != nil { // If there is an error executing the query
		fmt.Println(err)                       // Print the error
		return 0, 0, likedUsers, dislikedUsers // Return default values
	}
	defer rows.Close() // Ensure rows are closed when function exits

	// Iterate over the rows and scan into dislikedUsers slice
	for rows.Next() {
		var userId int            // Variable to hold a user ID
		err := rows.Scan(&userId) // Scan row into userId
		if err != nil {           // If there is an error scanning the row
			fmt.Println(err)                       // Print the error
			return 0, 0, likedUsers, dislikedUsers // Return default values
		}
		dislikedUsers = append(dislikedUsers, userId) // Add the user ID to the dislikedUsers slice
	}
	return likes, dislikes, likedUsers, dislikedUsers // Return the number of likes, dislikes, and the user IDs
}

// getCommentLikes retrieves the number of likes and dislikes for a comment
func getCommentLikes(commentId int) (int, int) {
	var likes int    // Variable to hold the number of likes
	var dislikes int // Variable to hold the number of dislikes

	// Query to get the number of likes and dislikes for the comment
	err := DB.QueryRow("select (SELECT count(commentId) FROM postCommentLikes WHERE likeType != 1 and commentId = ?), (SELECT count(commentId) FROM postCommentLikes where likeType = 1 and  commentId = ?);", commentId, commentId).Scan(&likes, &dislikes)
	if err != nil { // If there is an error executing the query
		fmt.Println(err) // Print the error
		return 0, 0      // Return default values
	}
	return likes, dislikes // Return the number of likes and dislikes
}
