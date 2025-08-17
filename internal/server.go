package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Server represents the HTTP server
type Server struct {
	Config          *Config
	Storage         *FileStorage
	TemplateManager *TemplateManager
	ContentManager  *ContentManager
	SchemaManager   *SchemaManager
	AuthManager     *AuthManager
	Mux             *http.ServeMux
}

// NewServer creates a new server instance
func NewServer(config *Config) *Server {
	storage := NewFileStorage(config.DataDir)
	server := &Server{
		Config:          config,
		Storage:         storage,
		TemplateManager: NewTemplateManager(storage),
		ContentManager:  NewContentManager(storage, config.DataDir),
		SchemaManager:   NewSchemaManager(storage, config.DataDir),
		AuthManager:     NewAuthManager(config),
		Mux:             http.NewServeMux(),
	}

	// Set up routes
	server.setupRoutes()

	return server
}

// setupRoutes configures all the HTTP routes
func (s *Server) setupRoutes() {
	// Static file serving
	s.Mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(s.Config.StaticDir))))
	s.Mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(filepath.Join(s.Config.DataDir, "images")))))

	// Public routes
	s.Mux.HandleFunc("/", s.handlePublicPage)
	s.Mux.HandleFunc("/health", s.handleHealth)

	// Authentication routes (not protected)
	s.Mux.HandleFunc("/admin/login", s.handleAdminLogin)
	s.Mux.HandleFunc("/admin/logout", s.handleAdminLogout)

	// Protected admin routes
	s.Mux.HandleFunc("/admin", s.AuthManager.RequireAuth(s.handleAdminPanel))

	// File management test endpoints (protected)
	s.Mux.HandleFunc("/admin/files", s.AuthManager.RequireAuth(s.handleFilesList))
	s.Mux.HandleFunc("/admin/test-storage", s.AuthManager.RequireAuth(s.handleTestStorage))

	// Template management endpoints (protected)
	s.Mux.HandleFunc("/admin/template", s.AuthManager.RequireAuth(s.handleTemplate))
	s.Mux.HandleFunc("/admin/template/info", s.AuthManager.RequireAuth(s.handleTemplateInfo))
	s.Mux.HandleFunc("/admin/template/restore", s.AuthManager.RequireAuth(s.handleTemplateRestore))
	s.Mux.HandleFunc("/admin/test-template", s.AuthManager.RequireAuth(s.handleTestTemplate))

	// Content management endpoints (protected)
	s.Mux.HandleFunc("/admin/content", s.AuthManager.RequireAuth(s.handleContent))
	s.Mux.HandleFunc("/admin/content/info", s.AuthManager.RequireAuth(s.handleContentInfo))
	s.Mux.HandleFunc("/admin/content/restore", s.AuthManager.RequireAuth(s.handleContentRestore))
	s.Mux.HandleFunc("/admin/content/export", s.AuthManager.RequireAuth(s.handleContentExport))
	s.Mux.HandleFunc("/admin/content/import", s.AuthManager.RequireAuth(s.handleContentImport))
	s.Mux.HandleFunc("/admin/test-content", s.AuthManager.RequireAuth(s.handleTestContent))

	// Schema management endpoints (protected)
	s.Mux.HandleFunc("/admin/schema", s.AuthManager.RequireAuth(s.handleSchema))
	s.Mux.HandleFunc("/admin/schema/info", s.AuthManager.RequireAuth(s.handleSchemaInfo))
	s.Mux.HandleFunc("/admin/schema/restore", s.AuthManager.RequireAuth(s.handleSchemaRestore))
	s.Mux.HandleFunc("/admin/schema/export", s.AuthManager.RequireAuth(s.handleSchemaExport))
	s.Mux.HandleFunc("/admin/schema/import", s.AuthManager.RequireAuth(s.handleSchemaImport))
	s.Mux.HandleFunc("/admin/schema/validate", s.AuthManager.RequireAuth(s.handleSchemaValidate))
	s.Mux.HandleFunc("/admin/schema/form", s.AuthManager.RequireAuth(s.handleSchemaForm))
	s.Mux.HandleFunc("/admin/test-schema", s.AuthManager.RequireAuth(s.handleTestSchema))

	// Schema parser endpoints (protected)
	s.Mux.HandleFunc("/admin/schema/analyze", s.AuthManager.RequireAuth(s.handleSchemaAnalyze))
	s.Mux.HandleFunc("/admin/schema/field-metadata", s.AuthManager.RequireAuth(s.handleSchemaFieldMetadata))
	s.Mux.HandleFunc("/admin/schema/validation-rules", s.AuthManager.RequireAuth(s.handleSchemaValidationRules))
	s.Mux.HandleFunc("/admin/schema/field-types", s.AuthManager.RequireAuth(s.handleSchemaFieldTypes))
	s.Mux.HandleFunc("/admin/schema/required-fields", s.AuthManager.RequireAuth(s.handleSchemaRequiredFields))
	s.Mux.HandleFunc("/admin/schema/validate-field", s.AuthManager.RequireAuth(s.handleSchemaValidateField))

	// Schema validator endpoints (protected)
	s.Mux.HandleFunc("/admin/schema/validate-content", s.AuthManager.RequireAuth(s.handleSchemaValidateContent))
	s.Mux.HandleFunc("/admin/schema/validate-field-detailed", s.AuthManager.RequireAuth(s.handleSchemaValidateFieldDetailed))
	s.Mux.HandleFunc("/admin/schema/validation-report", s.AuthManager.RequireAuth(s.handleSchemaValidationReport))

	// Authentication status endpoints (protected)
	s.Mux.HandleFunc("/admin/auth/status", s.AuthManager.RequireAuth(s.handleAuthStatus))
	s.Mux.HandleFunc("/admin/auth/sessions", s.AuthManager.RequireAuth(s.handleAuthSessions))
	s.Mux.HandleFunc("/admin/auth/change-password", s.AuthManager.RequireAuth(s.handleChangePassword))

	log.Println("Routes configured:")
	log.Println("  GET  /               - Public page")
	log.Println("  GET  /health         - Health check")
	log.Println("  GET  /static/        - Static files")
	log.Println("  GET  /images/        - Image files")
	log.Println("  GET  /admin          - Admin panel")
	log.Println("  POST /admin/login    - Admin login")
	log.Println("  POST /admin/logout   - Admin logout")
	log.Println("  GET  /admin/files    - List files (test)")
	log.Println("  POST /admin/test-storage - Test storage operations")
	log.Println("  GET/POST /admin/template - Template management")
	log.Println("  GET  /admin/template/info - Template information")
	log.Println("  POST /admin/template/restore - Restore template")
	log.Println("  POST /admin/test-template - Test template operations")
	log.Println("  GET/POST /admin/content - Content management")
	log.Println("  GET  /admin/content/info - Content information")
	log.Println("  POST /admin/content/restore - Restore content")
	log.Println("  GET  /admin/content/export - Export content")
	log.Println("  POST /admin/content/import - Import content")
	log.Println("  POST /admin/test-content - Test content operations")
	log.Println("  GET/POST /admin/schema - Schema management")
	log.Println("  GET  /admin/schema/info - Schema information")
	log.Println("  POST /admin/schema/restore - Restore schema")
	log.Println("  GET  /admin/schema/export - Export schema")
	log.Println("  POST /admin/schema/import - Import schema")
	log.Println("  POST /admin/schema/validate - Validate data against schema")
	log.Println("  GET  /admin/schema/form - Generate form from schema")
	log.Println("  POST /admin/test-schema - Test schema operations")
	log.Println("  GET  /admin/schema/analyze - Comprehensive schema analysis")
	log.Println("  GET  /admin/schema/field-metadata - Get field metadata (query: field)")
	log.Println("  GET  /admin/schema/validation-rules - Get all validation rules")
	log.Println("  GET  /admin/schema/field-types - Get field types mapping")
	log.Println("  GET  /admin/schema/required-fields - Get required/optional fields")
	log.Println("  POST /admin/schema/validate-field - Validate single field value")
	log.Println("  POST /admin/schema/validate-content - Comprehensive content validation")
	log.Println("  POST /admin/schema/validate-field-detailed - Detailed field validation")
	log.Println("  POST /admin/schema/validation-report - Generate validation report")
	log.Println("  GET  /admin/auth/status - Authentication status")
	log.Println("  GET  /admin/auth/sessions - List active sessions")
	log.Println("  POST /admin/auth/change-password - Change password")
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Ensure data directories exist
	if err := s.ensureDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	addr := ":" + s.Config.Port
	log.Printf("Starting server on http://localhost%s", addr)
	log.Printf("Admin panel: http://localhost%s/admin", addr)

	return http.ListenAndServe(addr, s.Mux)
}

