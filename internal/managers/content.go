package managers

import (
	"encoding/json"
	"fmt"
	"time"

	"onepagems/internal/types"
)

// ContentManager handles content.json operations
type ContentManager struct {
	storage *FileStorage
	dataDir string
}

// NewContentManager creates a new content manager
func NewContentManager(storage *FileStorage, dataDir string) *ContentManager {
	return &ContentManager{
		storage: storage,
		dataDir: dataDir,
	}
}

// contentFilePath returns the filename for content.json
func (cm *ContentManager) contentFilePath() string {
	return "content.json"
}

// LoadContent loads content from content.json or creates default if not exists
func (cm *ContentManager) LoadContent() (*types.ContentData, error) {
	contentFilename := cm.contentFilePath()

	// Check if content.json exists
	if !cm.storage.FileExists(contentFilename) {
		// Create default content
		defaultContent := cm.createDefaultContent()
		if err := cm.SaveContent(defaultContent); err != nil {
			return nil, fmt.Errorf("failed to create default content: %w", err)
		}
		return defaultContent, nil
	}

	// Load existing content
	var content types.ContentData
	if err := cm.storage.ReadJSONFile(contentFilename, &content); err != nil {
		return nil, fmt.Errorf("failed to read content file: %w", err)
	}

	// Validate content structure
	if err := cm.validateContent(&content); err != nil {
		return nil, fmt.Errorf("content validation failed: %w", err)
	}

	return &content, nil
}

// SaveContent saves content to content.json with backup
func (cm *ContentManager) SaveContent(content *types.ContentData) error {
	if content == nil {
		return fmt.Errorf("content cannot be nil")
	}

	// Update last updated timestamp
	content.LastUpdated = time.Now()

	// Validate content before saving
	if err := cm.validateContent(content); err != nil {
		return fmt.Errorf("content validation failed: %w", err)
	}

	// Save with backup
	contentFilename := cm.contentFilePath()
	if err := cm.storage.WriteJSONFile(contentFilename, content); err != nil {
		return fmt.Errorf("failed to save content file: %w", err)
	}

	return nil
}

// UpdateContent updates specific fields in the content
func (cm *ContentManager) UpdateContent(updates map[string]interface{}) error {
	// Load current content
	content, err := cm.LoadContent()
	if err != nil {
		return fmt.Errorf("failed to load current content: %w", err)
	}

	// Apply updates
	for key, value := range updates {
		switch key {
		case "title":
			if title, ok := value.(string); ok {
				content.Title = title
			} else {
				return fmt.Errorf("title must be a string")
			}
		case "description":
			if description, ok := value.(string); ok {
				content.Description = description
			} else {
				return fmt.Errorf("description must be a string")
			}
		case "sections":
			if sections, ok := value.(map[string]interface{}); ok {
				content.Sections = sections
			} else {
				return fmt.Errorf("sections must be a map")
			}
		default:
			return fmt.Errorf("unknown field: %s", key)
		}
	}

	// Save updated content
	return cm.SaveContent(content)
}

// BackupContent creates a backup of the current content
func (cm *ContentManager) BackupContent() error {
	contentFilename := cm.contentFilePath()
	return cm.storage.CreateBackup(contentFilename)
}

// RestoreContent restores content from backup
func (cm *ContentManager) RestoreContent() error {
	contentFilename := cm.contentFilePath()
	return cm.storage.RestoreFromBackup(contentFilename)
}

// GetContentSummary returns a summary of the current content
func (cm *ContentManager) GetContentSummary() (map[string]interface{}, error) {
	content, err := cm.LoadContent()
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"title":        content.Title,
		"description":  content.Description,
		"sections":     len(content.Sections),
		"last_updated": content.LastUpdated,
	}

	// Add section names
	sectionNames := make([]string, 0, len(content.Sections))
	for name := range content.Sections {
		sectionNames = append(sectionNames, name)
	}
	summary["section_names"] = sectionNames

	return summary, nil
}

// createDefaultContent creates default content structure
func (cm *ContentManager) createDefaultContent() *types.ContentData {
	return &types.ContentData{
		Title:       "Welcome to OnePage CMS",
		Description: "A simple, lightweight content management system",
		Sections: map[string]interface{}{
			"hero": map[string]interface{}{
				"title":    "Your Website Title",
				"subtitle": "Welcome to your new website",
				"content":  "This is the main hero section of your website. You can edit this content through the admin panel.",
			},
			"about": map[string]interface{}{
				"title":   "About Us",
				"content": "Tell your visitors about yourself or your organization.",
			},
			"contact": map[string]interface{}{
				"title":   "Contact",
				"email":   "contact@example.com",
				"phone":   "+1 (555) 123-4567",
				"address": "123 Main St, Anytown USA",
			},
		},
		LastUpdated: time.Now(),
	}
}

// validateContent validates the content structure
func (cm *ContentManager) validateContent(content *types.ContentData) error {
	if content == nil {
		return fmt.Errorf("content cannot be nil")
	}

	if content.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	if content.Sections == nil {
		content.Sections = make(map[string]interface{})
	}

	// Validate sections structure
	for sectionName, sectionData := range content.Sections {
		if sectionName == "" {
			return fmt.Errorf("section name cannot be empty")
		}

		// Ensure section data is a map
		if _, ok := sectionData.(map[string]interface{}); !ok {
			return fmt.Errorf("section '%s' must be an object", sectionName)
		}
	}

	return nil
}

// ExportContent exports content as JSON for external use
func (cm *ContentManager) ExportContent() ([]byte, error) {
	content, err := cm.LoadContent()
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(content, "", "  ")
}

// ImportContent imports content from JSON data
func (cm *ContentManager) ImportContent(data []byte) error {
	var content types.ContentData
	if err := json.Unmarshal(data, &content); err != nil {
		return fmt.Errorf("failed to parse imported content: %w", err)
	}

	return cm.SaveContent(&content)
}
