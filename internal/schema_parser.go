package internal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// SchemaParser handles parsing and analysis of JSON Schema definitions
type SchemaParser struct {
	schema *SchemaData
}

// NewSchemaParser creates a new schema parser
func NewSchemaParser(schema *SchemaData) *SchemaParser {
	return &SchemaParser{
		schema: schema,
	}
}

// ParsedProperty represents a parsed schema property with all metadata
type ParsedProperty struct {
	Name                 string                     `json:"name"`
	Type                 string                     `json:"type"`
	Format               string                     `json:"format,omitempty"`
	Title                string                     `json:"title,omitempty"`
	Description          string                     `json:"description,omitempty"`
	Required             bool                       `json:"required"`
	Default              interface{}                `json:"default,omitempty"`
	Enum                 []interface{}              `json:"enum,omitempty"`
	Pattern              string                     `json:"pattern,omitempty"`
	MinLength            *int                       `json:"minLength,omitempty"`
	MaxLength            *int                       `json:"maxLength,omitempty"`
	Minimum              *float64                   `json:"minimum,omitempty"`
	Maximum              *float64                   `json:"maximum,omitempty"`
	Items                *ParsedProperty            `json:"items,omitempty"`      // For arrays
	Properties           map[string]*ParsedProperty `json:"properties,omitempty"` // For objects
	AdditionalProperties bool                       `json:"additionalProperties"`
	Examples             []interface{}              `json:"examples,omitempty"`
	Raw                  map[string]interface{}     `json:"raw"` // Original property definition
}

// ValidationRule represents a single validation rule extracted from schema
type ValidationRule struct {
	Type         string      `json:"type"`          // required, minLength, pattern, etc.
	Value        interface{} `json:"value"`         // the validation value
	Message      string      `json:"message"`       // human-readable validation message
	PropertyPath string      `json:"property_path"` // dot-notation path to property
}

// SchemaAnalysis contains comprehensive analysis of the schema
type SchemaAnalysis struct {
	TotalProperties int                        `json:"total_properties"`
	RequiredFields  []string                   `json:"required_fields"`
	OptionalFields  []string                   `json:"optional_fields"`
	PropertyTypes   map[string]string          `json:"property_types"`
	ValidationRules []ValidationRule           `json:"validation_rules"`
	NestedObjects   []string                   `json:"nested_objects"`
	Arrays          []string                   `json:"arrays"`
	EnumFields      map[string][]interface{}   `json:"enum_fields"`
	FormattedFields map[string]string          `json:"formatted_fields"`
	Properties      map[string]*ParsedProperty `json:"properties"`
}

// ParseSchema parses the entire schema and returns detailed analysis
func (sp *SchemaParser) ParseSchema() (*SchemaAnalysis, error) {
	if sp.schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}

	analysis := &SchemaAnalysis{
		PropertyTypes:   make(map[string]string),
		ValidationRules: make([]ValidationRule, 0),
		NestedObjects:   make([]string, 0),
		Arrays:          make([]string, 0),
		EnumFields:      make(map[string][]interface{}),
		FormattedFields: make(map[string]string),
		Properties:      make(map[string]*ParsedProperty),
		RequiredFields:  make([]string, 0),
		OptionalFields:  make([]string, 0),
	}

	// Get required fields from root level
	requiredFields := sp.extractRequiredFields(sp.schema.Properties)

	// Parse each property
	for propName, propData := range sp.schema.Properties {
		propMap, ok := propData.(map[string]interface{})
		if !ok {
			continue
		}

		parsedProp, err := sp.parseProperty(propName, propMap, "", requiredFields)
		if err != nil {
			return nil, fmt.Errorf("failed to parse property '%s': %w", propName, err)
		}

		analysis.Properties[propName] = parsedProp
		analysis.TotalProperties++

		// Categorize property
		if parsedProp.Required {
			analysis.RequiredFields = append(analysis.RequiredFields, propName)
		} else {
			analysis.OptionalFields = append(analysis.OptionalFields, propName)
		}

		analysis.PropertyTypes[propName] = parsedProp.Type

		// Track special property types
		if parsedProp.Type == "object" {
			analysis.NestedObjects = append(analysis.NestedObjects, propName)
		}
		if parsedProp.Type == "array" {
			analysis.Arrays = append(analysis.Arrays, propName)
		}
		if len(parsedProp.Enum) > 0 {
			analysis.EnumFields[propName] = parsedProp.Enum
		}
		if parsedProp.Format != "" {
			analysis.FormattedFields[propName] = parsedProp.Format
		}

		// Extract validation rules
		rules := sp.extractValidationRules(propName, parsedProp, "")
		analysis.ValidationRules = append(analysis.ValidationRules, rules...)
	}

	return analysis, nil
}

