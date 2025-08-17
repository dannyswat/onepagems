package managers

import (
	"encoding/json"
	"fmt"
	"onepagems/internal/types"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SchemaValidator handles comprehensive validation of content against JSON schema
type SchemaValidator struct {
	schema *types.SchemaData
	parser *SchemaParser
}

// NewSchemaValidator creates a new schema validator
func NewSchemaValidator(schema *types.SchemaData) *SchemaValidator {
	return &SchemaValidator{
		schema: schema,
		parser: NewSchemaParser(schema),
	}
}

// ValidationResult represents the result of content validation
type ValidationResult struct {
	Valid      bool                      `json:"valid"`
	Errors     []ValidationDetailError   `json:"errors"`
	Warnings   []types.ValidationWarning `json:"warnings,omitempty"`
	FieldCount int                       `json:"field_count"`
	Summary    string                    `json:"summary"`
}

// ValidationDetailError represents a validation error with detailed information
type ValidationDetailError struct {
	Field        string      `json:"field"`
	Code         string      `json:"code"`
	Message      string      `json:"message"`
	Value        interface{} `json:"value,omitempty"`
	Expected     interface{} `json:"expected,omitempty"`
	PropertyPath string      `json:"property_path"`
}

// ValidateContent validates an entire content object against the schema
func (sv *SchemaValidator) ValidateContent(content interface{}) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   make([]ValidationDetailError, 0),
		Warnings: make([]types.ValidationWarning, 0),
	}

	// Validate that content is an object if schema type is object
	if sv.schema.Type == "object" {
		contentMap, ok := content.(map[string]interface{})
		if !ok {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:    "_root",
				Code:     "invalid_type",
				Message:  "Content must be an object",
				Value:    reflect.TypeOf(content).String(),
				Expected: "object",
			})
			result.Summary = "Content type validation failed"
			return result
		}

		// Validate each property
		sv.validateObject(contentMap, "", sv.schema.Properties, result)

		// Check for required fields
		sv.validateRequiredFields(contentMap, result)
	}

	result.FieldCount = len(result.Errors) + len(result.Warnings)

	if result.Valid {
		result.Summary = "All validations passed"
	} else {
		result.Summary = fmt.Sprintf("%d validation errors found", len(result.Errors))
	}

	return result
}

// validateObject validates an object and its properties
func (sv *SchemaValidator) validateObject(obj map[string]interface{}, path string, schemaProps map[string]interface{}, result *ValidationResult) {
	// Validate each field in the object
	for fieldName, value := range obj {
		fieldPath := fieldName
		if path != "" {
			fieldPath = path + "." + fieldName
		}

		// Check if field is defined in schema
		if schemaProp, exists := schemaProps[fieldName]; exists {
			if propMap, ok := schemaProp.(map[string]interface{}); ok {
				sv.validateField(fieldName, value, propMap, fieldPath, result)
			}
		} else {
			// Check if additional properties are allowed
			// For now, we'll allow additional properties but add a warning
			result.Warnings = append(result.Warnings, types.ValidationWarning{
				Field:   fieldPath,
				Code:    "additional_property",
				Message: fmt.Sprintf("Field '%s' is not defined in schema but is allowed", fieldName),
			})
		}
	}
}

// validateField validates a single field against its schema definition
func (sv *SchemaValidator) validateField(fieldName string, value interface{}, schemaProp map[string]interface{}, fieldPath string, result *ValidationResult) {
	// Get field type
	fieldType := "string" // default
	if propType, ok := schemaProp["type"].(string); ok {
		fieldType = propType
	}

	// Type validation
	if !sv.validateType(value, fieldType) {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationDetailError{
			Field:        fieldName,
			Code:         "invalid_type",
			Message:      fmt.Sprintf("Field '%s' must be of type %s", fieldName, fieldType),
			Value:        value,
			Expected:     fieldType,
			PropertyPath: fieldPath,
		})
		return // Skip further validation if type is wrong
	}

	// String validations
	if fieldType == "string" && value != nil {
		sv.validateStringField(fieldName, value, schemaProp, fieldPath, result)
	}

	// Number validations
	if (fieldType == "number" || fieldType == "integer") && value != nil {
		sv.validateNumberField(fieldName, value, schemaProp, fieldPath, result)
	}

	// Array validations
	if fieldType == "array" && value != nil {
		sv.validateArrayField(fieldName, value, schemaProp, fieldPath, result)
	}

	// Object validations
	if fieldType == "object" && value != nil {
		sv.validateNestedObject(fieldName, value, schemaProp, fieldPath, result)
	}

	// Enum validation
	if enumValues, ok := schemaProp["enum"].([]interface{}); ok && len(enumValues) > 0 {
		sv.validateEnum(fieldName, value, enumValues, fieldPath, result)
	}

	// Format validation
	if format, ok := schemaProp["format"].(string); ok && format != "" {
		sv.validateFormat(fieldName, value, format, fieldPath, result)
	}

	// Pattern validation
	if pattern, ok := schemaProp["pattern"].(string); ok && pattern != "" {
		sv.validatePattern(fieldName, value, pattern, fieldPath, result)
	}
}

