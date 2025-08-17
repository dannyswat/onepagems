package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Schema Parser Handlers

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

// Schema Validator Handlers

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
