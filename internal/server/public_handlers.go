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

// handleAdminPanel serves the admin panel dashboard
func (s *Server) handleAdminPanel(w http.ResponseWriter, r *http.Request) {
	// Placeholder for admin panel
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Admin Panel - OnePage CMS</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 800px; margin: 0 auto; }
        .status { padding: 20px; background: #f0f8ff; border-radius: 5px; margin: 20px 0; }
        button { background: #007cba; color: white; border: none; padding: 8px 16px; border-radius: 4px; cursor: pointer; }
        button:hover { background: #005a87; }
        a { color: #007cba; text-decoration: none; }
        a:hover { text-decoration: underline; }
    </style>
    <script>
        function testStorage() {
            fetch('/admin/test-storage', { method: 'POST' })
                .then(response => response.json())
                .then(data => {
                    alert('Storage test completed! Check console for details.');
                    console.log('Storage test result:', data);
                })
                .catch(error => {
                    alert('Storage test failed: ' + error);
                    console.error('Storage test error:', error);
                });
        }
        
        function testTemplate() {
            fetch('/admin/test-template', { method: 'POST' })
                .then(response => response.json())
                .then(data => {
                    alert('Template test completed! Check console for details.');
                    console.log('Template test result:', data);
                })
                .catch(error => {
                    alert('Template test failed: ' + error);
                    console.error('Template test error:', error);
                });
        }
        
        function testContent() {
            fetch('/admin/test-content', { method: 'POST' })
                .then(response => response.json())
                .then(data => {
                    alert('Content test completed! Check console for details.');
                    console.log('Content test result:', data);
                })
                .catch(error => {
                    alert('Content test failed: ' + error);
                    console.error('Content test error:', error);
                });
        }
        
        function testSchema() {
            fetch('/admin/test-schema', { method: 'POST' })
                .then(response => response.json())
                .then(data => {
                    alert('Schema test completed! Check console for details.');
                    console.log('Schema test result:', data);
                })
                .catch(error => {
                    alert('Schema test failed: ' + error);
                    console.error('Schema test error:', error);
                });
        }
    </script>
</head>
<body>
    <div class="container">
        <h1>OnePage CMS - Admin Panel</h1>
        <div class="status">
            <h3>🚧 Under Construction</h3>
            <p>The admin panel is being built. Current status:</p>
            <ul>
                <li>✅ Basic HTTP server running</li>
                <li>✅ Static file serving</li>
                <li>✅ Directory structure created</li>
                <li>✅ File operations module (JSON/text read/write)</li>
                <li>✅ Backup system (.bak files)</li>
                <li>✅ Directory validation and creation</li>
                <li>✅ Template management (load/save/validate)</li>
                <li>✅ Default template generation</li>
                <li>✅ Content management (load/save/validate)</li>
                <li>✅ Content backup and restore</li>
                <li>✅ Schema management (JSON Schema validation)</li>
                <li>✅ Form generation from schema</li>
                <li>⏳ Authentication system (coming soon)</li>
                <li>⏳ Image management (coming soon)</li>
            </ul>
            <h4>Testing:</h4>
            <p><a href="/admin/files">📁 View Files</a> | <button onclick="testStorage()">🧪 Test Storage</button></p>
            <p><a href="/admin/template/info">📄 Template Info</a> | <button onclick="testTemplate()">🧪 Test Template</button></p>
            <p><a href="/admin/content/info">📝 Content Info</a> | <button onclick="testContent()">🧪 Test Content</button></p>
            <p><a href="/admin/schema/info">📋 Schema Info</a> | <button onclick="testSchema()">🧪 Test Schema</button></p>
        </div>
        <p><a href="/">← Back to public page</a></p>
    </div>
</body>
</html>`)
}