// validateType checks if value matches the expected type
func (sv *SchemaValidator) validateType(value interface{}, expectedType string) bool {
	if value == nil {
		return true // null is valid for any type unless required
	}

	switch expectedType {
	case "string":
		_, ok := value.(string)
		return ok
	case "number":
		return sv.isNumber(value)
	case "integer":
		return sv.isInteger(value)
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "array":
		return reflect.TypeOf(value).Kind() == reflect.Slice
	case "object":
		_, ok := value.(map[string]interface{})
		return ok
	default:
		return true
	}
}

// validateStringField validates string-specific constraints
func (sv *SchemaValidator) validateStringField(fieldName string, value interface{}, schemaProp map[string]interface{}, fieldPath string, result *ValidationResult) {
	str, ok := value.(string)
	if !ok {
		return
	}

	// MinLength validation
	if minLength, ok := schemaProp["minLength"].(float64); ok {
		if len(str) < int(minLength) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:        fieldName,
				Code:         "min_length",
				Message:      fmt.Sprintf("Field '%s' must be at least %d characters", fieldName, int(minLength)),
				Value:        len(str),
				Expected:     int(minLength),
				PropertyPath: fieldPath,
			})
		}
	}

	// MaxLength validation
	if maxLength, ok := schemaProp["maxLength"].(float64); ok {
		if len(str) > int(maxLength) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:        fieldName,
				Code:         "max_length",
				Message:      fmt.Sprintf("Field '%s' must be at most %d characters", fieldName, int(maxLength)),
				Value:        len(str),
				Expected:     int(maxLength),
				PropertyPath: fieldPath,
			})
		}
	}
}

// validateNumberField validates number-specific constraints
func (sv *SchemaValidator) validateNumberField(fieldName string, value interface{}, schemaProp map[string]interface{}, fieldPath string, result *ValidationResult) {
	num, ok := sv.toFloat64(value)
	if !ok {
		return
	}

	// Minimum validation
	if minimum, ok := schemaProp["minimum"].(float64); ok {
		if num < minimum {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:        fieldName,
				Code:         "minimum",
				Message:      fmt.Sprintf("Field '%s' must be at least %.2f", fieldName, minimum),
				Value:        num,
				Expected:     minimum,
				PropertyPath: fieldPath,
			})
		}
	}

	// Maximum validation
	if maximum, ok := schemaProp["maximum"].(float64); ok {
		if num > maximum {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:        fieldName,
				Code:         "maximum",
				Message:      fmt.Sprintf("Field '%s' must be at most %.2f", fieldName, maximum),
				Value:        num,
				Expected:     maximum,
				PropertyPath: fieldPath,
			})
		}
	}

	// ExclusiveMinimum validation
	if exclusiveMinimum, ok := schemaProp["exclusiveMinimum"].(float64); ok {
		if num <= exclusiveMinimum {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:        fieldName,
				Code:         "exclusive_minimum",
				Message:      fmt.Sprintf("Field '%s' must be greater than %.2f", fieldName, exclusiveMinimum),
				Value:        num,
				Expected:     exclusiveMinimum,
				PropertyPath: fieldPath,
			})
		}
	}

	// ExclusiveMaximum validation
	if exclusiveMaximum, ok := schemaProp["exclusiveMaximum"].(float64); ok {
		if num >= exclusiveMaximum {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:        fieldName,
				Code:         "exclusive_maximum",
				Message:      fmt.Sprintf("Field '%s' must be less than %.2f", fieldName, exclusiveMaximum),
				Value:        num,
				Expected:     exclusiveMaximum,
				PropertyPath: fieldPath,
			})
		}
	}

	// MultipleOf validation
	if multipleOf, ok := schemaProp["multipleOf"].(float64); ok && multipleOf > 0 {
		if remainder := num / multipleOf; remainder != float64(int64(remainder)) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:        fieldName,
				Code:         "multiple_of",
				Message:      fmt.Sprintf("Field '%s' must be a multiple of %.2f", fieldName, multipleOf),
				Value:        num,
				Expected:     multipleOf,
				PropertyPath: fieldPath,
			})
		}
	}
}

