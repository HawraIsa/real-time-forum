package forum

import (
	"database/sql"
	"time"
)

// Define a struct to represent a user
type User struct {
	UserID    int    `json:"user_id"`  // Primary key
	Username  string `json:"username"` // Nickname
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// Define a struct to represent a session
type Session struct {
	SessionID      int       `json:"session_id"` // Primary key
	SessionData    string    `json:"session_data"`
	CreatedAt      time.Time `json:"created_at"`
	LastAccessedAt time.Time `json:"last_accessed_at"`
	Expiry         time.Time `json:"expiry"`
	UserID         int       `json:"user_id"` // Foreign key referencing UserID in Users table
}

// Define a struct to represent a post
type Post struct {
	PostID         int        `json:"post_id"` // Primary key
	PostTitle      string     `json:"post_title"`
	PostContent    string     `json:"post_content"`
	TimePosted     time.Time  `json:"time_posted"`
	UserID         int        `json:"user_id"` // Foreign key referencing UserID in Users table
	UserName       string     `json:"user_name"`
	LikedByUser    bool       `json:"liked_by_user"`
	DislikedByUser bool       `json:"disliked_by_user"`
	Likes          int        `json:"likes"`
	Dislikes       int        `json:"dislikes"`
	Categories     []Category `json:"categories"`
	Comments       []Comment  `json:"comments"`
	CommentsCount  int        `json:"comments_count"`
}

// Define a struct to represent a comment
type Comment struct {
	CommentID      int    `json:"comment_id"` // Primary key
	CommentContent string `json:"comment_content"`
	Likes          int    `json:"likes"`
	Dislikes       int    `json:"dislikes"`
	UserID         int    `json:"user_id"` // Foreign key referencing UserID in Users table
	UserName       string `json:"user_name"`
	PostID         int    `json:"post_id"` // Foreign key referencing PostID in Posts table
}

// Define a struct to represent a like/dislike for a comment
type PostCommentLike struct {
	PostCommentLikeID int `json:"post_comment_like_id"` // Primary key
	LikeType          int `json:"like_type"`            // 0 for dislike, 1 for like
	CommentID         int `json:"comment_id"`           // Foreign key referencing CommentID in Comments table
	UserID            int `json:"user_id"`              // Foreign key referencing UserID in Users table
}

// Define a struct to represent a category
type Category struct {
	CategoryID int    `json:"category_id"` // Primary key
	CatName    string `json:"cat_name"`
}

type CategoriesPosts struct {
	CatPostsID int `json:"cat_posts_id"` // Primary key
	PostID     int `json:"post_id"`
	CatID      int `json:"cat_id"`
}

// Define a struct to represent a post image
type PostImage struct {
	ImageID  int    `json:"image_id"` // Primary key
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	FileSize int    `json:"file_size"`
	PostID   int    `json:"post_id"` // Foreign key referencing PostID in Posts table
}

// Define a struct to represent a private message
type PrivateMessages struct {
	MessageID        int       `json:"message_id"`       // Primary key
	SenderUsername   string    `json:"senderUsername"`   // Foreign key referencing UserID in Users table
	ReceiverUsername string    `json:"receiverUsername"` // Foreign key referencing UserID in Users table
	MessageText      string    `json:"messageText"`      // Message content
	TimeSent         time.Time `json:"timeSent"`         // Timestamp of when the message was sent
}

var DB *sql.DB