// ensureDirectories creates necessary directories if they don't exist
func (s *Server) ensureDirectories() error {
	// Use storage to ensure data directories
	if err := s.Storage.EnsureDirectories(); err != nil {
		return err
	}

	// Ensure other directories
	dirs := []string{
		s.Config.StaticDir,
		s.Config.TemplatesDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		log.Printf("Ensured directory exists: %s", dir)
	}

	return nil
}

// Basic handlers (will be implemented properly in later phases)

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

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok","message":"OnePage CMS is running"}`)
}

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

func (s *Server) handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Serve login form
		s.serveLoginForm(w, r)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse login credentials
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Attempt login
	session, err := s.AuthManager.Login(username, password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid credentials",
		})
		return
	}

	// Set session cookie
	cookie := s.AuthManager.CreateSessionCookie(session.ID)
	http.SetCookie(w, cookie)

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"message":    "Login successful",
		"session_id": session.ID,
		"expires_at": session.ExpiresAt,
	})
}

func (s *Server) handleAdminLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get session from request
	session, err := s.AuthManager.GetSessionFromRequest(r)
	if err == nil {
		// Logout the session
		s.AuthManager.Logout(session.ID)
	}

	// Clear session cookie
	cookie := s.AuthManager.ClearSessionCookie()
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}

