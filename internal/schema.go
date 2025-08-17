package internal

import (
	"encoding/json"
	"fmt"
)

// SchemaManager handles schema.json operations
type SchemaManager struct {
	storage *FileStorage
	dataDir string
}

// NewSchemaManager creates a new schema manager
func NewSchemaManager(storage *FileStorage, dataDir string) *SchemaManager {
	return &SchemaManager{
		storage: storage,
		dataDir: dataDir,
	}
}

// schemaFilePath returns the filename for schema.json
func (sm *SchemaManager) schemaFilePath() string {
	return "schema.json"
}

// LoadSchema loads schema from schema.json or creates default if not exists
func (sm *SchemaManager) LoadSchema() (*SchemaData, error) {
	schemaFilename := sm.schemaFilePath()

	// Check if schema.json exists
	if !sm.storage.FileExists(schemaFilename) {
		// Create default schema
		defaultSchema := sm.createDefaultSchema()
		if err := sm.SaveSchema(defaultSchema); err != nil {
			return nil, fmt.Errorf("failed to create default schema: %w", err)
		}
		return defaultSchema, nil
	}

	// Load existing schema
	var schema SchemaData
	if err := sm.storage.ReadJSONFile(schemaFilename, &schema); err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	// Validate schema structure
	if err := sm.validateSchema(&schema); err != nil {
		return nil, fmt.Errorf("schema validation failed: %w", err)
	}

	return &schema, nil
}

// SaveSchema saves schema to schema.json with backup
func (sm *SchemaManager) SaveSchema(schema *SchemaData) error {
	if schema == nil {
		return fmt.Errorf("schema cannot be nil")
	}

	// Validate schema before saving
	if err := sm.validateSchema(schema); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	// Save with backup
	schemaFilename := sm.schemaFilePath()
	if err := sm.storage.WriteJSONFile(schemaFilename, schema); err != nil {
		return fmt.Errorf("failed to save schema file: %w", err)
	}

	return nil
}

// UpdateSchema updates the schema with new properties
func (sm *SchemaManager) UpdateSchema(updates map[string]interface{}) error {
	// Load current schema
	schema, err := sm.LoadSchema()
	if err != nil {
		return fmt.Errorf("failed to load current schema: %w", err)
	}

	// Apply updates
	for key, value := range updates {
		switch key {
		case "properties":
			if properties, ok := value.(map[string]interface{}); ok {
				schema.Properties = properties
			} else {
				return fmt.Errorf("properties must be a map")
			}
		case "type":
			if schemaType, ok := value.(string); ok {
				schema.Type = schemaType
			} else {
				return fmt.Errorf("type must be a string")
			}
		case "$schema":
			if schemaVersion, ok := value.(string); ok {
				schema.Schema = schemaVersion
			} else {
				return fmt.Errorf("$schema must be a string")
			}
		default:
			return fmt.Errorf("unknown schema field: %s", key)
		}
	}

	// Save updated schema
	return sm.SaveSchema(schema)
}

// BackupSchema creates a backup of the current schema
func (sm *SchemaManager) BackupSchema() error {
	schemaFilename := sm.schemaFilePath()
	return sm.storage.CreateBackup(schemaFilename)
}

// RestoreSchema restores schema from backup
func (sm *SchemaManager) RestoreSchema() error {
	schemaFilename := sm.schemaFilePath()
	return sm.storage.RestoreFromBackup(schemaFilename)
}

// GetSchemaInfo returns information about the current schema
func (sm *SchemaManager) GetSchemaInfo() (map[string]interface{}, error) {
	schema, err := sm.LoadSchema()
	if err != nil {
		return nil, err
	}

	info := map[string]interface{}{
		"schema_version": schema.Schema,
		"type":           schema.Type,
		"properties":     len(schema.Properties),
	}

	// Add property names
	propertyNames := make([]string, 0, len(schema.Properties))
	for name := range schema.Properties {
		propertyNames = append(propertyNames, name)
	}
	info["property_names"] = propertyNames

	// Analyze property types
	propertyTypes := make(map[string]string)
	for name, prop := range schema.Properties {
		if propMap, ok := prop.(map[string]interface{}); ok {
			if propType, ok := propMap["type"].(string); ok {
				propertyTypes[name] = propType
			}
		}
	}
	info["property_types"] = propertyTypes

	return info, nil
}

// ValidateAgainstSchema validates data against the current schema
func (sm *SchemaManager) ValidateAgainstSchema(data interface{}) error {
	schema, err := sm.LoadSchema()
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	// Basic validation - check if data is an object when schema type is object
	if schema.Type == "object" {
		if _, ok := data.(map[string]interface{}); !ok {
			return fmt.Errorf("data must be an object according to schema")
		}

		dataMap := data.(map[string]interface{})

		// Check required properties (basic implementation)
		for propName, propSchema := range schema.Properties {
			if propMap, ok := propSchema.(map[string]interface{}); ok {
				if required, ok := propMap["required"].(bool); ok && required {
					if _, exists := dataMap[propName]; !exists {
						return fmt.Errorf("required property '%s' is missing", propName)
					}
				}
			}
		}
	}

	return nil
}

