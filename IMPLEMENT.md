# OnePage CMS - Implementation Guide

## Overview

This document outlines the step-by-step implementation plan for the OnePage CMS, a configurable content management system that generates static HTML pages from JSON schema-driven content.

## Prerequisites

- Go 1.21 or later
- Basic understanding of Go templates and JSON schema
- Web development knowledge (HTML, CSS, JavaScript)

## Implementation Phases

### Phase 1: Project Setup and Basic Structure

#### Step 1.1: Initialize Go Project
```bash
# Initialize Go module
go mod init onepagems

# Create basic directory structure
mkdir -p cmd internal templates data static/codemirror
mkdir -p data/images
```

#### Step 1.2: Create Basic File Structure
```
onepagems/
├── cmd/
│   └── main.go
├── internal/
│   ├── config.go
│   ├── server.go
│   └── types.go
├── templates/
├── data/
│   └── images/
├── static/
│   └── codemirror/
└── go.mod
```

#### Step 1.3: Define Core Data Types
Create `internal/types.go`:
- ContentData struct
- SchemaData struct  
- Config struct
- Session management types

#### Step 1.4: Basic HTTP Server Setup
Create `cmd/main.go`:
- Basic HTTP server with routing
- Static file serving
- Environment variable loading

**Deliverable**: Basic web server that serves static files and responds to health checks.

---

### Phase 2: File Management and Storage

#### Step 2.1: File Operations Module
Create `internal/storage.go`:
- JSON file read/write operations
- Backup file management (.bak files)
- Directory creation and validation
- File existence checks

#### Step 2.2: Content Management
Create `internal/content.go`:
- Load/save content.json
- Content validation
- Default content generation
- Backup content before changes

#### Step 2.3: Schema Management  
Create `internal/schema.go`:
- Load/save schema.json
- Schema validation using JSON Schema Draft 07
- Default schema generation
- Schema backup management

#### Step 2.4: Template Management
Create `internal/template.go`:
- Load/save template.html
- Template validation
- Default template generation
- Template backup management

**Deliverable**: File management system with backup functionality.

---

### Phase 3: Authentication System

#### Step 3.1: Session Management
Create `internal/auth.go`:
- Simple session storage (in-memory map)
- Session creation and validation
- Session timeout handling
- CSRF token generation (optional)

#### Step 3.2: Authentication Handlers
Create authentication endpoints:
- `POST /admin/login` - Username/password validation
- `POST /admin/logout` - Session invalidation
- Authentication middleware for admin routes

#### Step 3.3: Security Middleware
- Admin route protection
- Basic input sanitization
- Session validation middleware

**Deliverable**: Working login/logout system with session management.

---

### Phase 4: JSON Schema Engine

#### Step 4.1: Schema Parser
Create `internal/schema_parser.go`:
- Parse JSON schema definitions
- Extract field types and metadata
- Support for nested objects and arrays
- Field validation rules extraction

#### Step 4.2: Schema Validator
- Validate content against schema
- Generate validation error messages
- Support for required fields
- Type validation (string, email, etc.)

#### Step 4.3: Form Generator ✅ COMPLETED
Create `internal/managers/form_generator.go`:
- ✅ Generate HTML forms from schema
- ✅ Support different input types:
  - text, textarea, email, number, checkbox
  - select, multiselect options
  - image picker integration
  - rich text editor integration
  - array/object field handling
- ✅ Field labeling, descriptions, and help text
- ✅ Validation constraints integration
- ✅ Nested object support
- ✅ Schema integration with parser and validator

**Deliverable**: ✅ Dynamic form generation based on JSON schema.

---

### Phase 5: Admin Interface

#### Step 5.1: Admin Panel Template
Create `templates/admin.html`:
- Main admin dashboard layout
- Navigation between different sections
- Content editing interface
- Form styling and UX

#### Step 5.2: Admin Handlers
Create `internal/admin_handlers.go`:
- `GET /admin` - Admin dashboard
- `POST /admin/content` - Content update
- Form processing and validation
- Success/error message handling

#### Step 5.3: Content Editor Interface
- Schema-driven form rendering
- Field validation feedback
- Save/cancel functionality
- Preview capability (basic)

**Deliverable**: Functional admin interface for content editing.

---

### Phase 6: Image Management

#### Step 6.1: Image Upload
Create `internal/image_handler.go`:
- `POST /admin/upload` - File upload handling
- File type validation (JPEG, PNG, GIF, WebP)
- File size validation
- Unique filename generation
- Image storage in data/images/

#### Step 6.2: Image Browser
- `GET /admin/images` - Image listing API
- `DELETE /admin/images/:filename` - Image deletion
- Thumbnail generation (optional)
- Image selection interface

#### Step 6.3: Image Integration
- Image picker component for schema fields
- Image URL generation
- Integration with form generator
- Image preview in forms

**Deliverable**: Complete image upload and management system.

---

### Phase 7: Code Editors

#### Step 7.1: CodeMirror Integration
Set up CodeMirror in `static/codemirror/`:
- Download CodeMirror assets
- JSON mode for schema editing
- HTML mode for template editing
- Basic themes and styling

#### Step 7.2: Schema Editor
Create `templates/schema_editor.html`:
- `GET /admin/editor/schema` - Schema editor interface
- `POST /admin/schema` - Schema update handler
- Syntax highlighting for JSON
- Validation feedback
- Save/cancel functionality

#### Step 7.3: Template Editor
Create `templates/template_editor.html`:
- `GET /admin/editor/template` - Template editor interface
- `POST /admin/template` - Template update handler
- HTML syntax highlighting
- Go template syntax support
- Template validation

