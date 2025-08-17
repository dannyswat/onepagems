package internal

import (
	"context"
	"encoding/json"
	"time"
)

// SessionContextKey is the key used to store session in context
type SessionContextKey string

const SessionKey SessionContextKey = "session"

// SessionContext creates a new context with the session
func SessionContext(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, SessionKey, session)
}

// SessionFromContext retrieves the session from context
func SessionFromContext(ctx context.Context) (*Session, bool) {
	session, ok := ctx.Value(SessionKey).(*Session)
	return session, ok
}

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

// Config represents the application configuration
type Config struct {
	Port           string `json:"port"`
	AdminUsername  string `json:"admin_username"`
	AdminPassword  string `json:"admin_password"`
	UploadMaxSize  int64  `json:"upload_max_size"`
	SessionTimeout int    `json:"session_timeout"` // in minutes
	DataDir        string `json:"data_dir"`
	StaticDir      string `json:"static_dir"`
	TemplatesDir   string `json:"templates_dir"`
}

// Session represents a user session
type Session struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
}

// ImageInfo represents information about an uploaded image
type ImageInfo struct {
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
	URL          string    `json:"url"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// ValidationWarning represents a non-critical validation issue
type ValidationWarning struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message,omitempty"`
	Data    interface{}            `json:"data,omitempty"`
	Errors  []ValidationError      `json:"errors,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// FormField represents a field in a dynamically generated form
type FormField struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Label       string      `json:"label"`
	Required    bool        `json:"required"`
	Placeholder string      `json:"placeholder,omitempty"`
	Options     []string    `json:"options,omitempty"`
	Value       interface{} `json:"value,omitempty"`
	Format      string      `json:"format,omitempty"`
	Description string      `json:"description,omitempty"`
}

// GeneratedForm represents a complete form generated from schema
type GeneratedForm struct {
	Fields []FormField `json:"fields"`
	Action string      `json:"action"`
	Method string      `json:"method"`
}

// FileBackup represents backup file information
type FileBackup struct {
	OriginalPath string    `json:"original_path"`
	BackupPath   string    `json:"backup_path"`
	CreatedAt    time.Time `json:"created_at"`
	Size         int64     `json:"size"`
}

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

// FileInfo represents information about files in the system
type FileInfo struct {
	Path        string    `json:"path"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ModifiedAt  time.Time `json:"modified_at"`
	IsDirectory bool      `json:"is_directory"`
	ContentType string    `json:"content_type,omitempty"`
	HasBackup   bool      `json:"has_backup"`
	BackupAge   *int64    `json:"backup_age,omitempty"` // seconds since backup
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Port:           "8080",
		AdminUsername:  "admin",
		AdminPassword:  "",              // Will be set to hashed "admin123" in ValidateConfig
		UploadMaxSize:  5 * 1024 * 1024, // 5MB
		SessionTimeout: 60,              // 60 minutes
		DataDir:        "./data",
		StaticDir:      "./static",
		TemplatesDir:   "./templates",
	}
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

// IsExpired checks if a session is expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt) || !s.IsActive
}

// Extend extends the session expiration time
func (s *Session) Extend(minutes int) {
	s.ExpiresAt = time.Now().Add(time.Duration(minutes) * time.Minute)
}

// NewAPIResponse creates a new API response
func NewAPIResponse(success bool, message string) *APIResponse {
	return &APIResponse{
		Success: success,
		Message: message,
		Meta:    make(map[string]interface{}),
	}
}

// AddError adds a validation error to the API response
func (r *APIResponse) AddError(field, message, code string) {
	if r.Errors == nil {
		r.Errors = make([]ValidationError, 0)
	}
	r.Errors = append(r.Errors, ValidationError{
		Field:   field,
		Message: message,
		Code:    code,
	})
}

// SetData sets the data field of the API response
func (r *APIResponse) SetData(data interface{}) {
	r.Data = data
}
