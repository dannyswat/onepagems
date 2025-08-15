package main

import (
	"log"
	"onepagems/internal"
)

func main() {
	// Load configuration from environment variables
	config := internal.LoadConfig()

	// Validate configuration
	if err := internal.ValidateConfig(config); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Log configuration (without sensitive data)
	log.Printf("OnePage CMS starting with configuration:")
	log.Printf("  Port: %s", config.Port)
	log.Printf("  Data directory: %s", config.DataDir)
	log.Printf("  Static directory: %s", config.StaticDir)
	log.Printf("  Templates directory: %s", config.TemplatesDir)
	log.Printf("  Upload max size: %d bytes", config.UploadMaxSize)
	log.Printf("  Session timeout: %d minutes", config.SessionTimeout)
	log.Printf("  Admin username: %s", config.AdminUsername)

	// Create and start server
	server := internal.NewServer(config)

	log.Println("OnePage CMS server starting...")
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