**Deliverable**: Web-based code editors for schema and template customization.

---

### Phase 8: HTML Generation Engine

#### Step 8.1: Template Engine
Create `internal/generator.go`:
- Load HTML template from file
- Parse Go template syntax
- Content data binding
- Template execution
- Error handling

#### Step 8.2: Static Page Generation
- Combine template + content + schema
- Generate static index.html
- Handle template errors gracefully
- Backup previous generated file

#### Step 8.3: Generation Triggers
- Auto-generation on content save
- Manual generation endpoint
- Generation status feedback
- Error reporting

**Deliverable**: Working static HTML generation from templates and content.

---

### Phase 9: Default Content and Setup

#### Step 9.1: Default Schema
Create default schema in `data/schema.json`:
- Basic website structure (hero, about, services, contact)
- Common field types and formats
- Reasonable defaults and titles
- Documentation comments

#### Step 9.2: Default Template
Create default template in `data/template.html`:
- Modern, responsive HTML layout
- Bootstrap or simple CSS framework
- Go template placeholders
- Mobile-friendly design
- Basic styling

#### Step 9.3: Default Content
Create default content in `data/content.json`:
- Sample data matching default schema
- Placeholder text and images
- Working example content
- Proper structure validation

#### Step 9.4: Application Initialization
- Auto-create data directory structure
- Generate default files if missing
- Environment setup validation
- First-run setup process

**Deliverable**: Working application with sensible defaults.

---

### Phase 10: Integration and Polish

#### Step 10.1: Error Handling
- Comprehensive error handling throughout
- User-friendly error messages
- Logging for debugging
- Graceful failure recovery

#### Step 10.2: Validation and Security
- Input sanitization
- File upload security
- Session security
- CSRF protection (basic)

#### Step 10.3: User Experience
- Loading indicators
- Success/error notifications
- Form validation feedback
- Intuitive navigation

#### Step 10.4: Documentation
- README.md with setup instructions
- API documentation
- Schema format documentation
- Template syntax guide

**Deliverable**: Production-ready application with documentation.

---

## Implementation Order

### Week 1: Foundation
- Days 1-2: Phase 1 (Project Setup)
- Days 3-4: Phase 2 (File Management)
- Days 5-7: Phase 3 (Authentication)

### Week 2: Core Features
- Days 1-3: Phase 4 (Schema Engine)
- Days 4-5: Phase 5 (Admin Interface)
- Days 6-7: Phase 6 (Image Management)

### Week 3: Advanced Features
- Days 1-3: Phase 7 (Code Editors)
- Days 4-5: Phase 8 (HTML Generation)
- Days 6-7: Phase 9 (Default Content)

### Week 4: Polish and Testing
- Days 1-3: Phase 10 (Integration)
- Days 4-5: Testing and bug fixes
- Days 6-7: Documentation and deployment prep

## Key Dependencies

### Go Packages
```go
// Standard library
"encoding/json"
"html/template"
"net/http"
"os"
"path/filepath"
"time"

// Third-party (minimal)
// Consider: github.com/xeipuuv/gojsonschema for JSON schema validation
```

### Frontend Assets
- CodeMirror (for code editors)
- Basic CSS framework (or custom CSS)
- Minimal JavaScript for UX enhancements

## Testing Strategy

### Unit Tests
- File operations
- Schema validation
- Template generation
- Authentication logic

### Integration Tests
- End-to-end content management flow
- API endpoint testing
- File upload/download
- HTML generation

### Manual Testing
- Admin interface usability
- Content editing workflow
- Image management
- Generated site validation

## Deployment Considerations

### Build Process
```bash
# Build single binary
go build -o onepagems cmd/main.go

# Embed static assets (using go:embed)
# Bundle templates and static files
```

### Runtime Requirements
- Write permissions to data/ directory
- Network access on specified port
- Minimal memory footprint
- Single binary deployment

## Success Criteria

### Minimum Viable Product (MVP)
- ✅ Admin login/logout works
- ✅ Content can be edited through forms
- ✅ Images can be uploaded and managed
- ✅ Static HTML is generated correctly
- ✅ Schema and template are editable
- ✅ Application runs as single binary

### Full Feature Set
- ✅ Dynamic form generation from schema
- ✅ Code editors with syntax highlighting
- ✅ Real-time validation feedback
- ✅ Backup and recovery functionality
- ✅ Professional admin interface
- ✅ Mobile-responsive generated sites
- ✅ Comprehensive error handling
- ✅ Security best practices

## Risk Mitigation

### Technical Risks
- **JSON Schema Complexity**: Start with simple schema support, expand gradually
- **Template Security**: Use Go's built-in template security features
- **File Handling**: Implement robust error handling and validation
- **Browser Compatibility**: Test code editors across browsers

### Scope Risks
- **Feature Creep**: Stick to MVP first, then iterate
- **Over-Engineering**: Keep it simple, avoid unnecessary complexity
- **Performance**: Profile and optimize only when needed

## Getting Started

1. **Clone/Initialize**: Set up the project structure
2. **Environment**: Configure development environment
3. **Phase 1**: Start with basic HTTP server
4. **Iterate**: Complete each phase before moving to next
5. **Test**: Validate functionality at each step
6. **Deploy**: Test deployment process early

This implementation plan provides a clear roadmap from initial setup to production-ready application, with each phase building upon the previous one while maintaining the core principles of simplicity, lightweight architecture, and security.
