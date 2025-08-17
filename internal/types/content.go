package types

import (
	"encoding/json"
	"time"
)

// ContentData represents the main content structure stored in content.json
type ContentData struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Sections    map[string]interface{} `json:"sections"`
	LastUpdated time.Time              `json:"last_updated"`
}

// SchemaData represents the JSON schema structure stored in schema.json
type SchemaData struct {
	Schema     string                 `json:"$schema"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

// ToJSON converts any struct to JSON string
func (c *ContentData) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ToJSON converts schema to JSON string
func (s *SchemaData) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
