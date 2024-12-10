package main

import (
	"database/sql" // Package for SQL database operations
	"fmt"          // Package for formatted I/O

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// getDatabase initializes the database and creates necessary tables
func getDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "forum.db") // Open a connection to the SQLite database
	if err != nil {                            // If there is an error opening the database
		return nil, err // Return the error
	}

	// SQL commands to create tables if they do not exist
	sqlCommands := []string{
		`CREATE TABLE IF NOT EXISTS Users(
			userId INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			age INTEGER NOT NULL,
			gender TEXT NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS Sessions(
			sessionId INTEGER PRIMARY KEY AUTOINCREMENT,
			sessionData TEXT,
			createdAt TEXT,
			lastAccessedAt TEXT,
			expiry TEXT,
			userId INT,
			FOREIGN KEY (userId) REFERENCES Users(userId)
		);`,
		`CREATE TABLE IF NOT EXISTS Posts(
			postId INTEGER PRIMARY KEY AUTOINCREMENT,
			postTitle TEXT NOT NULL,
			postContent TEXT NOT NULL,
			timePosted TEXT,
			userId INT,
			FOREIGN KEY (userId) REFERENCES Users(userId)
		);`,
		`CREATE TABLE IF NOT EXISTS PostLikes(
			postLikeId INTEGER PRIMARY KEY AUTOINCREMENT,
			postIsDisliked INT,
			postId INT,
			userId INT,
			FOREIGN KEY (postId) REFERENCES Posts(postId),
			FOREIGN KEY (userId) REFERENCES Users(userId)
		);`,
		`CREATE TABLE IF NOT EXISTS Comments(
			commentId INTEGER PRIMARY KEY AUTOINCREMENT,
			commentContent TEXT,
			likes INT,
			dislikes INT,
			userId INT,
			postId INT,
			FOREIGN KEY (userId) REFERENCES Users(userId),
			FOREIGN KEY (postId) REFERENCES Posts(postId)
		);`,
		`CREATE TABLE IF NOT EXISTS PostCommentLikes(
			postCommentLikeId INTEGER PRIMARY KEY AUTOINCREMENT,
			likeType INT,
			commentId INT,
			userId INT,
			FOREIGN KEY (commentId) REFERENCES Comments(commentId),
			FOREIGN KEY (userId) REFERENCES Users(userId)
		);`,
		`CREATE TABLE IF NOT EXISTS Categories(
			catId INTEGER PRIMARY KEY AUTOINCREMENT,
			catName TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS CategoriesPosts(
			catPostsId INTEGER PRIMARY KEY AUTOINCREMENT,
			postId INT,
			catId INT,
			FOREIGN KEY (postId) REFERENCES Posts(postId),
			FOREIGN KEY (catId) REFERENCES Categories(catId)
		);`,
		`CREATE TABLE IF NOT EXISTS Images(
			imageId INTEGER PRIMARY KEY AUTOINCREMENT,
			fileName TEXT,
			fileType TEXT,
			fileSize TEXT,
			postId INT,
			FOREIGN KEY (postId) REFERENCES Posts(postId)
		);`,
		`CREATE TABLE IF NOT EXISTS PrivateMessage (
			messageId INTEGER PRIMARY KEY AUTOINCREMENT,
			senderUsername TEXT,
			receiverUsername TEXT,
			messageText TEXT,
			timeSent DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (senderUsername) REFERENCES Users(username),
			FOREIGN KEY (receiverUsername) REFERENCES Users(username)
		);`,
	}

	// Execute each SQL command to create tables
	for _, command := range sqlCommands {
		statement, err := db.Prepare(command) // Prepare the SQL statement
		if err != nil {                       // If there is an error preparing the statement
			fmt.Println("error while running the following command:", command) // Print the command that caused the error
			fmt.Println(err)                                                   // Print the error
			return nil, err                                                    // Return the error
		}
		statement.Exec() // Execute the SQL statement
	}

	// Add default categories if there are none
	categories := []string{"Tech", "Sport", "News", "Culture", "Cars", "Nature", "Art"} // Built-in categories
	var count int
	err = db.QueryRow("SELECT count(catId) FROM Categories").Scan(&count) // Query to count the number of categories
	if err != nil && err != sql.ErrNoRows {                               // If there is an error executing the query
		return nil, err // Return the error
	}
	if count == 0 { // If there are no categories, add the default ones
		for _, cat := range categories {
			_, err = db.Exec("INSERT INTO Categories (catName) VALUES (?)", cat) // Insert the default categories
			if err != nil {                                                      // If there is an error inserting the categories
				return nil, err // Return the error
			}
		}
	}

	return db, nil // Return the database connection
}
