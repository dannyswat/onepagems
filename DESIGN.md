# OnePage CMS - Simplified Design Document

## Overview

OnePage CMS is a minimal content management system for managing a single HTML page. Built with Go, it generates a static `index.html` file after content updates.

## Design Principles

- **Simple**: Minimal features, easy to use
- **Lightweight**: No database, file-based storage
- **Secure**: Basic authentication only

## Architecture

### High-Level Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Admin Panel   │    │   index.html    │    │   Images        │
│   (Web UI)      │    │   (Generated)   │    │   (Static)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Go Backend    │
                    │   (HTTP Server) │
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │   JSON Config   │
                    │   (Content)     │
                    └─────────────────┘
```

### Component Breakdown

#### 1. HTTP Server (Go)
- **Port**: 8080
- **Framework**: Standard library `net/http`
- **Features**: Basic authentication, file serving
- **Routes**: Admin panel and static file serving

#### 2. File Storage
- **Content**: JSON file for content storage
- **Template**: Single predefined HTML template
- **Images**: Static image directory with browser

#### 3. Template Engine
- **Engine**: Go's built-in `html/template`
- **Template**: Single predefined layout
- **Output**: Generates static `index.html`

#### 4. Authentication
- **Method**: Simple session-based auth
- **Credential**: Single admin user (hardcoded or env var)

## Data Model

### Configurable Content Structure

The content structure is defined by a JSON schema that can be edited through the admin interface. This allows complete customization of the content fields and structure.

#### Default JSON Schema (`schema.json`)

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "title": {
      "type": "string",
      "title": "Website Title",
      "default": "My Website"
    },
    "description": {
      "type": "string",
      "title": "Website Description",
      "default": "Welcome to my website"
    },
    "sections": {
      "type": "object",
      "properties": {
        "hero": {
          "type": "object",
          "title": "Hero Section",
          "properties": {
            "title": {"type": "string", "title": "Hero Title"},
            "subtitle": {"type": "string", "title": "Hero Subtitle"},
            "image": {"type": "string", "format": "image", "title": "Hero Image"},
            "button_text": {"type": "string", "title": "Button Text"},
            "button_link": {"type": "string", "title": "Button Link"}
          }
        },
        "about": {
          "type": "object",
          "title": "About Section",
          "properties": {
            "title": {"type": "string", "title": "About Title"},
            "content": {"type": "string", "format": "textarea", "title": "About Content"},
            "image": {"type": "string", "format": "image", "title": "About Image"}
          }
        },
        "services": {
          "type": "object",
          "title": "Services Section",
          "properties": {
            "title": {"type": "string", "title": "Services Title"},
            "items": {
              "type": "array",
              "title": "Service Items",
              "items": {
                "type": "object",
                "properties": {
                  "title": {"type": "string", "title": "Service Title"},
                  "description": {"type": "string", "title": "Service Description"},
                  "image": {"type": "string", "format": "image", "title": "Service Image"}
                }
              }
            }
          }
        },
        "contact": {
          "type": "object",
          "title": "Contact Section",
          "properties": {
            "title": {"type": "string", "title": "Contact Title"},
            "email": {"type": "string", "format": "email", "title": "Email"},
            "phone": {"type": "string", "title": "Phone"},
            "address": {"type": "string", "title": "Address"}
          }
        }
      }
    }
  }
}
```

#### Content Data (`content.json`)

```json
{
  "title": "Website Title",
  "description": "Website Description",
  "sections": {
    "hero": {
      "title": "Hero Title",
      "subtitle": "Hero Subtitle", 
      "image": "/images/hero.jpg",
      "button_text": "Call to Action",
      "button_link": "#contact"
    },
    "about": {
      "title": "About Us",
      "content": "About content here...",
      "image": "/images/about.jpg"
    },
    "services": {
      "title": "Our Services",
      "items": [
        {
          "title": "Service 1",
          "description": "Service description",
          "image": "/images/service1.jpg"
        }
      ]
    },
    "contact": {
      "title": "Contact Us",
      "email": "contact@example.com",
      "phone": "+1234567890",
      "address": "Your Address"
    }
  },
  "last_updated": "2025-08-15T10:30:00Z"
}
```

### File Structure

```
data/
├── content.json          # Main content file
├── content.json.bak     # Previous version backup
├── schema.json          # JSON schema definition
├── schema.json.bak      # Schema backup
├── template.html        # HTML template
├── template.html.bak    # Template backup
└── images/              # Image storage
    ├── hero.jpg
    ├── about.jpg
    └── ...
```

## API Design

### Admin API Endpoints

```
POST   /admin/login              # Admin authentication
POST   /admin/logout             # Admin logout
GET    /admin                    # Admin panel interface
POST   /admin/content            # Update content and regenerate
POST   /admin/upload             # Image upload
GET    /admin/images             # Image browser
DELETE /admin/images/:filename   # Delete image
GET    /admin/schema             # Get JSON schema
POST   /admin/schema             # Update JSON schema
GET    /admin/template           # Get HTML template
POST   /admin/template           # Update HTML template
GET    /admin/editor/schema      # Schema code editor interface
GET    /admin/editor/template    # Template code editor interface
```

### Public Endpoints

```
GET    /                         # Serves index.html
GET    /images/*                 # Static images
```