// parseProperty parses a single property recursively
func (sp *SchemaParser) parseProperty(name string, prop map[string]interface{}, path string, requiredFields []string) (*ParsedProperty, error) {
	parsed := &ParsedProperty{
		Name:                 name,
		Type:                 "string", // default
		AdditionalProperties: true,     // default
		Raw:                  prop,
	}

	// Build property path
	if path != "" {
		parsed.Name = path + "." + name
	}

	// Check if required
	parsed.Required = sp.isRequired(name, requiredFields)

	// Extract basic properties
	if propType, ok := prop["type"].(string); ok {
		parsed.Type = propType
	}

	if format, ok := prop["format"].(string); ok {
		parsed.Format = format
	}

	if title, ok := prop["title"].(string); ok {
		parsed.Title = title
	}

	if description, ok := prop["description"].(string); ok {
		parsed.Description = description
	}

	if defaultVal, ok := prop["default"]; ok {
		parsed.Default = defaultVal
	}

	if pattern, ok := prop["pattern"].(string); ok {
		parsed.Pattern = pattern
	}

	// Extract numeric constraints
	if minLength, ok := prop["minLength"].(float64); ok {
		val := int(minLength)
		parsed.MinLength = &val
	}

	if maxLength, ok := prop["maxLength"].(float64); ok {
		val := int(maxLength)
		parsed.MaxLength = &val
	}

	if minimum, ok := prop["minimum"].(float64); ok {
		parsed.Minimum = &minimum
	}

	if maximum, ok := prop["maximum"].(float64); ok {
		parsed.Maximum = &maximum
	}

	// Extract enum values
	if enumData, ok := prop["enum"].([]interface{}); ok {
		parsed.Enum = enumData
	}

	// Extract examples
	if examples, ok := prop["examples"].([]interface{}); ok {
		parsed.Examples = examples
	}

	// Handle array type
	if parsed.Type == "array" {
		if itemsData, ok := prop["items"].(map[string]interface{}); ok {
			itemsProp, err := sp.parseProperty("items", itemsData, parsed.Name, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to parse array items: %w", err)
			}
			parsed.Items = itemsProp
		}
	}

	// Handle object type
	if parsed.Type == "object" {
		parsed.Properties = make(map[string]*ParsedProperty)

		if properties, ok := prop["properties"].(map[string]interface{}); ok {
			// Get required fields for this nested object
			nestedRequired := sp.extractRequiredFields(map[string]interface{}{"required": prop["required"]})

			for nestedName, nestedData := range properties {
				if nestedProp, ok := nestedData.(map[string]interface{}); ok {
					nestedParsed, err := sp.parseProperty(nestedName, nestedProp, parsed.Name, nestedRequired)
					if err != nil {
						return nil, fmt.Errorf("failed to parse nested property '%s': %w", nestedName, err)
					}
					parsed.Properties[nestedName] = nestedParsed
				}
			}
		}

		if additionalProps, ok := prop["additionalProperties"].(bool); ok {
			parsed.AdditionalProperties = additionalProps
		}
	}

	return parsed, nil
}

// extractRequiredFields extracts required field names from schema
func (sp *SchemaParser) extractRequiredFields(schemaProps map[string]interface{}) []string {
	// First check if required is defined at the schema root level
	if sp.schema.Properties != nil {
		if requiredData, ok := sp.schema.Properties["required"]; ok {
			if requiredArray, ok := requiredData.([]interface{}); ok {
				required := make([]string, 0, len(requiredArray))
				for _, item := range requiredArray {
					if fieldName, ok := item.(string); ok {
						required = append(required, fieldName)
					}
				}
				return required
			}
		}
	}

	// If not found, check in the passed properties (for nested objects)
	if requiredData, ok := schemaProps["required"]; ok {
		if requiredArray, ok := requiredData.([]interface{}); ok {
			required := make([]string, 0, len(requiredArray))
			for _, item := range requiredArray {
				if fieldName, ok := item.(string); ok {
					required = append(required, fieldName)
				}
			}
			return required
		}
	}

	return []string{}
}