// createDefaultSchema creates a default JSON schema structure for content
func (sm *SchemaManager) createDefaultSchema() *SchemaData {
	return &SchemaData{
		Schema: "https://json-schema.org/draft/2020-12/schema",
		Type:   "object",
		Properties: map[string]interface{}{
			"title": map[string]interface{}{
				"type":        "string",
				"title":       "Page Title",
				"description": "The main title of your website",
				"minLength":   1,
				"maxLength":   100,
			},
			"description": map[string]interface{}{
				"type":        "string",
				"title":       "Page Description",
				"description": "A brief description of your website",
				"maxLength":   500,
			},
			"sections": map[string]interface{}{
				"type":        "object",
				"title":       "Content Sections",
				"description": "Various content sections of your website",
				"properties": map[string]interface{}{
					"hero": map[string]interface{}{
						"type":        "object",
						"title":       "Hero Section",
						"description": "The main hero/banner section",
						"properties": map[string]interface{}{
							"title": map[string]interface{}{
								"type":        "string",
								"title":       "Hero Title",
								"description": "Main headline for the hero section",
							},
							"subtitle": map[string]interface{}{
								"type":        "string",
								"title":       "Hero Subtitle",
								"description": "Subtitle or tagline for the hero section",
							},
							"content": map[string]interface{}{
								"type":        "string",
								"title":       "Hero Content",
								"description": "Main content text for the hero section",
								"format":      "textarea",
							},
						},
					},
					"about": map[string]interface{}{
						"type":        "object",
						"title":       "About Section",
						"description": "Information about you or your organization",
						"properties": map[string]interface{}{
							"title": map[string]interface{}{
								"type":        "string",
								"title":       "About Title",
								"description": "Title for the about section",
							},
							"content": map[string]interface{}{
								"type":        "string",
								"title":       "About Content",
								"description": "Main content for the about section",
								"format":      "textarea",
							},
						},
					},
					"contact": map[string]interface{}{
						"type":        "object",
						"title":       "Contact Section",
						"description": "Contact information and details",
						"properties": map[string]interface{}{
							"title": map[string]interface{}{
								"type":        "string",
								"title":       "Contact Title",
								"description": "Title for the contact section",
							},
							"email": map[string]interface{}{
								"type":        "string",
								"title":       "Email Address",
								"description": "Contact email address",
								"format":      "email",
							},
							"phone": map[string]interface{}{
								"type":        "string",
								"title":       "Phone Number",
								"description": "Contact phone number",
							},
							"address": map[string]interface{}{
								"type":        "string",
								"title":       "Address",
								"description": "Physical address",
								"format":      "textarea",
							},
						},
					},
				},
			},
		},
	}
}

// validateSchema validates the schema structure
func (sm *SchemaManager) validateSchema(schema *SchemaData) error {
	if schema == nil {
		return fmt.Errorf("schema cannot be nil")
	}

	if schema.Schema == "" {
		schema.Schema = "https://json-schema.org/draft/2020-12/schema"
	}

	if schema.Type == "" {
		schema.Type = "object"
	}

	if schema.Properties == nil {
		schema.Properties = make(map[string]interface{})
	}

	return nil
}

// ExportSchema exports schema as JSON for external use
func (sm *SchemaManager) ExportSchema() ([]byte, error) {
	schema, err := sm.LoadSchema()
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(schema, "", "  ")
}

// ImportSchema imports schema from JSON data
func (sm *SchemaManager) ImportSchema(data []byte) error {
	var schema SchemaData
	if err := json.Unmarshal(data, &schema); err != nil {
		return fmt.Errorf("failed to parse imported schema: %w", err)
	}

	return sm.SaveSchema(&schema)
}

// GenerateFormFromSchema generates form field definitions from the schema
func (sm *SchemaManager) GenerateFormFromSchema() ([]FormField, error) {
	schema, err := sm.LoadSchema()
	if err != nil {
		return nil, err
	}

	var fields []FormField

	// Generate fields from schema properties
	for propName, propData := range schema.Properties {
		if propMap, ok := propData.(map[string]interface{}); ok {
			field := sm.createFormFieldFromProperty(propName, propMap)
			fields = append(fields, field)
		}
	}

	return fields, nil
}

