package managers

import (
	"fmt"
	"sort"
	"strings"

	"onepagems/internal/types"
)

// FormGenerator handles dynamic form generation from JSON schema
type FormGenerator struct {
	schema      *types.SchemaData
	parser      *SchemaParser
	validator   *SchemaValidator
	imageFields []string // tracks fields that should be image pickers
}

// NewFormGenerator creates a new form generator
func NewFormGenerator(schema *types.SchemaData) *FormGenerator {
	parser := NewSchemaParser(schema)
	validator := NewSchemaValidator(schema)

	return &FormGenerator{
		schema:      schema,
		parser:      parser,
		validator:   validator,
		imageFields: make([]string, 0),
	}
}

// GenerateForm generates a complete form from the JSON schema
func (fg *FormGenerator) GenerateForm() (*types.GeneratedForm, error) {
	if fg.schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}

	fields, err := fg.generateFormFields("", fg.schema.Properties, false)
	if err != nil {
		return nil, fmt.Errorf("failed to generate form fields: %w", err)
	}

	// Sort fields to ensure consistent ordering
	sort.Slice(fields, func(i, j int) bool {
		return fg.getFieldPriority(fields[i]) < fg.getFieldPriority(fields[j])
	})

	form := &types.GeneratedForm{
		Fields: fields,
		Action: "/admin/content",
		Method: "POST",
	}

	return form, nil
}

// generateFormFields recursively generates form fields from schema properties
func (fg *FormGenerator) generateFormFields(prefix string, properties map[string]interface{}, isNested bool) ([]types.FormField, error) {
	var fields []types.FormField

	for fieldName, propData := range properties {
		propMap, ok := propData.(map[string]interface{})
		if !ok {
			continue
		}

		fullFieldName := fieldName
		if prefix != "" {
			fullFieldName = prefix + "." + fieldName
		}

		field, err := fg.createFormField(fullFieldName, fieldName, propMap, isNested)
		if err != nil {
			return nil, fmt.Errorf("failed to create field %s: %w", fieldName, err)
		}

		fields = append(fields, field)

		// Handle nested objects
		if field.Type == "object" {
			if nestedProps, ok := propMap["properties"].(map[string]interface{}); ok {
				nestedFields, err := fg.generateFormFields(fullFieldName, nestedProps, true)
				if err != nil {
					return nil, fmt.Errorf("failed to generate nested fields for %s: %w", fieldName, err)
				}
				fields = append(fields, nestedFields...)
			}
		}
	}

	return fields, nil
}

// createFormField creates a single form field from a schema property
func (fg *FormGenerator) createFormField(fullName, displayName string, prop map[string]interface{}, isNested bool) (types.FormField, error) {
	field := types.FormField{
		Name: fullName,
		Type: "text", // default
	}

	// Extract basic properties
	fg.extractBasicProperties(&field, displayName, prop)
	fg.extractTypeAndFormat(&field, prop)
	fg.extractValidationConstraints(&field, prop)
	fg.extractEnumOptions(&field, prop)

	// Handle special cases
	fg.handleSpecialFieldTypes(&field, prop)

	// Apply nested field styling
	if isNested {
		field.Label = fg.formatNestedLabel(displayName)
	}

	return field, nil
}

// extractBasicProperties extracts title, description, and basic metadata
func (fg *FormGenerator) extractBasicProperties(field *types.FormField, displayName string, prop map[string]interface{}) {
	// Extract label (from title or use display name)
	if title, ok := prop["title"].(string); ok {
		field.Label = title
	} else {
		field.Label = fg.formatFieldLabel(displayName)
	}

	// Extract description
	if description, ok := prop["description"].(string); ok {
		field.Description = description
		if field.Placeholder == "" {
			field.Placeholder = description
		}
	}

	// Extract default value
	if defaultValue, ok := prop["default"]; ok {
		field.Value = defaultValue
	}
}

// extractTypeAndFormat determines the field type based on schema type and format
func (fg *FormGenerator) extractTypeAndFormat(field *types.FormField, prop map[string]interface{}) {
	fieldType, _ := prop["type"].(string)
	format, _ := prop["format"].(string)

	switch fieldType {
	case "string":
		fg.handleStringField(field, format, prop)
	case "number", "integer":
		field.Type = "number"
		if fieldType == "integer" {
			field.Type = "number"
			field.Format = "integer"
		}
	case "boolean":
		field.Type = "checkbox"
	case "array":
		field.Type = "array"
		fg.handleArrayField(field, prop)
	case "object":
		field.Type = "object"
	default:
		field.Type = "text"
	}
}