// serveLoginForm serves the login HTML form
func (s *Server) serveLoginForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OnePage CMS - Admin Login</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            margin: 0;
            padding: 0;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .login-container {
            background: white;
            padding: 2rem;
            border-radius: 10px;
            box-shadow: 0 15px 35px rgba(0,0,0,0.1);
            width: 100%%;
            max-width: 400px;
        }
        .login-header {
            text-align: center;
            margin-bottom: 2rem;
        }
        .login-header h1 {
            color: #333;
            margin: 0;
            font-size: 1.8rem;
        }
        .login-header p {
            color: #666;
            margin: 0.5rem 0 0 0;
        }
        .form-group {
            margin-bottom: 1.5rem;
        }
        .form-group label {
            display: block;
            margin-bottom: 0.5rem;
            color: #333;
            font-weight: 500;
        }
        .form-group input {
            width: 100%%;
            padding: 0.75rem;
            border: 1px solid #ddd;
            border-radius: 5px;
            font-size: 1rem;
            box-sizing: border-box;
        }
        .form-group input:focus {
            border-color: #667eea;
            outline: none;
            box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
        }
        .login-button {
            width: 100%%;
            padding: 0.75rem;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 5px;
            font-size: 1rem;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        .login-button:hover {
            background: #5a6fd8;
        }
        .error-message {
            background: #fee;
            color: #c00;
            padding: 0.75rem;
            border-radius: 5px;
            margin-bottom: 1rem;
            display: none;
        }
        .success-message {
            background: #efe;
            color: #060;
            padding: 0.75rem;
            border-radius: 5px;
            margin-bottom: 1rem;
            display: none;
        }
    </style>
</head>
<body>
    <div class="login-container">
        <div class="login-header">
            <h1>OnePage CMS</h1>
            <p>Admin Login</p>
        </div>
        
        <div id="error-message" class="error-message"></div>
        <div id="success-message" class="success-message"></div>
        
        <form id="login-form">
            <div class="form-group">
                <label for="username">Username:</label>
                <input type="text" id="username" name="username" required>
            </div>
            
            <div class="form-group">
                <label for="password">Password:</label>
                <input type="password" id="password" name="password" required>
            </div>
            
            <button type="submit" class="login-button">Login</button>
        </form>
    </div>

    <script>
        document.getElementById('login-form').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const formData = new FormData(this);
            const errorDiv = document.getElementById('error-message');
            const successDiv = document.getElementById('success-message');
            
            // Hide previous messages
            errorDiv.style.display = 'none';
            successDiv.style.display = 'none';
            
            try {
                const response = await fetch('/admin/login', {
                    method: 'POST',
                    body: formData
                });
                
                const data = await response.json();
                
                if (response.ok && data.success) {
                    successDiv.textContent = 'Login successful! Redirecting...';
                    successDiv.style.display = 'block';
                    setTimeout(() => {
                        window.location.href = '/admin';
                    }, 1000);
                } else {
                    errorDiv.textContent = data.error || 'Login failed';
                    errorDiv.style.display = 'block';
                }
            } catch (error) {
                errorDiv.textContent = 'Network error: ' + error.message;
                errorDiv.style.display = 'block';
            }
        });
    </script>