## Security Features

### Authentication
- Simple session-based authentication
- Environment variable for admin credentials
- Session timeout

### Input Validation
- Basic HTML sanitization
- Image file type validation
- File size limits

## File Structure

```
onepagems/
├── cmd/
│   └── main.go                  # Application entry point
├── internal/
│   ├── auth.go                  # Simple authentication
│   ├── handlers.go              # HTTP handlers
│   ├── content.go               # Content management
│   ├── schema.go                # Schema validation & management
│   └── generator.go             # HTML generation
├── templates/
│   ├── admin.html               # Admin panel template
│   ├── editor.html              # Code editor template
│   └── default_layout.html      # Default public page template
├── data/
│   ├── content.json             # Content storage
│   ├── content.json.bak         # Content backup
│   ├── schema.json              # JSON schema definition
│   ├── schema.json.bak          # Schema backup
│   ├── template.html            # Custom HTML template
│   ├── template.html.bak        # Template backup
│   └── images/                  # Uploaded images
├── static/
│   ├── admin.css                # Admin panel styles
│   ├── admin.js                 # Admin panel scripts
│   ├── codemirror/              # Code editor assets
│   │   ├── codemirror.css
│   │   ├── codemirror.js
│   │   ├── mode-html.js
│   │   └── mode-json.js
├── index.html                   # Generated public page
├── go.mod
├── README.md
└── DESIGN.md
```

## Configuration

### Environment Variables

```bash
# Server Configuration
PORT=8080
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your-password

# File Upload
UPLOAD_MAX_SIZE=5242880  # 5MB
```

## Content Management Features

### Admin Interface
- **Content Editor**: Dynamic form generation based on JSON schema
- **Schema Editor**: Web-based code editor for modifying JSON schema
- **Template Editor**: Web-based code editor for HTML template
- **Image Browser**: Upload and manage images
- **Real-time Preview**: Preview changes before saving

### Content Management Flow
1. **Schema Definition**: Admin can modify the JSON schema to define content structure
2. **Template Editing**: Admin can edit the HTML template with custom layout and styling
3. **Content Editing**: Dynamic forms are generated based on the current schema
4. **Image Management**: Upload and select images through the browser
5. **Generation**: Static HTML is generated using template + content + schema

### Schema-Driven Features
- **Dynamic Forms**: Content editing forms are automatically generated from schema
- **Field Types**: Support for text, textarea, email, image, array, object types
- **Validation**: Client and server-side validation based on schema rules
- **Field Metadata**: Titles, descriptions, defaults defined in schema

### Code Editors
- **JSON Schema Editor**: Syntax highlighting, validation, auto-completion
- **HTML Template Editor**: HTML syntax highlighting, Go template syntax support
- **Live Validation**: Real-time validation feedback
- **Error Highlighting**: Visual indicators for syntax errors

### Image Management
- Upload images via web interface
- Browse existing images with thumbnail preview
- Select images from browser for content fields
- Delete unused images
- Automatic file type validation

## Deployment

### Standalone Binary
- Single executable
- Self-contained with embedded templates
- Creates data directory on first run
- No external dependencies

## Development Plan

### Implementation Steps
1. **Basic Structure**: Set up Go project with basic HTTP server
2. **File Operations**: JSON content, schema, and template management
3. **Schema Engine**: JSON schema parsing and validation
4. **Template Engine**: HTML template processing with Go templates
5. **Dynamic Forms**: Schema-driven form generation
6. **Code Editors**: Web-based editors with syntax highlighting
7. **Image Management**: Upload and browse functionality
8. **Authentication**: Basic login/logout functionality

## Core Functionality

### Content Flow
1. Admin logs in via web interface
2. Admin can edit JSON schema to define content structure
3. Admin can edit HTML template for custom layout
4. Admin edits content through dynamically generated forms
5. Admin selects images via image browser
6. On save, content is validated against schema
7. Static `index.html` is generated from custom template
8. Previous versions are backed up to `.bak` files

### Schema-to-Form Generation
The system dynamically generates content editing forms based on the JSON schema:
- **String fields** → Text inputs
- **String with format:textarea** → Textarea inputs  
- **String with format:email** → Email inputs
- **String with format:image** → Image picker with browser
- **Array fields** → Dynamic add/remove item interface
- **Object fields** → Nested fieldsets
- **Field titles** → Form labels
- **Field defaults** → Initial values

### Template System
The HTML template has access to the complete content structure:
- Content is passed as JSON to the template
- Template uses Go template syntax: `{{.title}}`, `{{.sections.hero.title}}`
- Schema can define any structure, template adapts accordingly
- Custom CSS and JavaScript can be included in template

### Code Editor Features
- **Syntax Highlighting**: JSON and HTML syntax highlighting
- **Auto-completion**: Basic auto-completion for common patterns
- **Error Detection**: Real-time validation and error highlighting
- **Customizable**: Uses CodeMirror for extensible editing experience

---

## Conclusion

This enhanced design provides a fully configurable one-page CMS where both the content structure and presentation can be customized through web-based editors. The JSON schema defines the content fields and validation rules, while the HTML template controls the layout and styling. The dynamic form generation ensures that the admin interface adapts to any schema changes, making the system highly flexible while maintaining simplicity.