// isRequired checks if a field name is in the required fields list
func (sp *SchemaParser) isRequired(fieldName string, requiredFields []string) bool {
	for _, required := range requiredFields {
		if required == fieldName {
			return true
		}
	}
	return false
}

// extractValidationRules extracts all validation rules from a parsed property
func (sp *SchemaParser) extractValidationRules(propertyName string, prop *ParsedProperty, parentPath string) []ValidationRule {
	rules := make([]ValidationRule, 0)

	// Build full property path
	fullPath := propertyName
	if parentPath != "" {
		fullPath = parentPath + "." + propertyName
	}

	// Required validation
	if prop.Required {
		rules = append(rules, ValidationRule{
			Type:         "required",
			Value:        true,
			Message:      fmt.Sprintf("Field '%s' is required", propertyName),
			PropertyPath: fullPath,
		})
	}

	// String length validations
	if prop.MinLength != nil {
		rules = append(rules, ValidationRule{
			Type:         "minLength",
			Value:        *prop.MinLength,
			Message:      fmt.Sprintf("Field '%s' must be at least %d characters", propertyName, *prop.MinLength),
			PropertyPath: fullPath,
		})
	}

	if prop.MaxLength != nil {
		rules = append(rules, ValidationRule{
			Type:         "maxLength",
			Value:        *prop.MaxLength,
			Message:      fmt.Sprintf("Field '%s' must be at most %d characters", propertyName, *prop.MaxLength),
			PropertyPath: fullPath,
		})
	}

	// Numeric validations
	if prop.Minimum != nil {
		rules = append(rules, ValidationRule{
			Type:         "minimum",
			Value:        *prop.Minimum,
			Message:      fmt.Sprintf("Field '%s' must be at least %.2f", propertyName, *prop.Minimum),
			PropertyPath: fullPath,
		})
	}

	if prop.Maximum != nil {
		rules = append(rules, ValidationRule{
			Type:         "maximum",
			Value:        *prop.Maximum,
			Message:      fmt.Sprintf("Field '%s' must be at most %.2f", propertyName, *prop.Maximum),
			PropertyPath: fullPath,
		})
	}

	// Pattern validation
	if prop.Pattern != "" {
		rules = append(rules, ValidationRule{
			Type:         "pattern",
			Value:        prop.Pattern,
			Message:      fmt.Sprintf("Field '%s' must match the required pattern", propertyName),
			PropertyPath: fullPath,
		})
	}

	// Format validation
	if prop.Format != "" {
		rules = append(rules, ValidationRule{
			Type:         "format",
			Value:        prop.Format,
			Message:      fmt.Sprintf("Field '%s' must be a valid %s", propertyName, prop.Format),
			PropertyPath: fullPath,
		})
	}

	// Enum validation
	if len(prop.Enum) > 0 {
		rules = append(rules, ValidationRule{
			Type:         "enum",
			Value:        prop.Enum,
			Message:      fmt.Sprintf("Field '%s' must be one of the allowed values", propertyName),
			PropertyPath: fullPath,
		})
	}

	// Type validation
	rules = append(rules, ValidationRule{
		Type:         "type",
		Value:        prop.Type,
		Message:      fmt.Sprintf("Field '%s' must be of type %s", propertyName, prop.Type),
		PropertyPath: fullPath,
	})

	// Recursively extract rules from nested objects
	if prop.Type == "object" && prop.Properties != nil {
		for nestedName, nestedProp := range prop.Properties {
			nestedRules := sp.extractValidationRules(nestedName, nestedProp, fullPath)
			rules = append(rules, nestedRules...)
		}
	}

	// Extract rules from array items
	if prop.Type == "array" && prop.Items != nil {
		itemRules := sp.extractValidationRules("items", prop.Items, fullPath)
		rules = append(rules, itemRules...)
	}

	return rules
}

// GetFieldMetadata returns metadata for a specific field by name
func (sp *SchemaParser) GetFieldMetadata(fieldName string) (*ParsedProperty, error) {
	analysis, err := sp.ParseSchema()
	if err != nil {
		return nil, err
	}

	if prop, exists := analysis.Properties[fieldName]; exists {
		return prop, nil
	}

	return nil, fmt.Errorf("field '%s' not found in schema", fieldName)
}

