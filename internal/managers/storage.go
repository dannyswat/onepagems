package managers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"onepagems/internal/types"
)

// FileStorage handles all file operations for the CMS
type FileStorage struct {
	dataDir string
}

// NewFileStorage creates a new file storage instance
func NewFileStorage(dataDir string) *FileStorage {
	return &FileStorage{
		dataDir: dataDir,
	}
}

// EnsureDirectories creates all necessary directories if they don't exist
func (fs *FileStorage) EnsureDirectories() error {
	dirs := []string{
		fs.dataDir,
		filepath.Join(fs.dataDir, "images"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// FileExists checks if a file exists
func (fs *FileStorage) FileExists(filename string) bool {
	fullPath := filepath.Join(fs.dataDir, filename)
	_, err := os.Stat(fullPath)
	return err == nil
}

// GetFilePath returns the full path for a file in the data directory
func (fs *FileStorage) GetFilePath(filename string) string {
	return filepath.Join(fs.dataDir, filename)
}

// ReadJSONFile reads and unmarshals a JSON file
func (fs *FileStorage) ReadJSONFile(filename string, target interface{}) error {
	fullPath := fs.GetFilePath(filename)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", filename)
		}
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %w", filename, err)
	}

	return nil
}

// WriteJSONFile marshals and writes data to a JSON file
func (fs *FileStorage) WriteJSONFile(filename string, data interface{}) error {
	// Create backup before writing
	if err := fs.CreateBackup(filename); err != nil {
		// Log the error but don't fail the write operation
		fmt.Printf("Warning: failed to create backup for %s: %v\n", filename, err)
	}

	fullPath := fs.GetFilePath(filename)

	// Marshal data with indentation for readability
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data for %s: %w", filename, err)
	}

	// Write to temporary file first, then rename (atomic operation)
	tempPath := fullPath + ".tmp"
	if err := os.WriteFile(tempPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file %s: %w", tempPath, err)
	}

	// Rename temporary file to final file (atomic on most filesystems)
	if err := os.Rename(tempPath, fullPath); err != nil {
		// Clean up temporary file on failure
		os.Remove(tempPath)
		return fmt.Errorf("failed to rename temporary file %s to %s: %w", tempPath, fullPath, err)
	}

	return nil
}

// ReadTextFile reads a text file and returns its contents
func (fs *FileStorage) ReadTextFile(filename string) (string, error) {
	fullPath := fs.GetFilePath(filename)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file %s does not exist", filename)
		}
		return "", fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return string(data), nil
}

// WriteTextFile writes text content to a file
func (fs *FileStorage) WriteTextFile(filename string, content string) error {
	// Create backup before writing
	if err := fs.CreateBackup(filename); err != nil {
		// Log the error but don't fail the write operation
		fmt.Printf("Warning: failed to create backup for %s: %v\n", filename, err)
	}

	fullPath := fs.GetFilePath(filename)

	// Write to temporary file first, then rename (atomic operation)
	tempPath := fullPath + ".tmp"
	if err := os.WriteFile(tempPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write temporary file %s: %w", tempPath, err)
	}

	// Rename temporary file to final file (atomic on most filesystems)
	if err := os.Rename(tempPath, fullPath); err != nil {
		// Clean up temporary file on failure
		os.Remove(tempPath)
		return fmt.Errorf("failed to rename temporary file %s to %s: %w", tempPath, fullPath, err)
	}

	return nil
}

// CreateBackup creates a backup of a file with .bak extension
func (fs *FileStorage) CreateBackup(filename string) error {
	sourcePath := fs.GetFilePath(filename)
	backupPath := sourcePath + ".bak"

	// Check if source file exists
	if !fs.FileExists(filename) {
		// No file to backup, which is fine
		return nil
	}

	// Open source file
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", sourcePath, err)
	}
	defer sourceFile.Close()

	// Create backup file
	backupFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file %s: %w", backupPath, err)
	}
	defer backupFile.Close()

	// Copy contents
	if _, err := io.Copy(backupFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy contents to backup file %s: %w", backupPath, err)
	}

	return nil
}

