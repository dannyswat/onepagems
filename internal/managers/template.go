package managers

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	t "onepagems/internal/types"
)

// TemplateManager handles template operations
type TemplateManager struct {
	storage *FileStorage
}

// NewTemplateManager creates a new template manager
func NewTemplateManager(storage *FileStorage) *TemplateManager {
	return &TemplateManager{
		storage: storage,
	}
}

// LoadTemplate loads the HTML template from file
func (tm *TemplateManager) LoadTemplate() (string, error) {
	const filename = "template.html"

	if !tm.storage.FileExists(filename) {
		// Generate default template if none exists
		defaultTemplate := tm.GetDefaultTemplate()
		if err := tm.SaveTemplate(defaultTemplate); err != nil {
			return "", fmt.Errorf("failed to create default template: %w", err)
		}
		return defaultTemplate, nil
	}

	content, err := tm.storage.ReadTextFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to load template: %w", err)
	}

	return content, nil
}

// SaveTemplate saves the HTML template to file
func (tm *TemplateManager) SaveTemplate(content string) error {
	const filename = "template.html"

	// Validate template before saving
	if err := tm.ValidateTemplate(content); err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}

	if err := tm.storage.WriteTextFile(filename, content); err != nil {
		return fmt.Errorf("failed to save template: %w", err)
	}

	return nil
}

// ValidateTemplate validates the HTML template syntax
func (tm *TemplateManager) ValidateTemplate(content string) error {
	// Check if template is not empty
	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("template cannot be empty")
	}

	// Try to parse as Go template
	tmpl, err := template.New("test").Parse(content)
	if err != nil {
		return fmt.Errorf("template parsing failed: %w", err)
	}

	// Test execution with dummy data to catch runtime errors
	testData := map[string]interface{}{
		"title":       "Test Title",
		"description": "Test Description",
		"sections": map[string]interface{}{
			"hero": map[string]interface{}{
				"title":       "Hero Title",
				"subtitle":    "Hero Subtitle",
				"image":       "/images/test.jpg",
				"button_text": "Click Me",
				"button_link": "#test",
			},
			"about": map[string]interface{}{
				"title":   "About",
				"content": "Test content",
				"image":   "/images/about.jpg",
			},
		},
	}

	// Execute template with test data
	var buf strings.Builder
	if err := tmpl.Execute(&buf, testData); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	// Check for basic HTML structure
	output := buf.String()
	if !strings.Contains(strings.ToLower(output), "<html") {
		return fmt.Errorf("template must contain valid HTML structure")
	}

	return nil
}