// handleStringField handles string type fields with various formats
func (fg *FormGenerator) handleStringField(field *types.FormField, format string, prop map[string]interface{}) {
	switch format {
	case "email":
		field.Type = "email"
	case "password":
		field.Type = "password"
	case "textarea":
		field.Type = "textarea"
	case "url":
		field.Type = "url"
	case "tel":
		field.Type = "tel"
	case "date":
		field.Type = "date"
	case "datetime-local":
		field.Type = "datetime-local"
	case "time":
		field.Type = "time"
	case "color":
		field.Type = "color"
	case "image":
		field.Type = "image"
		fg.imageFields = append(fg.imageFields, field.Name)
	default:
		// Check for textarea hint in description or title
		if fg.isTextAreaField(prop) {
			field.Type = "textarea"
		} else {
			field.Type = "text"
		}
	}

	field.Format = format
}

// handleArrayField handles array type fields
func (fg *FormGenerator) handleArrayField(field *types.FormField, prop map[string]interface{}) {
	field.Type = "array"

	// Extract array item type
	if items, ok := prop["items"].(map[string]interface{}); ok {
		if itemType, ok := items["type"].(string); ok {
			field.Format = itemType
		}

		// For string arrays with enum, convert to multi-select
		if itemType, ok := items["type"].(string); ok && itemType == "string" {
			if enum, ok := items["enum"].([]interface{}); ok {
				field.Type = "multiselect"
				field.Options = fg.convertEnumToOptions(enum)
			}
		}
	}
}

// extractValidationConstraints extracts validation rules and converts them to form constraints
func (fg *FormGenerator) extractValidationConstraints(field *types.FormField, prop map[string]interface{}) {
	// Required field (this might be set at the parent level, but we can also check here)
	if required, ok := prop["required"].(bool); ok {
		field.Required = required
	}

	// String length constraints
	if minLength, ok := prop["minLength"]; ok {
		if minLen, ok := minLength.(float64); ok {
			field.Required = field.Required || minLen > 0
		}
	}

	// Number constraints
	if field.Type == "number" {
		if minimum, ok := prop["minimum"].(float64); ok {
			if field.Placeholder == "" {
				field.Placeholder = fmt.Sprintf("Minimum: %.0f", minimum)
			}
		}
		if maximum, ok := prop["maximum"].(float64); ok {
			if field.Placeholder == "" {
				field.Placeholder = fmt.Sprintf("Maximum: %.0f", maximum)
			} else {
				field.Placeholder += fmt.Sprintf(", Maximum: %.0f", maximum)
			}
		}
	}
}

// extractEnumOptions extracts enum values and converts them to select options
func (fg *FormGenerator) extractEnumOptions(field *types.FormField, prop map[string]interface{}) {
	if enum, ok := prop["enum"].([]interface{}); ok {
		field.Type = "select"
		field.Options = fg.convertEnumToOptions(enum)
	}
}

// handleSpecialFieldTypes handles special field types and formats
func (fg *FormGenerator) handleSpecialFieldTypes(field *types.FormField, prop map[string]interface{}) {
	// Check for image field hints in title or description
	if fg.isImageField(prop) && field.Type == "text" {
		field.Type = "image"
		fg.imageFields = append(fg.imageFields, field.Name)
	}

	// Check for rich text editor hint
	if fg.isRichTextField(prop) {
		field.Type = "richtext"
	}
}

// Helper methods

// isTextAreaField determines if a field should be a textarea based on hints
func (fg *FormGenerator) isTextAreaField(prop map[string]interface{}) bool {
	title, _ := prop["title"].(string)
	description, _ := prop["description"].(string)

	textAreaKeywords := []string{"content", "description", "text", "message", "body", "summary"}

	for _, keyword := range textAreaKeywords {
		if strings.Contains(strings.ToLower(title), keyword) ||
			strings.Contains(strings.ToLower(description), keyword) {
			return true
		}
	}

	// Check if maxLength suggests a textarea (longer text)
	if maxLength, ok := prop["maxLength"].(float64); ok && maxLength > 100 {
		return true
	}

	return false
}

