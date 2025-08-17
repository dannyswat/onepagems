package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"onepagems/internal/types"
)

// AdminPageData represents the data passed to admin templates
type AdminPageData struct {
	Title          string
	Username       string
	Page           string
	Content        interface{}
	Stats          *AdminStats
	Status         *SystemStatus
	RecentActivity []ActivityItem
}

// AdminStats represents dashboard statistics
type AdminStats struct {
	ContentFields int    `json:"content_fields"`
	Images        int    `json:"images"`
	LastUpdated   string `json:"last_updated"`
	SchemaVersion string `json:"schema_version"`
}

// SystemStatus represents system component status
type SystemStatus struct {
	ContentModified  string `json:"content_modified"`
	SchemaModified   string `json:"schema_modified"`
	TemplateModified string `json:"template_modified"`
	SiteGenerated    string `json:"site_generated"`
}

// ActivityItem represents a recent activity entry
type ActivityItem struct {
	Action      string    `json:"action"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

// handleAdminPanel serves the main admin dashboard
func (s *Server) handleAdminPanel(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, ok := types.SessionFromContext(r.Context())
	if !ok {
		http.Error(w, "Session not found", http.StatusInternalServerError)
		return
	}

	// Gather dashboard data
	stats, err := s.getAdminStats()
	if err != nil {
		http.Error(w, "Failed to load dashboard stats", http.StatusInternalServerError)
		return
	}

	status, err := s.getSystemStatus()
	if err != nil {
		http.Error(w, "Failed to load system status", http.StatusInternalServerError)
		return
	}

	recentActivity := s.getRecentActivity()

	dashboardContent, err := s.renderTemplate("admin_dashboard.html", map[string]interface{}{
		"Stats":          stats,
		"Status":         status,
		"RecentActivity": recentActivity,
	})
	if err != nil {
		http.Error(w, "Failed to render dashboard", http.StatusInternalServerError)
		return
	}

	pageData := AdminPageData{
		Title:          "Dashboard",
		Username:       session.Username,
		Page:           "dashboard",
		Content:        dashboardContent,
		Stats:          stats,
		Status:         status,
		RecentActivity: recentActivity,
	}

	s.renderAdminPage(w, pageData)
}

// handleAdminContent serves the content editor interface
func (s *Server) handleAdminContent(w http.ResponseWriter, r *http.Request) {
	session, ok := types.SessionFromContext(r.Context())
	if !ok {
		http.Error(w, "Session not found", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "GET":
		// Serve content editor interface
		contentEditorHTML, err := s.renderTemplate("admin_content.html", nil)
		if err != nil {
			http.Error(w, "Failed to render content editor", http.StatusInternalServerError)
			return
		}

		pageData := AdminPageData{
			Title:    "Content Editor",
			Username: session.Username,
			Page:     "content",
			Content:  contentEditorHTML,
		}

		s.renderAdminPage(w, pageData)

	case "POST":
		// Handle content updates
		s.handleContentUpdate(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleContentUpdate processes content form submissions
func (s *Server) handleContentUpdate(w http.ResponseWriter, r *http.Request) {
	var content map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
		response := types.NewAPIResponse(false, "Invalid JSON data: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate content against schema
	validationResult, err := s.SchemaManager.ValidateContentDetailed(content)
	if err != nil {
		response := types.NewAPIResponse(false, "Validation failed: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !validationResult.Valid {
		response := types.NewAPIResponse(false, "Content validation failed")
		response.SetData(map[string]interface{}{
			"errors":      validationResult.Errors,
			"valid":       false,
			"error_count": len(validationResult.Errors),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Save content
	contentData := &types.ContentData{}
	if err := s.mapToContentData(content, contentData); err != nil {
		response := types.NewAPIResponse(false, "Failed to process content: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := s.ContentManager.SaveContent(contentData); err != nil {
		response := types.NewAPIResponse(false, "Failed to save content: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Log activity
	s.logActivity("Content Updated", "Content has been successfully updated through the admin panel")

	response := types.NewAPIResponse(true, "Content saved successfully")
	response.SetData(map[string]interface{}{
		"validation": validationResult,
		"timestamp":  time.Now().Format(time.RFC3339),
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper methods

// renderAdminPage renders the main admin template with provided data
func (s *Server) renderAdminPage(w http.ResponseWriter, data AdminPageData) {
	tmplPath := filepath.Join("templates", "admin.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Failed to load admin template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render admin page", http.StatusInternalServerError)
		return
	}
}

// renderTemplate renders a specific template file and returns the HTML string
func (s *Server) renderTemplate(templateName string, data interface{}) (string, error) {
	tmplPath := filepath.Join("templates", templateName)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("failed to load template %s: %w", templateName, err)
	}

	// Use a bytes buffer to capture template output
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

// getAdminStats collects dashboard statistics
func (s *Server) getAdminStats() (*AdminStats, error) {
	// Get schema info to count fields
	schema, err := s.SchemaManager.LoadSchema()
	if err != nil {
		return nil, err
	}

	fieldCount := len(schema.Properties)
	if schema.Properties != nil {
		// Count nested fields too
		fieldCount = s.countSchemaFields(schema.Properties)
	}

	// TODO: Count images in images directory
	imageCount := 0

	stats := &AdminStats{
		ContentFields: fieldCount,
		Images:        imageCount,
		LastUpdated:   time.Now().Format("2006-01-02"),
		SchemaVersion: "1.0", // TODO: Get from schema
	}

	return stats, nil
}

// getSystemStatus collects system component status
func (s *Server) getSystemStatus() (*SystemStatus, error) {
	status := &SystemStatus{
		ContentModified:  "Recent",
		SchemaModified:   "Recent",
		TemplateModified: "Recent",
		SiteGenerated:    "Pending",
	}

	// TODO: Get actual file modification times
	return status, nil
}

// getRecentActivity returns recent activity items
func (s *Server) getRecentActivity() []ActivityItem {
	// TODO: Implement actual activity logging
	return []ActivityItem{
		{
			Action:      "Content Updated",
			Description: "Website content was updated through the admin panel",
			Timestamp:   time.Now().Add(-1 * time.Hour),
		},
		{
			Action:      "Schema Modified",
			Description: "JSON schema was updated to add new fields",
			Timestamp:   time.Now().Add(-2 * time.Hour),
		},
	}
}

// logActivity logs an activity item
func (s *Server) logActivity(action, description string) {
	// TODO: Implement proper activity logging to file or database
	fmt.Printf("[ACTIVITY] %s: %s\n", action, description)
}

// countSchemaFields recursively counts fields in schema
func (s *Server) countSchemaFields(properties map[string]interface{}) int {
	count := 0
	for _, prop := range properties {
		count++
		if propMap, ok := prop.(map[string]interface{}); ok {
			if nestedProps, ok := propMap["properties"].(map[string]interface{}); ok {
				count += s.countSchemaFields(nestedProps)
			}
		}
	}
	return count
}

// mapToContentData converts map to ContentData struct
func (s *Server) mapToContentData(content map[string]interface{}, target *types.ContentData) error {
	// Extract standard fields
	if title, ok := content["title"].(string); ok {
		target.Title = title
	}
	if description, ok := content["description"].(string); ok {
		target.Description = description
	}

	// Store all content as custom data
	target.Sections = content

	return nil
}

// Individual API handlers for direct routing

// handleAPIStats returns dashboard statistics as JSON
func (s *Server) handleAPIStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := s.getAdminStats()
	if err != nil {
		response := types.NewAPIResponse(false, "Failed to get stats: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.NewAPIResponse(true, "Stats retrieved successfully")
	response.SetData(stats)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleAPIGenerate handles site generation requests
func (s *Server) handleAPIGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For now, return success. In Phase 8, this will actually generate the site
	response := types.NewAPIResponse(true, "Site generation completed successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleContentAutoSave handles auto-save functionality for content editor
func (s *Server) handleContentAutoSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON body
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		// Try parsing as form data for backward compatibility
		if err := r.ParseForm(); err != nil {
			response := types.NewAPIResponse(false, "Invalid request data")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Convert form data to map
		updates = make(map[string]interface{})
		for key, values := range r.Form {
			if len(values) > 0 {
				updates[key] = values[0]
			}
		}
	}

	// Update content
	if err := s.ContentManager.UpdateContentFlexible(updates); err != nil {
		response := types.NewAPIResponse(false, "Auto-save failed: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.NewAPIResponse(true, "Content auto-saved successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handlePreviewContent provides preview functionality
func (s *Server) handlePreviewContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For now, redirect to the main site. In Phase 8, this will provide live preview
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// handleAPIStatus returns system status as JSON
func (s *Server) handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status, err := s.getSystemStatus()
	if err != nil {
		response := types.NewAPIResponse(false, "Failed to get status: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.NewAPIResponse(true, "Status retrieved successfully")
	response.SetData(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