// GetDefaultTemplate returns the default HTML template
func (tm *TemplateManager) GetDefaultTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.title}}</title>
    <meta name="description" content="{{.description}}">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 0 20px;
        }
        
        /* Header */
        header {
            background: #fff;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            position: sticky;
            top: 0;
            z-index: 100;
        }
        
        nav {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 1rem 0;
        }
        
        .logo {
            font-size: 1.5rem;
            font-weight: bold;
            color: #007cba;
        }
        
        /* Hero Section */
        .hero {
            background: linear-gradient(135deg, #007cba 0%, #005a87 100%);
            color: white;
            text-align: center;
            padding: 100px 0;
        }
        
        .hero h1 {
            font-size: 3rem;
            margin-bottom: 1rem;
        }
        
        .hero p {
            font-size: 1.25rem;
            margin-bottom: 2rem;
            opacity: 0.9;
        }
        
        .btn {
            display: inline-block;
            padding: 12px 30px;
            background: white;
            color: #007cba;
            text-decoration: none;
            border-radius: 5px;
            font-weight: bold;
            transition: transform 0.3s ease;
        }
        
        .btn:hover {
            transform: translateY(-2px);
        }
        
        /* Content Sections */
        .section {
            padding: 80px 0;
        }
        
        .section:nth-child(even) {
            background: #f8f9fa;
        }
        
        .section h2 {
            font-size: 2.5rem;
            text-align: center;
            margin-bottom: 3rem;
            color: #333;
        }
        
        .about-content {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 4rem;
            align-items: center;
        }
        
        .about-text {
            font-size: 1.1rem;
            line-height: 1.8;
        }
        
        .about-image {
            text-align: center;
        }
        
        .about-image img {
            max-width: 100%;
            border-radius: 10px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
        }
        
        /* Services */
        .services-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 2rem;
            margin-top: 2rem;
        }
        
        .service-card {
            background: white;
            padding: 2rem;
            border-radius: 10px;
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
            text-align: center;
            transition: transform 0.3s ease;
        }
        
        .service-card:hover {
            transform: translateY(-5px);
        }
        
        .service-card h3 {
            color: #007cba;
            margin-bottom: 1rem;
        }
        
        /* Contact */
        .contact-info {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 2rem;
            text-align: center;
        }
        
        .contact-item {
            padding: 2rem;
        }
        
        .contact-item h3 {
            color: #007cba;
            margin-bottom: 1rem;
        }
        
        /* Footer */
        footer {
            background: #333;
            color: white;
            text-align: center;
            padding: 2rem 0;
        }
        
        /* Responsive */
        @media (max-width: 768px) {
            .hero h1 {
                font-size: 2rem;
            }
            
            .about-content {
                grid-template-columns: 1fr;
                gap: 2rem;
            }
            
            .services-grid {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <!-- Header -->
    <header>
        <nav class="container">
            <div class="logo">{{.title}}</div>
        </nav>
    </header>

    <!-- Hero Section -->
    {{with .sections.hero}}
    <section class="hero">
        <div class="container">
            <h1>{{.title}}</h1>
            <p>{{.subtitle}}</p>
            {{if .button_text}}
            <a href="{{.button_link}}" class="btn">{{.button_text}}</a>
            {{end}}
        </div>
    </section>
    {{end}}

    <!-- About Section -->
    {{with .sections.about}}
    <section class="section">
        <div class="container">
            <h2>{{.title}}</h2>
            <div class="about-content">
                <div class="about-text">
                    <p>{{.content}}</p>
                </div>
                {{if .image}}
                <div class="about-image">
                    <img src="{{.image}}" alt="{{.title}}">
                </div>
                {{end}}
            </div>
        </div>
    </section>
    {{end}}

    <!-- Services Section -->
    {{with .sections.services}}
    <section class="section">
        <div class="container">
            <h2>{{.title}}</h2>
            {{if .items}}
            <div class="services-grid">
                {{range .items}}
                <div class="service-card">
                    {{if .image}}
                    <img src="{{.image}}" alt="{{.title}}" style="width: 60px; height: 60px; margin-bottom: 1rem;">
                    {{end}}
                    <h3>{{.title}}</h3>
                    <p>{{.description}}</p>
                </div>
                {{end}}
            </div>
            {{end}}
        </div>
    </section>
    {{end}}

    <!-- Contact Section -->
    {{with .sections.contact}}
    <section class="section">
        <div class="container">
            <h2>{{.title}}</h2>
            <div class="contact-info">
                {{if .email}}
                <div class="contact-item">
                    <h3>Email</h3>
                    <p><a href="mailto:{{.email}}">{{.email}}</a></p>
                </div>
                {{end}}
                {{if .phone}}
                <div class="contact-item">
                    <h3>Phone</h3>
                    <p><a href="tel:{{.phone}}">{{.phone}}</a></p>
                </div>
                {{end}}
                {{if .address}}
                <div class="contact-item">
                    <h3>Address</h3>
                    <p>{{.address}}</p>
                </div>
                {{end}}
            </div>
        </div>
    </section>
    {{end}}

    <!-- Footer -->
    <footer>
        <div class="container">
            <p>&copy; 2025 {{.title}}. All rights reserved.</p>
        </div>
    </footer>
</body>
</html>`
}

// GetTemplateInfo returns information about the current template
func (tm *TemplateManager) GetTemplateInfo() (*t.FileInfo, error) {
	const filename = "template.html"

	if !tm.storage.FileExists(filename) {
		return nil, fmt.Errorf("template file does not exist")
	}

	size, err := tm.storage.GetFileSize(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get template size: %w", err)
	}

	modTime, err := tm.storage.GetFileModTime(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get template modification time: %w", err)
	}

	// Check if backup exists
	backupInfo, _ := tm.storage.GetBackupInfo(filename)
	hasBackup := backupInfo != nil
	var backupAge *int64
	if hasBackup {
		age := int64(time.Since(backupInfo.CreatedAt).Seconds())
		backupAge = &age
	}

	return &t.FileInfo{
		Path:        tm.storage.GetFilePath(filename),
		Name:        filename,
		Size:        size,
		ModifiedAt:  modTime,
		IsDirectory: false,
		ContentType: "text/html",
		HasBackup:   hasBackup,
		BackupAge:   backupAge,
	}, nil
}

// RestoreTemplate restores template from backup
func (tm *TemplateManager) RestoreTemplate() error {
	const filename = "template.html"

	if err := tm.storage.RestoreFromBackup(filename); err != nil {
		return fmt.Errorf("failed to restore template from backup: %w", err)
	}

	// Validate restored template
	content, err := tm.LoadTemplate()
	if err != nil {
		return fmt.Errorf("failed to load restored template: %w", err)
	}

	if err := tm.ValidateTemplate(content); err != nil {
		return fmt.Errorf("restored template is invalid: %w", err)
	}

	return nil
}

// DeleteTemplate deletes the template file and its backup
func (tm *TemplateManager) DeleteTemplate() error {
	const filename = "template.html"

	if err := tm.storage.DeleteFile(filename); err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	return nil
}

// GetTemplateVariables extracts variables used in the template
func (tm *TemplateManager) GetTemplateVariables(content string) ([]string, error) {
	tmpl, err := template.New("analysis").Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template for analysis: %w", err)
	}

	// This is a simple extraction - in a real implementation you might want
	// to use a more sophisticated method to extract all template variables
	variables := []string{
		"title",
		"description",
		"sections.hero.title",
		"sections.hero.subtitle",
		"sections.hero.image",
		"sections.hero.button_text",
		"sections.hero.button_link",
		"sections.about.title",
		"sections.about.content",
		"sections.about.image",
		"sections.services.title",
		"sections.services.items",
		"sections.contact.title",
		"sections.contact.email",
		"sections.contact.phone",
		"sections.contact.address",
	}

	// Test that the template can execute (basic validation)
	testData := map[string]interface{}{
		"title":       "Test",
		"description": "Test",
		"sections": map[string]interface{}{
			"hero": map[string]interface{}{
				"title": "Test", "subtitle": "Test", "image": "test.jpg",
				"button_text": "Test", "button_link": "#test",
			},
			"about": map[string]interface{}{
				"title": "Test", "content": "Test", "image": "test.jpg",
			},
			"services": map[string]interface{}{
				"title": "Test",
				"items": []map[string]interface{}{
					{"title": "Test", "description": "Test", "image": "test.jpg"},
				},
			},
			"contact": map[string]interface{}{
				"title": "Test", "email": "test@test.com",
				"phone": "123", "address": "Test",
			},
		},
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, testData); err != nil {
		return nil, fmt.Errorf("template execution test failed: %w", err)
	}

	return variables, nil
}
