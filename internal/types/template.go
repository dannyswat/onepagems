package types

import "time"

// TemplateData represents data passed to HTML templates
type TemplateData struct {
	Title       string                 `json:"title"`
	Content     ContentData            `json:"content"`
	Schema      SchemaData             `json:"schema"`
	Form        *GeneratedForm         `json:"form,omitempty"`
	Images      []ImageInfo            `json:"images,omitempty"`
	Session     *Session               `json:"session,omitempty"`
	Errors      []ValidationError      `json:"errors,omitempty"`
	Messages    []string               `json:"messages,omitempty"`
	CSRFToken   string                 `json:"csrf_token,omitempty"`
	CurrentPage string                 `json:"current_page,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
}

// GenerationResult represents the result of HTML generation
type GenerationResult struct {
	Success     bool      `json:"success"`
	OutputPath  string    `json:"output_path,omitempty"`
	GeneratedAt time.Time `json:"generated_at"`
	Errors      []string  `json:"errors,omitempty"`
	Size        int64     `json:"size,omitempty"`
}
