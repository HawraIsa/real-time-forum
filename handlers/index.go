package forum

import (
	"net/http"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {

	// Check if this is a direct post URL (e.g., /post/123)
	if strings.HasPrefix(r.URL.Path, "/post/") {
		// Still serve the SPA, but let the frontend handle the routing
		r.URL.Path = "/"
	} else if r.URL.Path == "/404" {
		http.ServeFile(w, r, "template/SPA.html")
	} else if r.URL.Path != "/" {
		// Handle 404 for any other unknown paths
		// http.NotFound(w, r)
		http.Redirect(w, r, "/404", http.StatusTemporaryRedirect)
		return
	}

	// All other paths return 404
	// http.NotFound(w, r)
	http.ServeFile(w, r, "template/SPA.html")
}