// createFormFieldFromProperty creates a form field from a schema property
func (sm *SchemaManager) createFormFieldFromProperty(name string, prop map[string]interface{}) FormField {
	field := FormField{
		Name: name,
		Type: "text", // default
	}

	// Extract type
	if fieldType, ok := prop["type"].(string); ok {
		switch fieldType {
		case "string":
			field.Type = "text"
			if format, ok := prop["format"].(string); ok {
				switch format {
				case "email":
					field.Type = "email"
				case "textarea":
					field.Type = "textarea"
				}
			}
		case "number", "integer":
			field.Type = "number"
		case "boolean":
			field.Type = "checkbox"
		case "object":
			field.Type = "object"
		}
	}

	// Extract label (from title or use name)
	if title, ok := prop["title"].(string); ok {
		field.Label = title
	} else {
		field.Label = name
	}

	// Extract description
	if description, ok := prop["description"].(string); ok {
		field.Description = description
		field.Placeholder = description
	}

	// Extract required
	if required, ok := prop["required"].(bool); ok {
		field.Required = required
	}

	return field
}

// ParseSchemaDetailed returns comprehensive schema analysis using the schema parser
func (sm *SchemaManager) ParseSchemaDetailed() (*SchemaAnalysis, error) {
	schema, err := sm.LoadSchema()
	if err != nil {
		return nil, err
	}

	parser := NewSchemaParser(schema)
	return parser.ParseSchema()
}

// GetFieldMetadata returns detailed metadata for a specific field
func (sm *SchemaManager) GetFieldMetadata(fieldName string) (*ParsedProperty, error) {
	schema, err := sm.LoadSchema()
	if err != nil {
		return nil, err
	}

	parser := NewSchemaParser(schema)
	return parser.GetFieldMetadata(fieldName)
}

// GetValidationRules returns all validation rules for the schema
func (sm *SchemaManager) GetValidationRules() ([]ValidationRule, error) {
	analysis, err := sm.ParseSchemaDetailed()
	if err != nil {
		return nil, err
	}

	return analysis.ValidationRules, nil
}

// ValidateFieldValue validates a single field value against the schema
func (sm *SchemaManager) ValidateFieldValue(fieldName string, value interface{}) ([]ValidationRule, error) {
	schema, err := sm.LoadSchema()
	if err != nil {
		return nil, err
	}

	parser := NewSchemaParser(schema)
	failures := parser.ValidateFieldValue(fieldName, value)
	return failures, nil
}

// GetSchemaFieldTypes returns a map of field names to their types
func (sm *SchemaManager) GetSchemaFieldTypes() (map[string]string, error) {
	analysis, err := sm.ParseSchemaDetailed()
	if err != nil {
		return nil, err
	}

	return analysis.PropertyTypes, nil
}

// GetRequiredFields returns list of required field names
func (sm *SchemaManager) GetRequiredFields() ([]string, error) {
	analysis, err := sm.ParseSchemaDetailed()
	if err != nil {
		return nil, err
	}

	return analysis.RequiredFields, nil
}

// GetOptionalFields returns list of optional field names
func (sm *SchemaManager) GetOptionalFields() ([]string, error) {
	analysis, err := sm.ParseSchemaDetailed()
	if err != nil {
		return nil, err
	}

	return analysis.OptionalFields, nil
}

// GetNestedObjects returns list of fields that are nested objects
func (sm *SchemaManager) GetNestedObjects() ([]string, error) {
	analysis, err := sm.ParseSchemaDetailed()
	if err != nil {
		return nil, err
	}

	return analysis.NestedObjects, nil
}

// GetArrayFields returns list of fields that are arrays
func (sm *SchemaManager) GetArrayFields() ([]string, error) {
	analysis, err := sm.ParseSchemaDetailed()
	if err != nil {
		return nil, err
	}

	return analysis.Arrays, nil
}

// GetEnumFields returns map of field names to their enum values
func (sm *SchemaManager) GetEnumFields() (map[string][]interface{}, error) {
	analysis, err := sm.ParseSchemaDetailed()
	if err != nil {
		return nil, err
	}

	return analysis.EnumFields, nil
}

// ValidateContentDetailed validates content using the comprehensive schema validator
func (sm *SchemaManager) ValidateContentDetailed(content interface{}) (*ValidationResult, error) {
	schema, err := sm.LoadSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to load schema: %w", err)
	}

	validator := NewSchemaValidator(schema)
	result := validator.ValidateContent(content)
	return result, nil
}

// ValidateFieldValueDetailed validates a single field value using the comprehensive validator
func (sm *SchemaManager) ValidateFieldValueDetailed(fieldName string, value interface{}) (*ValidationResult, error) {
	schema, err := sm.LoadSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to load schema: %w", err)
	}

	validator := NewSchemaValidator(schema)
	result := validator.ValidateFieldValue(fieldName, value)
	return result, nil
}

// GenerateValidationReport generates a detailed validation report for content
func (sm *SchemaManager) GenerateValidationReport(content interface{}) (map[string]interface{}, error) {
	schema, err := sm.LoadSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to load schema: %w", err)
	}

	validator := NewSchemaValidator(schema)
	report := validator.GenerateValidationReport(content)
	return report, nil
}
