package forum

import (
	"encoding/json" // Package for encoding and decoding JSON
	"net/http"      // Package for HTTP client and server implementations
)

// GetCategoriesHandler handles HTTP requests to fetch all categories
func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := getCategories() // Call function to get all categories
	if err != nil {                    // If there is an error fetching categories
		http.Error(w, err.Error(), http.StatusInternalServerError) // Return error response
		return                                                     // Exit function
	}
	w.Header().Set("Content-Type", "application/json") // Set response content type to JSON
	json.NewEncoder(w).Encode(categories)              // Encode categories slice to JSON and write to response
}

// getCategories retrieves all categories from the database
func getCategories() ([]Category, error) {
	categories := []Category{}                                     // Slice to hold categories
	rows, err := DB.Query("SELECT catid, catname FROM Categories") // Execute query to get all categories
	if err != nil {                                                // If there is an error executing the query
		return nil, err // Return error
	}
	defer rows.Close() // Ensure rows are closed when function exits

	// Iterate over the rows and scan into Category structs
	for rows.Next() { // Loop through each row in the result set
		var category Category                                     // Variable to hold a category
		err := rows.Scan(&category.CategoryID, &category.CatName) // Scan row into category struct
		if err != nil {                                           // If there is an error scanning the row
			return nil, err // Return error
		}
		categories = append(categories, category) // Add the category to the slice
	}

	// Check for any errors encountered while iterating over rows
	if err := rows.Err(); err != nil { // If there is an error iterating over the rows
		return nil, err // Return error
	}
	return categories, nil // Return the slice of categories
}

// getPostCategories retrieves categories for a specific post from the database
func getPostCategories(postId int) ([]Category, error) {
	categories := []Category{}                                                                                                               // Slice to hold categories
	rows, err := DB.Query("SELECT c.catid, c.catname FROM Categories c, categoriesposts p where c.catid = p.catid and p.postid = ?", postId) // Execute query to get categories for a specific post
	if err != nil {                                                                                                                          // If there is an error executing the query
		return nil, err // Return error
	}
	defer rows.Close() // Ensure rows are closed when function exits

	// Iterate over the rows and scan into Category structs
	for rows.Next() { // Loop through each row in the result set
		var category Category                                     // Variable to hold a category
		err := rows.Scan(&category.CategoryID, &category.CatName) // Scan row into category struct
		if err != nil {                                           // If there is an error scanning the row
			return nil, err // Return error
		}
		categories = append(categories, category) // Add the category to the slice
	}

	// Check for any errors encountered while iterating over rows
	if err := rows.Err(); err != nil { // If there is an error iterating over the rows
		return nil, err // Return error
	}
	return categories, nil // Return the slice of categories
}
