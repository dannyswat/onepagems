package types

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
