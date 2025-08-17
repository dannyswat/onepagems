package types

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