// isImageField determines if a field should be an image picker
func (fg *FormGenerator) isImageField(prop map[string]interface{}) bool {
	title, _ := prop["title"].(string)
	description, _ := prop["description"].(string)

	imageKeywords := []string{"image", "photo", "picture", "avatar", "logo", "banner", "background"}

	for _, keyword := range imageKeywords {
		if strings.Contains(strings.ToLower(title), keyword) ||
			strings.Contains(strings.ToLower(description), keyword) {
			return true
		}
	}

	return false
}

// isRichTextField determines if a field should use a rich text editor
func (fg *FormGenerator) isRichTextField(prop map[string]interface{}) bool {
	title, _ := prop["title"].(string)
	description, _ := prop["description"].(string)
	format, _ := prop["format"].(string)

	if format == "html" || format == "richtext" {
		return true
	}

	richTextKeywords := []string{"html", "rich", "formatted", "wysiwyg"}

	for _, keyword := range richTextKeywords {
		if strings.Contains(strings.ToLower(title), keyword) ||
			strings.Contains(strings.ToLower(description), keyword) {
			return true
		}
	}

	return false
}

// convertEnumToOptions converts enum values to string options
func (fg *FormGenerator) convertEnumToOptions(enum []interface{}) []string {
	options := make([]string, len(enum))
	for i, val := range enum {
		options[i] = fmt.Sprintf("%v", val)
	}
	return options
}

// formatFieldLabel creates a human-readable label from a field name
func (fg *FormGenerator) formatFieldLabel(fieldName string) string {
	// Split by dots for nested fields
	parts := strings.Split(fieldName, ".")
	lastPart := parts[len(parts)-1]

	// Split by underscores and capitalize
	words := strings.Split(lastPart, "_")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}

// formatNestedLabel formats labels for nested fields
func (fg *FormGenerator) formatNestedLabel(fieldName string) string {
	label := fg.formatFieldLabel(fieldName)
	return "  " + label // Indent nested fields
}

// getFieldPriority determines the display order priority of fields
func (fg *FormGenerator) getFieldPriority(field types.FormField) int {
	// Primary fields first
	if field.Name == "title" {
		return 1
	}
	if field.Name == "description" {
		return 2
	}

	// Top-level fields before nested fields
	if !strings.Contains(field.Name, ".") {
		return 10
	}

	// Nested fields ordered by depth and name
	depth := strings.Count(field.Name, ".")
	return 100 + (depth * 10)
}

// GetImageFields returns a list of fields that should use image pickers
func (fg *FormGenerator) GetImageFields() []string {
	return fg.imageFields
}

// GenerateFieldValidationRules generates validation rules for a specific field
func (fg *FormGenerator) GenerateFieldValidationRules(fieldName string) ([]ValidationRule, error) {
	metadata, err := fg.parser.GetFieldMetadata(fieldName)
	if err != nil {
		return nil, fmt.Errorf("failed to get field metadata: %w", err)
	}

	var rules []ValidationRule

	// Convert parsed metadata to validation rules
	if metadata.Required {
		rules = append(rules, ValidationRule{
			Type:         "required",
			Message:      fmt.Sprintf("Field '%s' is required", fieldName),
			PropertyPath: fieldName,
		})
	}

	if metadata.Type == "string" {
		if *metadata.MinLength > 0 {
			rules = append(rules, ValidationRule{
				Type:         "minLength",
				Value:        *metadata.MinLength,
				Message:      fmt.Sprintf("Field '%s' must be at least %d characters", fieldName, *metadata.MinLength),
				PropertyPath: fieldName,
			})
		}
		if *metadata.MaxLength > 0 {
			rules = append(rules, ValidationRule{
				Type:         "maxLength",
				Value:        *metadata.MaxLength,
				Message:      fmt.Sprintf("Field '%s' must be at most %d characters", fieldName, *metadata.MaxLength),
				PropertyPath: fieldName,
			})
		}
	}

	return rules, nil
}

// ValidateFormField validates a form field value using the schema validator
func (fg *FormGenerator) ValidateFormField(fieldName string, value interface{}) *ValidationResult {
	return fg.validator.ValidateFieldValue(fieldName, value)
}