// validateArrayField validates array-specific constraints
func (sv *SchemaValidator) validateArrayField(fieldName string, value interface{}, schemaProp map[string]interface{}, fieldPath string, result *ValidationResult) {
	arr := reflect.ValueOf(value)
	if arr.Kind() != reflect.Slice {
		return
	}

	arrayLen := arr.Len()

	// MinItems validation
	if minItems, ok := schemaProp["minItems"].(float64); ok {
		if arrayLen < int(minItems) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:        fieldName,
				Code:         "min_items",
				Message:      fmt.Sprintf("Field '%s' must have at least %d items", fieldName, int(minItems)),
				Value:        arrayLen,
				Expected:     int(minItems),
				PropertyPath: fieldPath,
			})
		}
	}

	// MaxItems validation
	if maxItems, ok := schemaProp["maxItems"].(float64); ok {
		if arrayLen > int(maxItems) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:        fieldName,
				Code:         "max_items",
				Message:      fmt.Sprintf("Field '%s' must have at most %d items", fieldName, int(maxItems)),
				Value:        arrayLen,
				Expected:     int(maxItems),
				PropertyPath: fieldPath,
			})
		}
	}

	// UniqueItems validation
	if uniqueItems, ok := schemaProp["uniqueItems"].(bool); ok && uniqueItems {
		seen := make(map[string]bool)
		for i := 0; i < arrayLen; i++ {
			item := arr.Index(i).Interface()
			itemStr := fmt.Sprintf("%v", item)
			if seen[itemStr] {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationDetailError{
					Field:        fieldName,
					Code:         "unique_items",
					Message:      fmt.Sprintf("Field '%s' must have unique items", fieldName),
					Value:        item,
					PropertyPath: fieldPath,
				})
				break
			}
			seen[itemStr] = true
		}
	}

	// Validate array items against items schema
	if items, ok := schemaProp["items"].(map[string]interface{}); ok {
		for i := 0; i < arrayLen; i++ {
			item := arr.Index(i).Interface()
			itemPath := fmt.Sprintf("%s[%d]", fieldPath, i)
			sv.validateField(fmt.Sprintf("%s[%d]", fieldName, i), item, items, itemPath, result)
		}
	}
}

// validateNestedObject validates nested object fields
func (sv *SchemaValidator) validateNestedObject(fieldName string, value interface{}, schemaProp map[string]interface{}, fieldPath string, result *ValidationResult) {
	objMap, ok := value.(map[string]interface{})
	if !ok {
		return
	}

	// Get nested properties
	if properties, ok := schemaProp["properties"].(map[string]interface{}); ok {
		sv.validateObject(objMap, fieldPath, properties, result)
	}

	// Validate required fields for this nested object
	if required, ok := schemaProp["required"].([]interface{}); ok {
		for _, reqField := range required {
			if reqFieldName, ok := reqField.(string); ok {
				if _, exists := objMap[reqFieldName]; !exists {
					result.Valid = false
					result.Errors = append(result.Errors, ValidationDetailError{
						Field:        fmt.Sprintf("%s.%s", fieldName, reqFieldName),
						Code:         "required",
						Message:      fmt.Sprintf("Required field '%s.%s' is missing", fieldName, reqFieldName),
						PropertyPath: fmt.Sprintf("%s.%s", fieldPath, reqFieldName),
					})
				}
			}
		}
	}
}

