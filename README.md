# OnePage CMS

A simple, lightweight, and secure content management system for managing single-page websites.

## Current Status: Step 2.2 Complete ✅

### Implemented Features:
- ✅ Basic HTTP server with routing
- ✅ Configuration management from environment variables
- ✅ Static file serving
- ✅ Directory structure auto-creation
- ✅ Health check endpoint
- ✅ Placeholder admin panel
- ✅ Core data types defined
- ✅ **File operations module (JSON/text read/write)**
- ✅ **Backup file management (.bak files)**
- ✅ **Directory creation and validation**
- ✅ **File existence checks and listing**
- ✅ **Atomic file operations (temp + rename)**
- ✅ **Template management system with validation**
- ✅ **Default responsive HTML template generation**
- ✅ **Content management (load/save/validate)**
- ✅ **Default content structure generation**
- ✅ **Content backup and restore functionality**
- ✅ **Content export/import (JSON)**

## Quick Start

1. **Run the application:**
   ```bash
   go run cmd/main.go
   ```

2. **Access the application:**
   - Public page: http://localhost:8080
   - Admin panel: http://localhost:8080/admin
   - Health check: http://localhost:8080/health

## Configuration

The application can be configured using environment variables:

```bash
# Server Configuration
export PORT=8080
export ADMIN_USERNAME=admin
export ADMIN_PASSWORD=your-secure-password

# File Upload
export UPLOAD_MAX_SIZE=5242880  # 5MB

# Session
export SESSION_TIMEOUT=60  # minutes

# Directories
export DATA_DIR=./data
export STATIC_DIR=./static
export TEMPLATES_DIR=./templates
```

## Current Endpoints

- `GET /` - Public page (placeholder)
- `GET /health` - Health check
- `GET /static/*` - Static files
- `GET /images/*` - Image files
- `GET /admin` - Admin panel with testing interface
- `POST /admin/login` - Login (placeholder)
- `POST /admin/logout` - Logout (placeholder)

### File Management
- `GET /admin/files` - List files (test endpoint)
- `POST /admin/test-storage` - Test storage operations

### Template Management
- `GET/POST /admin/template` - Template management
- `GET /admin/template/info` - Template information
- `POST /admin/template/restore` - Restore template from backup
- `POST /admin/test-template` - Test template operations

### Content Management
- `GET/POST /admin/content` - Content management
- `GET /admin/content/info` - Content information and summary
- `POST /admin/content/restore` - Restore content from backup
- `GET /admin/content/export` - Export content as JSON
- `POST /admin/content/import` - Import content from JSON
- `POST /admin/test-content` - Test content operations

## Testing the System

You can test the functionality using the built-in test endpoints:

```bash
# Test storage operations
curl -X POST http://localhost:8080/admin/test-storage

# Test template management
curl -X POST http://localhost:8080/admin/test-template

# Test content management
curl -X POST http://localhost:8080/admin/test-content

# Get content information
curl http://localhost:8080/admin/content/info

# Export current content
curl http://localhost:8080/admin/content/export

# List files in data directory
curl http://localhost:8080/admin/files
```

## Next Steps (Phase 3)

- Authentication system with session management
- Login/logout functionality
- Admin route protection
- CSRF protection

## Project Structure

```
onepagems/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config.go            # Configuration management
│   ├── server.go            # HTTP server implementation
│   └── types.go             # Core data types
├── templates/               # HTML templates (empty)
├── data/                    # Data storage
│   └── images/              # Image uploads
├── static/                  # Static assets
│   └── codemirror/          # Code editor assets (empty)
├── go.mod                   # Go module definition
├── DESIGN.md                # Design documentation
├── IMPLEMENT.md             # Implementation plan
└── README.md                # This file
```

## Development

This project follows a phased development approach as outlined in `IMPLEMENT.md`. Each phase builds upon the previous one while maintaining working software at each step.

**Current Phase**: Phase 2 Complete ✅  
**Next Phase**: Phase 3 - Authentication System
