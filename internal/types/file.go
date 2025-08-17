package types

import "time"

// ImageInfo represents information about an uploaded image
type ImageInfo struct {
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
	URL          string    `json:"url"`
}

// FileBackup represents backup file information
type FileBackup struct {
	OriginalPath string    `json:"original_path"`
	BackupPath   string    `json:"backup_path"`
	CreatedAt    time.Time `json:"created_at"`
	Size         int64     `json:"size"`
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
