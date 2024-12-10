package forum

import (
	"fmt"      // Package for formatted I/O
	"net/http" // Package for HTTP client and server implementations
)

// returnJson sends a JSON response with a specified status code
func returnJson(w http.ResponseWriter, json string, statusCode int) {
	w.WriteHeader(statusCode)                          // Set the HTTP status code for the response
	w.Header().Set("Content-Type", "application/json") // Set the Content-Type header to application/json
	fmt.Fprint(w, json)                                // Write the JSON string to the response
}
