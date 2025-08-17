package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"onepagems/internal/types"
)

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

	response := types.NewAPIResponse(true, "Template loaded successfully")
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

	response := types.NewAPIResponse(true, "Template saved successfully")
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

	response := types.NewAPIResponse(true, "Template info retrieved successfully")
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

	response := types.NewAPIResponse(true, "Template restored from backup successfully")
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

	response := types.NewAPIResponse(true, "Template test completed")
	response.SetData(results)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