// GetNestedFieldMetadata returns metadata for a nested field using dot notation
func (sp *SchemaParser) GetNestedFieldMetadata(fieldPath string) (*ParsedProperty, error) {
	parts := strings.Split(fieldPath, ".")
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid field path: %s", fieldPath)
	}

	analysis, err := sp.ParseSchema()
	if err != nil {
		return nil, err
	}

	// Start with root property
	currentProp, exists := analysis.Properties[parts[0]]
	if !exists {
		return nil, fmt.Errorf("field '%s' not found in schema", parts[0])
	}

	// Navigate through nested properties
	for i := 1; i < len(parts); i++ {
		if currentProp.Type != "object" || currentProp.Properties == nil {
			return nil, fmt.Errorf("cannot navigate to '%s': parent is not an object", strings.Join(parts[:i+1], "."))
		}

		if nextProp, exists := currentProp.Properties[parts[i]]; exists {
			currentProp = nextProp
		} else {
			return nil, fmt.Errorf("field '%s' not found", strings.Join(parts[:i+1], "."))
		}
	}

	return currentProp, nil
}

// GetValidationRulesForField returns all validation rules for a specific field
func (sp *SchemaParser) GetValidationRulesForField(fieldName string) ([]ValidationRule, error) {
	analysis, err := sp.ParseSchema()
	if err != nil {
		return nil, err
	}

	rules := make([]ValidationRule, 0)
	for _, rule := range analysis.ValidationRules {
		// Check if rule applies to this field (considering nested paths)
		if rule.PropertyPath == fieldName || strings.HasPrefix(rule.PropertyPath, fieldName+".") {
			rules = append(rules, rule)
		}
	}

	return rules, nil
}

// ValidateFieldValue validates a field value against its schema definition
func (sp *SchemaParser) ValidateFieldValue(fieldName string, value interface{}) []ValidationRule {
	failures := make([]ValidationRule, 0)

	_, err := sp.GetFieldMetadata(fieldName)
	if err != nil {
		// If field not found, it might be valid if additional properties are allowed
		return failures
	}

	// Get validation rules for this field
	rules, err := sp.GetValidationRulesForField(fieldName)
	if err != nil {
		return failures
	}

	// Check each validation rule
	for _, rule := range rules {
		if !sp.validateSingleRule(rule, value) {
			failures = append(failures, rule)
		}
	}

	return failures
}

// validateSingleRule validates a single rule against a value
func (sp *SchemaParser) validateSingleRule(rule ValidationRule, value interface{}) bool {
	switch rule.Type {
	case "required":
		return value != nil && value != ""

	case "type":
		expectedType := rule.Value.(string)
		return sp.checkType(value, expectedType)

	case "minLength":
		if str, ok := value.(string); ok {
			minLen := rule.Value.(int)
			return len(str) >= minLen
		}
		return false

	case "maxLength":
		if str, ok := value.(string); ok {
			maxLen := rule.Value.(int)
			return len(str) <= maxLen
		}
		return false

	case "minimum":
		if num, ok := sp.toFloat64(value); ok {
			min := rule.Value.(float64)
			return num >= min
		}
		return false

	case "maximum":
		if num, ok := sp.toFloat64(value); ok {
			max := rule.Value.(float64)
			return num <= max
		}
		return false

	case "enum":
		enumValues := rule.Value.([]interface{})
		for _, enumVal := range enumValues {
			if reflect.DeepEqual(value, enumVal) {
				return true
			}
		}
		return false

	case "pattern":
		// Pattern validation would require regex - simplified for now
		return true

	case "format":
		// Format validation (email, date, etc.) - simplified for now
		return true

	default:
		return true
	}
}

// checkType checks if a value matches the expected JSON Schema type
func (sp *SchemaParser) checkType(value interface{}, expectedType string) bool {
	if value == nil {
		return true // null is valid for any type unless required
	}

	switch expectedType {
	case "string":
		_, ok := value.(string)
		return ok
	case "number":
		_, ok := sp.toFloat64(value)
		return ok
	case "integer":
		if num, ok := sp.toFloat64(value); ok {
			return num == float64(int64(num)) // Check if it's a whole number
		}
		return false
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "array":
		return reflect.TypeOf(value).Kind() == reflect.Slice
	case "object":
		return reflect.TypeOf(value).Kind() == reflect.Map
	default:
		return true
	}
}

// toFloat64 converts various numeric types to float64
func (sp *SchemaParser) toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case json.Number:
		if f, err := v.Float64(); err == nil {
			return f, true
		}
	}
	return 0, false
}