// RestoreFromBackup restores a file from its backup
func (fs *FileStorage) RestoreFromBackup(filename string) error {
	sourcePath := fs.GetFilePath(filename)
	backupPath := sourcePath + ".bak"

	// Check if backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file %s does not exist", backupPath)
	}

	// Open backup file
	backupFile, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file %s: %w", backupPath, err)
	}
	defer backupFile.Close()

	// Create/overwrite main file
	mainFile, err := os.Create(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to create main file %s: %w", sourcePath, err)
	}
	defer mainFile.Close()

	// Copy contents
	if _, err := io.Copy(mainFile, backupFile); err != nil {
		return fmt.Errorf("failed to copy contents from backup to main file: %w", err)
	}

	return nil
}

// GetBackupInfo returns information about a backup file
func (fs *FileStorage) GetBackupInfo(filename string) (*types.FileBackup, error) {
	sourcePath := fs.GetFilePath(filename)
	backupPath := sourcePath + ".bak"

	// Check if backup exists
	info, err := os.Stat(backupPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("backup file does not exist")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get backup file info: %w", err)
	}

	return &types.FileBackup{
		OriginalPath: sourcePath,
		BackupPath:   backupPath,
		CreatedAt:    info.ModTime(),
		Size:         info.Size(),
	}, nil
}

// ListFiles returns a list of files in the data directory with their info
func (fs *FileStorage) ListFiles() ([]types.FileInfo, error) {
	entries, err := os.ReadDir(fs.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	var files []types.FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Skip backup files from main listing
		if filepath.Ext(entry.Name()) == ".bak" {
			continue
		}

		// Check if backup exists
		backupPath := filepath.Join(fs.dataDir, entry.Name()+".bak")
		hasBackup := false
		var backupAge *int64
		if backupInfo, err := os.Stat(backupPath); err == nil {
			hasBackup = true
			age := int64(time.Since(backupInfo.ModTime()).Seconds())
			backupAge = &age
		}

		fileInfo := types.FileInfo{
			Path:        filepath.Join(fs.dataDir, entry.Name()),
			Name:        entry.Name(),
			Size:        info.Size(),
			ModifiedAt:  info.ModTime(),
			IsDirectory: false,
			HasBackup:   hasBackup,
			BackupAge:   backupAge,
		}

		// Determine content type based on extension
		switch filepath.Ext(entry.Name()) {
		case ".json":
			fileInfo.ContentType = "application/json"
		case ".html":
			fileInfo.ContentType = "text/html"
		case ".txt":
			fileInfo.ContentType = "text/plain"
		default:
			fileInfo.ContentType = "application/octet-stream"
		}

		files = append(files, fileInfo)
	}

	return files, nil
}

// DeleteFile deletes a file and its backup if it exists
func (fs *FileStorage) DeleteFile(filename string) error {
	sourcePath := fs.GetFilePath(filename)
	backupPath := sourcePath + ".bak"

	// Delete main file
	if err := os.Remove(sourcePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete main file %s: %w", sourcePath, err)
	}

	// Delete backup file if it exists
	if _, err := os.Stat(backupPath); err == nil {
		if err := os.Remove(backupPath); err != nil {
			// Log warning but don't fail
			fmt.Printf("Warning: failed to delete backup file %s: %v\n", backupPath, err)
		}
	}

	return nil
}

// ValidateJSON checks if a string contains valid JSON
func (fs *FileStorage) ValidateJSON(data string) error {
	var temp interface{}
	if err := json.Unmarshal([]byte(data), &temp); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

// GetFileSize returns the size of a file in bytes
func (fs *FileStorage) GetFileSize(filename string) (int64, error) {
	fullPath := fs.GetFilePath(filename)
	info, err := os.Stat(fullPath)
	if err != nil {
		return 0, fmt.Errorf("failed to get file size for %s: %w", filename, err)
	}
	return info.Size(), nil
}

// GetFileModTime returns the modification time of a file
func (fs *FileStorage) GetFileModTime(filename string) (time.Time, error) {
	fullPath := fs.GetFilePath(filename)
	info, err := os.Stat(fullPath)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get file modification time for %s: %w", filename, err)
	}
	return info.ModTime(), nil
}