// validateEnum validates that value is one of the allowed enum values
func (sv *SchemaValidator) validateEnum(fieldName string, value interface{}, enumValues []interface{}, fieldPath string, result *ValidationResult) {
	for _, enumValue := range enumValues {
		if reflect.DeepEqual(value, enumValue) {
			return // Valid enum value found
		}
	}

	result.Valid = false
	result.Errors = append(result.Errors, ValidationDetailError{
		Field:        fieldName,
		Code:         "enum",
		Message:      fmt.Sprintf("Field '%s' must be one of the allowed values", fieldName),
		Value:        value,
		Expected:     enumValues,
		PropertyPath: fieldPath,
	})
}

// validateFormat validates string format constraints (email, date, etc.)
func (sv *SchemaValidator) validateFormat(fieldName string, value interface{}, format string, fieldPath string, result *ValidationResult) {
	str, ok := value.(string)
	if !ok || str == "" {
		return // Skip format validation for non-strings or empty strings
	}

	switch format {
	case "email":
		if !sv.isValidEmail(str) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:        fieldName,
				Code:         "format_email",
				Message:      fmt.Sprintf("Field '%s' must be a valid email address", fieldName),
				Value:        str,
				Expected:     "valid email format",
				PropertyPath: fieldPath,
			})
		}
	case "date":
		if !sv.isValidDate(str) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:        fieldName,
				Code:         "format_date",
				Message:      fmt.Sprintf("Field '%s' must be a valid date (YYYY-MM-DD)", fieldName),
				Value:        str,
				Expected:     "YYYY-MM-DD format",
				PropertyPath: fieldPath,
			})
		}
	case "date-time":
		if !sv.isValidDateTime(str) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:        fieldName,
				Code:         "format_datetime",
				Message:      fmt.Sprintf("Field '%s' must be a valid date-time (RFC3339)", fieldName),
				Value:        str,
				Expected:     "RFC3339 format",
				PropertyPath: fieldPath,
			})
		}
	case "uri":
		if !sv.isValidURI(str) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:        fieldName,
				Code:         "format_uri",
				Message:      fmt.Sprintf("Field '%s' must be a valid URI", fieldName),
				Value:        str,
				Expected:     "valid URI format",
				PropertyPath: fieldPath,
			})
		}
	case "ipv4":
		if !sv.isValidIPv4(str) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:        fieldName,
				Code:         "format_ipv4",
				Message:      fmt.Sprintf("Field '%s' must be a valid IPv4 address", fieldName),
				Value:        str,
				Expected:     "IPv4 format",
				PropertyPath: fieldPath,
			})
		}
	case "ipv6":
		if !sv.isValidIPv6(str) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationDetailError{
				Field:        fieldName,
				Code:         "format_ipv6",
				Message:      fmt.Sprintf("Field '%s' must be a valid IPv6 address", fieldName),
				Value:        str,
				Expected:     "IPv6 format",
				PropertyPath: fieldPath,
			})
		}
	}
}

// validatePattern validates string against regex pattern
func (sv *SchemaValidator) validatePattern(fieldName string, value interface{}, pattern string, fieldPath string, result *ValidationResult) {
	str, ok := value.(string)
	if !ok {
		return
	}

	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		result.Warnings = append(result.Warnings, types.ValidationWarning{
			Field:   fieldName,
			Code:    "invalid_pattern",
			Message: fmt.Sprintf("Invalid regex pattern for field '%s': %s", fieldName, err.Error()),
		})
		return
	}

	if !matched {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationDetailError{
			Field:        fieldName,
			Code:         "pattern",
			Message:      fmt.Sprintf("Field '%s' must match the required pattern", fieldName),
			Value:        str,
			Expected:     pattern,
			PropertyPath: fieldPath,
		})
	}
}

// validateRequiredFields checks that all required fields are present
func (sv *SchemaValidator) validateRequiredFields(content map[string]interface{}, result *ValidationResult) {
	// Check required fields defined at property level
	for propName, propData := range sv.schema.Properties {
		if propMap, ok := propData.(map[string]interface{}); ok {
			if required, ok := propMap["required"].(bool); ok && required {
				if _, exists := content[propName]; !exists {
					result.Valid = false
					result.Errors = append(result.Errors, ValidationDetailError{
						Field:        propName,
						Code:         "required",
						Message:      fmt.Sprintf("Required field '%s' is missing", propName),
						PropertyPath: propName,
					})
				}
			}
		}
	}
}