</body>
</html>`)
}

// handleAuthStatus returns current authentication status
func (s *Server) handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, ok := SessionFromContext(r.Context())
	if !ok {
		http.Error(w, "No session found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated":   true,
		"username":        session.Username,
		"session_id":      session.ID,
		"created_at":      session.CreatedAt,
		"expires_at":      session.ExpiresAt,
		"active_sessions": s.AuthManager.GetActiveSessions(),
	})
}

// handleAuthSessions lists all active sessions
func (s *Server) handleAuthSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessions := s.AuthManager.ListSessions()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sessions": sessions,
		"count":    len(sessions),
	})
}

// handleChangePassword changes the admin password
func (s *Server) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	if currentPassword == "" || newPassword == "" || confirmPassword == "" {
		http.Error(w, "All password fields are required", http.StatusBadRequest)
		return
	}

	if newPassword != confirmPassword {
		http.Error(w, "New passwords do not match", http.StatusBadRequest)
		return
	}

	if err := s.AuthManager.ChangePassword(currentPassword, newPassword); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password changed successfully",
	})
}

// handleFilesList lists all files in the data directory
func (s *Server) handleFilesList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	files, err := s.Storage.ListFiles()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list files: %v", err), http.StatusInternalServerError)
		return
	}

	response := NewAPIResponse(true, "Files listed successfully")
	response.SetData(files)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleTestStorage demonstrates file storage operations
func (s *Server) handleTestStorage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Test data
	testData := map[string]interface{}{
		"message":    "Hello from file storage test",
		"timestamp":  time.Now().Format(time.RFC3339),
		"test_array": []string{"item1", "item2", "item3"},
		"test_object": map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	// Test JSON operations
	filename := "test.json"

	// Write test file
	if err := s.Storage.WriteJSONFile(filename, testData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write test file: %v", err), http.StatusInternalServerError)
		return
	}

	// Read it back
	var readData map[string]interface{}
	if err := s.Storage.ReadJSONFile(filename, &readData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to read test file: %v", err), http.StatusInternalServerError)
		return
	}

	// Get backup info
	backupInfo, _ := s.Storage.GetBackupInfo(filename)

	// Test text file operations
	textFilename := "test.txt"
	textContent := fmt.Sprintf("Test text file created at %s\nThis demonstrates text file operations.", time.Now().Format(time.RFC3339))

	if err := s.Storage.WriteTextFile(textFilename, textContent); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write text file: %v", err), http.StatusInternalServerError)
		return
	}

	readText, err := s.Storage.ReadTextFile(textFilename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read text file: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare response
	result := map[string]interface{}{
		"json_write_success":  true,
		"json_read_success":   true,
		"json_data_matches":   fmt.Sprintf("%v", testData["message"]) == fmt.Sprintf("%v", readData["message"]),
		"text_write_success":  true,
		"text_read_success":   true,
		"text_content_length": len(readText),
		"backup_created":      backupInfo != nil,
		"files_created":       []string{filename, textFilename},
	}

	if backupInfo != nil {
		result["backup_info"] = map[string]interface{}{
			"backup_path": backupInfo.BackupPath,
			"created_at":  backupInfo.CreatedAt,
			"size":        backupInfo.Size,
		}
	}

	response := NewAPIResponse(true, "File storage test completed successfully")
	response.SetData(result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleTemplate handles template operations (GET to load, POST to save)
func (s *Server) handleTemplate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleTemplateGet(w, r)
	case "POST":
		s.handleTemplatePost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTemplateGet loads and returns the current template
func (s *Server) handleTemplateGet(w http.ResponseWriter, r *http.Request) {
	content, err := s.TemplateManager.LoadTemplate()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load template: %v", err), http.StatusInternalServerError)
		return
	}

	response := NewAPIResponse(true, "Template loaded successfully")
	response.SetData(map[string]interface{}{
		"content": content,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleTemplatePost saves a new template
func (s *Server) handleTemplatePost(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Template content is required", http.StatusBadRequest)
		return
	}

	// Save template
	if err := s.TemplateManager.SaveTemplate(content); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save template: %v", err), http.StatusBadRequest)
		return
	}

	response := NewAPIResponse(true, "Template saved successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleTemplateInfo returns information about the current template
func (s *Server) handleTemplateInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	info, err := s.TemplateManager.GetTemplateInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get template info: %v", err), http.StatusInternalServerError)
		return
	}

	// Get template variables
	content, err := s.TemplateManager.LoadTemplate()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load template for analysis: %v", err), http.StatusInternalServerError)
		return
	}

	variables, err := s.TemplateManager.GetTemplateVariables(content)
	if err != nil {
		// Don't fail completely if variable analysis fails
		variables = []string{"Error analyzing variables: " + err.Error()}
	}

	result := map[string]interface{}{
		"file_info": info,
		"variables": variables,
		"content_preview": func() string {
			if len(content) > 200 {
				return content[:200] + "..."
			}
			return content
		}(),
	}

	response := NewAPIResponse(true, "Template info retrieved successfully")
	response.SetData(result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleTemplateRestore restores template from backup
func (s *Server) handleTemplateRestore(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := s.TemplateManager.RestoreTemplate(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to restore template: %v", err), http.StatusInternalServerError)
		return
	}

	response := NewAPIResponse(true, "Template restored from backup successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleTestTemplate tests template functionality
func (s *Server) handleTestTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Test template operations
	results := make(map[string]interface{})

	// Test 1: Load default template
	content, err := s.TemplateManager.LoadTemplate()
	if err != nil {
		results["load_template"] = "Failed: " + err.Error()
	} else {
		results["load_template"] = "Success"
		results["template_size"] = len(content)
	}

	// Test 2: Validate template
	if content != "" {
		if err := s.TemplateManager.ValidateTemplate(content); err != nil {
			results["validate_template"] = "Failed: " + err.Error()
		} else {
			results["validate_template"] = "Success"
		}
	}

	// Test 3: Get template info
	if info, err := s.TemplateManager.GetTemplateInfo(); err != nil {
		results["template_info"] = "Failed: " + err.Error()
	} else {
		results["template_info"] = "Success"
		results["has_backup"] = info.HasBackup
		results["file_size"] = info.Size
	}

	// Test 4: Get template variables
	if variables, err := s.TemplateManager.GetTemplateVariables(content); err != nil {
		results["template_variables"] = "Failed: " + err.Error()
	} else {
		results["template_variables"] = "Success"
		results["variable_count"] = len(variables)
		results["variables"] = variables
	}

	// Test 5: Save a test template (minor modification)
	testContent := content + "\n<!-- Test modification at " + time.Now().Format(time.RFC3339) + " -->"
	if err := s.TemplateManager.SaveTemplate(testContent); err != nil {
		results["save_template"] = "Failed: " + err.Error()
	} else {
		results["save_template"] = "Success"
		results["backup_created"] = true
	}

	response := NewAPIResponse(true, "Template test completed")
	response.SetData(results)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleContent handles content management requests
func (s *Server) handleContent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Load and return current content
		content, err := s.ContentManager.LoadContent()
		if err != nil {
			response := NewAPIResponse(false, "Failed to load content: "+err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := NewAPIResponse(true, "Content loaded successfully")
		response.SetData(content)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	case "POST":
		// Update content
		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			response := NewAPIResponse(false, "Invalid JSON in request body: "+err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if err := s.ContentManager.UpdateContent(updates); err != nil {
			response := NewAPIResponse(false, "Failed to update content: "+err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := NewAPIResponse(true, "Content updated successfully")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleContentInfo returns information about the current content
func (s *Server) handleContentInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	summary, err := s.ContentManager.GetContentSummary()
	if err != nil {
		response := NewAPIResponse(false, "Failed to get content information: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := NewAPIResponse(true, "Content information retrieved")
	response.SetData(summary)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleContentRestore restores content from backup
func (s *Server) handleContentRestore(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := s.ContentManager.RestoreContent(); err != nil {
		response := NewAPIResponse(false, "Failed to restore content: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := NewAPIResponse(true, "Content restored from backup successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleContentExport exports content as JSON
func (s *Server) handleContentExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := s.ContentManager.ExportContent()
	if err != nil {
		response := NewAPIResponse(false, "Failed to export content: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=content-export.json")
	w.Write(data)
}

// handleContentImport imports content from JSON
func (s *Server) handleContentImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	var requestData struct {
		Content json.RawMessage `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		response := NewAPIResponse(false, "Invalid JSON in request body: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := s.ContentManager.ImportContent(requestData.Content); err != nil {
		response := NewAPIResponse(false, "Failed to import content: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := NewAPIResponse(true, "Content imported successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleTestContent tests content management operations
func (s *Server) handleTestContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	results := make(map[string]interface{})

	// Test 1: Load current content
	content, err := s.ContentManager.LoadContent()
	if err != nil {
		results["load_content"] = "Failed: " + err.Error()
	} else {
		results["load_content"] = "Success"
		results["content_title"] = content.Title
		results["sections_count"] = len(content.Sections)
	}

	// Test 2: Get content summary
	if summary, err := s.ContentManager.GetContentSummary(); err != nil {
		results["content_summary"] = "Failed: " + err.Error()
	} else {
		results["content_summary"] = "Success"
		results["summary"] = summary
	}

	// Test 3: Update content (test update)
	testUpdates := map[string]interface{}{
		"description": "Test description updated at " + time.Now().Format(time.RFC3339),
	}
	if err := s.ContentManager.UpdateContent(testUpdates); err != nil {
		results["update_content"] = "Failed: " + err.Error()
	} else {
		results["update_content"] = "Success"
		results["backup_created"] = true
	}

	// Test 4: Export content
	if data, err := s.ContentManager.ExportContent(); err != nil {
		results["export_content"] = "Failed: " + err.Error()
	} else {
		results["export_content"] = "Success"
		results["export_size"] = len(data)
	}

	// Test 5: Backup content
	if err := s.ContentManager.BackupContent(); err != nil {
		results["backup_content"] = "Failed: " + err.Error()
	} else {
		results["backup_content"] = "Success"
	}

	response := NewAPIResponse(true, "Content test completed")
	response.SetData(results)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSchema handles schema management requests
func (s *Server) handleSchema(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Load and return current schema
		schema, err := s.SchemaManager.LoadSchema()
		if err != nil {
			response := NewAPIResponse(false, "Failed to load schema: "+err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := NewAPIResponse(true, "Schema loaded successfully")
		response.SetData(schema)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	case "POST":
		// Update schema
		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			response := NewAPIResponse(false, "Invalid JSON in request body: "+err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if err := s.SchemaManager.UpdateSchema(updates); err != nil {
			response := NewAPIResponse(false, "Failed to update schema: "+err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := NewAPIResponse(true, "Schema updated successfully")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleSchemaInfo returns information about the current schema
func (s *Server) handleSchemaInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	info, err := s.SchemaManager.GetSchemaInfo()
	if err != nil {
		response := NewAPIResponse(false, "Failed to get schema information: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := NewAPIResponse(true, "Schema information retrieved")
	response.SetData(info)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSchemaRestore restores schema from backup
func (s *Server) handleSchemaRestore(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := s.SchemaManager.RestoreSchema(); err != nil {
		response := NewAPIResponse(false, "Failed to restore schema: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := NewAPIResponse(true, "Schema restored from backup successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSchemaExport exports schema as JSON
func (s *Server) handleSchemaExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := s.SchemaManager.ExportSchema()
	if err != nil {
		response := NewAPIResponse(false, "Failed to export schema: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=schema-export.json")
	w.Write(data)
}

// handleSchemaImport imports schema from JSON
func (s *Server) handleSchemaImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	var requestData struct {
		Schema json.RawMessage `json:"schema"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		response := NewAPIResponse(false, "Invalid JSON in request body: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := s.SchemaManager.ImportSchema(requestData.Schema); err != nil {
		response := NewAPIResponse(false, "Failed to import schema: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := NewAPIResponse(true, "Schema imported successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSchemaValidate validates content against the current schema
func (s *Server) handleSchemaValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	var requestData struct {
		Data interface{} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		response := NewAPIResponse(false, "Invalid JSON in request body: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := s.SchemaManager.ValidateAgainstSchema(requestData.Data); err != nil {
		response := NewAPIResponse(false, "Validation failed: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := NewAPIResponse(true, "Data is valid against schema")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSchemaForm generates form fields from schema
func (s *Server) handleSchemaForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fields, err := s.SchemaManager.GenerateFormFromSchema()
	if err != nil {
		response := NewAPIResponse(false, "Failed to generate form from schema: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	form := GeneratedForm{
		Fields: fields,
		Action: "/admin/content",
		Method: "POST",
	}

	response := NewAPIResponse(true, "Form generated from schema")
	response.SetData(form)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleTestSchema tests schema management operations
func (s *Server) handleTestSchema(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	results := make(map[string]interface{})

	// Test 1: Load current schema
	schema, err := s.SchemaManager.LoadSchema()
	if err != nil {
		results["load_schema"] = "Failed: " + err.Error()
	} else {
		results["load_schema"] = "Success"
		results["schema_type"] = schema.Type
		results["properties_count"] = len(schema.Properties)
	}

	// Test 2: Get schema info
	if info, err := s.SchemaManager.GetSchemaInfo(); err != nil {
		results["schema_info"] = "Failed: " + err.Error()
	} else {
		results["schema_info"] = "Success"
		results["info"] = info
	}

	// Test 3: Generate form from schema
	if fields, err := s.SchemaManager.GenerateFormFromSchema(); err != nil {
		results["generate_form"] = "Failed: " + err.Error()
	} else {
		results["generate_form"] = "Success"
		results["form_fields_count"] = len(fields)
	}

	// Test 4: Export schema
	if data, err := s.SchemaManager.ExportSchema(); err != nil {
		results["export_schema"] = "Failed: " + err.Error()
	} else {
		results["export_schema"] = "Success"
		results["export_size"] = len(data)
	}

	// Test 5: Validate current content against schema
	if content, err := s.ContentManager.LoadContent(); err != nil {
		results["validate_content"] = "Failed to load content: " + err.Error()
	} else {
		if err := s.SchemaManager.ValidateAgainstSchema(map[string]interface{}{
			"title":       content.Title,
			"description": content.Description,
			"sections":    content.Sections,
		}); err != nil {
			results["validate_content"] = "Validation failed: " + err.Error()
		} else {
			results["validate_content"] = "Success"
		}
	}

	// Test 6: Backup schema
	if err := s.SchemaManager.BackupSchema(); err != nil {
		results["backup_schema"] = "Failed: " + err.Error()
	} else {
		results["backup_schema"] = "Success"
	}

	response := NewAPIResponse(true, "Schema test completed")
	response.SetData(results)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSchemaAnalyze returns comprehensive schema analysis
func (s *Server) handleSchemaAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	analysis, err := s.SchemaManager.ParseSchemaDetailed()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to analyze schema: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analysis)
}

// handleSchemaFieldMetadata returns metadata for a specific field
func (s *Server) handleSchemaFieldMetadata(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fieldName := r.URL.Query().Get("field")
	if fieldName == "" {
		http.Error(w, "Field name is required", http.StatusBadRequest)
		return
	}

	metadata, err := s.SchemaManager.GetFieldMetadata(fieldName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get field metadata: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metadata)
}

// handleSchemaValidationRules returns all validation rules for the schema
func (s *Server) handleSchemaValidationRules(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rules, err := s.SchemaManager.GetValidationRules()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get validation rules: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rules": rules,
		"count": len(rules),
	})
}

// handleSchemaFieldTypes returns field types mapping
func (s *Server) handleSchemaFieldTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fieldTypes, err := s.SchemaManager.GetSchemaFieldTypes()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get field types: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"field_types": fieldTypes,
		"count":       len(fieldTypes),
	})
}

// handleSchemaRequiredFields returns required and optional fields
func (s *Server) handleSchemaRequiredFields(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requiredFields, err := s.SchemaManager.GetRequiredFields()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get required fields: %v", err), http.StatusInternalServerError)
		return
	}

	optionalFields, err := s.SchemaManager.GetOptionalFields()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get optional fields: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"required": requiredFields,
		"optional": optionalFields,
		"total":    len(requiredFields) + len(optionalFields),
	})
}

// handleSchemaValidateField validates a field value against schema
func (s *Server) handleSchemaValidateField(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		FieldName string      `json:"field_name"`
		Value     interface{} `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if requestData.FieldName == "" {
		http.Error(w, "Field name is required", http.StatusBadRequest)
		return
	}

	validationFailures, err := s.SchemaManager.ValidateFieldValue(requestData.FieldName, requestData.Value)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to validate field: %v", err), http.StatusInternalServerError)
		return
	}

	isValid := len(validationFailures) == 0

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":    isValid,
		"failures": validationFailures,
		"field":    requestData.FieldName,
		"value":    requestData.Value,
	})
}

// handleSchemaValidateContent validates entire content using comprehensive validator
func (s *Server) handleSchemaValidateContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Content interface{} `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	validationResult, err := s.SchemaManager.ValidateContentDetailed(requestData.Content)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to validate content: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(validationResult)
}

// handleSchemaValidateFieldDetailed validates a field value using comprehensive validator
func (s *Server) handleSchemaValidateFieldDetailed(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		FieldName string      `json:"field_name"`
		Value     interface{} `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if requestData.FieldName == "" {
		http.Error(w, "Field name is required", http.StatusBadRequest)
		return
	}

	validationResult, err := s.SchemaManager.ValidateFieldValueDetailed(requestData.FieldName, requestData.Value)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to validate field: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(validationResult)
}

// handleSchemaValidationReport generates a comprehensive validation report
func (s *Server) handleSchemaValidationReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Content interface{} `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	report, err := s.SchemaManager.GenerateValidationReport(requestData.Content)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate validation report: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
