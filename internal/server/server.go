package server

import (
	"fmt"
	"log"
	"net/http"
	"onepagems/internal/managers"
	"onepagems/internal/types"
	"os"
)

// Server represents the HTTP server
type Server struct {
	Config          *types.Config
	Storage         *managers.FileStorage
	TemplateManager *managers.TemplateManager
	ContentManager  *managers.ContentManager
	SchemaManager   *managers.SchemaManager
	AuthManager     *managers.AuthManager
	Mux             *http.ServeMux
}

// NewServer creates a new server instance
func NewServer(config *types.Config) *Server {
	storage := managers.NewFileStorage(config.DataDir)
	server := &Server{
		Config:          config,
		Storage:         storage,
		TemplateManager: managers.NewTemplateManager(storage),
		ContentManager:  managers.NewContentManager(storage, config.DataDir),
		SchemaManager:   managers.NewSchemaManager(storage, config.DataDir),
		AuthManager:     managers.NewAuthManager(config),
		Mux:             http.NewServeMux(),
	}

	// Set up routes
	server.setupRoutes()

	return server
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Ensure data directories exist
	if err := s.ensureDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	addr := ":" + s.Config.Port
	log.Printf("Starting server on http://localhost%s", addr)
	log.Printf("Admin panel: http://localhost%s/admin", addr)

	return http.ListenAndServe(addr, s.Mux)
}

// ensureDirectories creates necessary directories if they don't exist
func (s *Server) ensureDirectories() error {
	// Use storage to ensure data directories
	if err := s.Storage.EnsureDirectories(); err != nil {
		return err
	}

	// Ensure other directories
	dirs := []string{
		s.Config.StaticDir,
		s.Config.TemplatesDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		log.Printf("Ensured directory exists: %s", dir)
	}

	return nil
}