// ValidateFieldValue validates a single field value against schema
func (sv *SchemaValidator) ValidateFieldValue(fieldName string, value interface{}) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   make([]ValidationDetailError, 0),
		Warnings: make([]types.ValidationWarning, 0),
	}

	// Find the field schema
	if schemaProp, exists := sv.schema.Properties[fieldName]; exists {
		if propMap, ok := schemaProp.(map[string]interface{}); ok {
			sv.validateField(fieldName, value, propMap, fieldName, result)
		}
	} else {
		result.Warnings = append(result.Warnings, types.ValidationWarning{
			Field:   fieldName,
			Code:    "unknown_field",
			Message: fmt.Sprintf("Field '%s' is not defined in schema", fieldName),
		})
	}

	if result.Valid {
		result.Summary = "Field validation passed"
	} else {
		result.Summary = fmt.Sprintf("Field validation failed with %d errors", len(result.Errors))
	}

	return result
}

// Helper functions for type checking and format validation

func (sv *SchemaValidator) isNumber(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	case json.Number:
		return true
	default:
		// Try to parse as number
		if str, ok := value.(string); ok {
			_, err := strconv.ParseFloat(str, 64)
			return err == nil
		}
		return false
	}
}

func (sv *SchemaValidator) isInteger(value interface{}) bool {
	if !sv.isNumber(value) {
		return false
	}

	num, _ := sv.toFloat64(value)
	return num == float64(int64(num))
}

func (sv *SchemaValidator) toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case json.Number:
		if f, err := v.Float64(); err == nil {
			return f, true
		}
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

func (sv *SchemaValidator) isValidEmail(email string) bool {
	// Basic email validation regex
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}

func (sv *SchemaValidator) isValidDate(date string) bool {
	// YYYY-MM-DD format
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func (sv *SchemaValidator) isValidDateTime(datetime string) bool {
	// RFC3339 format
	_, err := time.Parse(time.RFC3339, datetime)
	return err == nil
}

func (sv *SchemaValidator) isValidURI(uri string) bool {
	// Basic URI validation - contains scheme and host
	uriRegex := `^[a-zA-Z][a-zA-Z0-9+.-]*://[^\s]+$`
	matched, _ := regexp.MatchString(uriRegex, uri)
	return matched
}

func (sv *SchemaValidator) isValidIPv4(ip string) bool {
	// IPv4 validation
	ipv4Regex := `^(\d{1,3}\.){3}\d{1,3}$`
	matched, _ := regexp.MatchString(ipv4Regex, ip)
	if !matched {
		return false
	}

	// Check if each octet is <= 255
	parts := strings.Split(ip, ".")
	for _, part := range parts {
		if num, err := strconv.Atoi(part); err != nil || num > 255 {
			return false
		}
	}
	return true
}

func (sv *SchemaValidator) isValidIPv6(ip string) bool {
	// Basic IPv6 validation
	ipv6Regex := `^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|^::1$|^::$`
	matched, _ := regexp.MatchString(ipv6Regex, ip)
	return matched
}

// GenerateValidationReport generates a detailed validation report
func (sv *SchemaValidator) GenerateValidationReport(content interface{}) map[string]interface{} {
	result := sv.ValidateContent(content)

	report := map[string]interface{}{
		"validation_result": result,
		"timestamp":         time.Now().UTC().Format(time.RFC3339),
		"schema_info": map[string]interface{}{
			"type":             sv.schema.Type,
			"properties_count": len(sv.schema.Properties),
		},
		"error_summary": map[string]interface{}{
			"total_errors":   len(result.Errors),
			"total_warnings": len(result.Warnings),
			"error_codes":    sv.getErrorCodes(result.Errors),
		},
	}

	return report
}

// getErrorCodes extracts unique error codes from validation errors
func (sv *SchemaValidator) getErrorCodes(errors []ValidationDetailError) []string {
	codes := make(map[string]bool)
	for _, err := range errors {
		codes[err.Code] = true
	}

	result := make([]string, 0, len(codes))
	for code := range codes {
		result = append(result, code)
	}
	return result
}
