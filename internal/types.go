// Package internal provides backward compatibility for the refactored types.
// All types have been moved to the internal/types package for better organization.
package internal

import "onepagems/internal/types"

// Re-export all types from the types package for backward compatibility

// Session types
type SessionContextKey = types.SessionContextKey
type Session = types.Session

// Content types
type ContentData = types.ContentData
type SchemaData = types.SchemaData

// Configuration types
type Config = types.Config

// Validation types
type ValidationError = types.ValidationError
type ValidationWarning = types.ValidationWarning

// API types
type APIResponse = types.APIResponse
type FormField = types.FormField
type GeneratedForm = types.GeneratedForm

// File types
type ImageInfo = types.ImageInfo
type FileBackup = types.FileBackup
type FileInfo = types.FileInfo

// Template types
type TemplateData = types.TemplateData
type GenerationResult = types.GenerationResult

// Re-export constants and functions
const SessionKey = types.SessionKey

var (
	SessionContext     = types.SessionContext
	SessionFromContext = types.SessionFromContext
	DefaultConfig      = types.DefaultConfig
	NewAPIResponse     = types.NewAPIResponse
)
