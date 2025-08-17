package server

import (
	"encoding/json"
	"net/http"
	"time"

	"onepagems/internal/types"
)

// handleContent handles content management requests
func (s *Server) handleContent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Load and return current content
		content, err := s.ContentManager.LoadContent()
		if err != nil {
			response := types.NewAPIResponse(false, "Failed to load content: "+err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := types.NewAPIResponse(true, "Content loaded successfully")
		response.SetData(content)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	case "POST":
		// Update content
		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			response := types.NewAPIResponse(false, "Invalid JSON in request body: "+err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if err := s.ContentManager.UpdateContent(updates); err != nil {
			response := types.NewAPIResponse(false, "Failed to update content: "+err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := types.NewAPIResponse(true, "Content updated successfully")
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
		response := types.NewAPIResponse(false, "Failed to get content information: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.NewAPIResponse(true, "Content information retrieved")
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
		response := types.NewAPIResponse(false, "Failed to restore content: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.NewAPIResponse(true, "Content restored from backup successfully")
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
		response := types.NewAPIResponse(false, "Failed to export content: "+err.Error())
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
		response := types.NewAPIResponse(false, "Invalid JSON in request body: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := s.ContentManager.ImportContent(requestData.Content); err != nil {
		response := types.NewAPIResponse(false, "Failed to import content: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.NewAPIResponse(true, "Content imported successfully")
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

	response := types.NewAPIResponse(true, "Content test completed")
	response.SetData(results)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
