package server

import (
	"encoding/json"
	"net/http"

	"onepagems/internal/types"
)

// handleSchema handles schema management requests
func (s *Server) handleSchema(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Load and return current schema
		schema, err := s.SchemaManager.LoadSchema()
		if err != nil {
			response := types.NewAPIResponse(false, "Failed to load schema: "+err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := types.NewAPIResponse(true, "Schema loaded successfully")
		response.SetData(schema)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	case "POST":
		// Update schema
		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			response := types.NewAPIResponse(false, "Invalid JSON in request body: "+err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if err := s.SchemaManager.UpdateSchema(updates); err != nil {
			response := types.NewAPIResponse(false, "Failed to update schema: "+err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := types.NewAPIResponse(true, "Schema updated successfully")
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
		response := types.NewAPIResponse(false, "Failed to get schema information: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.NewAPIResponse(true, "Schema information retrieved")
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
		response := types.NewAPIResponse(false, "Failed to restore schema: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.NewAPIResponse(true, "Schema restored from backup successfully")
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
		response := types.NewAPIResponse(false, "Failed to export schema: "+err.Error())
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
		response := types.NewAPIResponse(false, "Invalid JSON in request body: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := s.SchemaManager.ImportSchema(requestData.Schema); err != nil {
		response := types.NewAPIResponse(false, "Failed to import schema: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.NewAPIResponse(true, "Schema imported successfully")
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
		response := types.NewAPIResponse(false, "Invalid JSON in request body: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := s.SchemaManager.ValidateAgainstSchema(requestData.Data); err != nil {
		response := types.NewAPIResponse(false, "Validation failed: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.NewAPIResponse(true, "Data is valid against schema")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSchemaForm generates complete form structure from schema
func (s *Server) handleSchemaForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	form, err := s.SchemaManager.GenerateCompleteForm()
	if err != nil {
		response := types.NewAPIResponse(false, "Failed to generate form from schema: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.NewAPIResponse(true, "Form generated from schema")
	response.SetData(form)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSchemaFormFields generates just the form fields array from schema
func (s *Server) handleSchemaFormFields(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fields, err := s.SchemaManager.GenerateFormFromSchema()
	if err != nil {
		response := types.NewAPIResponse(false, "Failed to generate form fields from schema: "+err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.NewAPIResponse(true, "Form fields generated from schema")
	response.SetData(map[string]interface{}{
		"fields": fields,
		"count":  len(fields),
	})
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

	response := types.NewAPIResponse(true, "Schema test completed")
	response.SetData(results)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
