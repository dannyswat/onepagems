package server

import (
	"fmt"
	"net/http"
	"os"
)

// handlePublicPage serves the main public page
func (s *Server) handlePublicPage(w http.ResponseWriter, r *http.Request) {
	// For now, serve a simple placeholder
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Check if index.html exists
	indexPath := "index.html"
	if _, err := os.Stat(indexPath); err == nil {
		http.ServeFile(w, r, indexPath)
		return
	}

	// Serve placeholder content
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>OnePage CMS</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
    <h1>Welcome to OnePage CMS</h1>
    <p>Your website will appear here after you configure it through the admin panel.</p>
    <p><a href="/admin">Go to Admin Panel</a></p>
</body>
</html>`)
}

// handleHealth returns health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok","message":"OnePage CMS is running"}`)
}
