package server

import (
	"log"
	"net/http"
	"path/filepath"
)

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
