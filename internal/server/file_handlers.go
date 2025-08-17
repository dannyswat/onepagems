package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"onepagems/internal/types"
)

// handleFilesList lists all files in the data directory
func (s *Server) handleFilesList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	files, err := s.Storage.ListFiles()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list files: %v", err), http.StatusInternalServerError)
		return
	}

	response := types.NewAPIResponse(true, "Files listed successfully")
	response.SetData(files)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleTestStorage demonstrates file storage operations
func (s *Server) handleTestStorage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Test data
	testData := map[string]interface{}{
		"message":    "Hello from file storage test",
		"timestamp":  time.Now().Format(time.RFC3339),
		"test_array": []string{"item1", "item2", "item3"},
		"test_object": map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	// Test JSON operations
	filename := "test.json"

	// Write test file
	if err := s.Storage.WriteJSONFile(filename, testData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write test file: %v", err), http.StatusInternalServerError)
		return
	}

	// Read it back
	var readData map[string]interface{}
	if err := s.Storage.ReadJSONFile(filename, &readData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to read test file: %v", err), http.StatusInternalServerError)
		return
	}

	// Get backup info
	backupInfo, _ := s.Storage.GetBackupInfo(filename)

	// Test text file operations
	textFilename := "test.txt"
	textContent := fmt.Sprintf("Test text file created at %s\nThis demonstrates text file operations.", time.Now().Format(time.RFC3339))

	if err := s.Storage.WriteTextFile(textFilename, textContent); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write text file: %v", err), http.StatusInternalServerError)
		return
	}

	readText, err := s.Storage.ReadTextFile(textFilename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read text file: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare response
	result := map[string]interface{}{
		"json_write_success":  true,
		"json_read_success":   true,
		"json_data_matches":   fmt.Sprintf("%v", testData["message"]) == fmt.Sprintf("%v", readData["message"]),
		"text_write_success":  true,
		"text_read_success":   true,
		"text_content_length": len(readText),
		"backup_created":      backupInfo != nil,
		"files_created":       []string{filename, textFilename},
	}

	if backupInfo != nil {
		result["backup_info"] = map[string]interface{}{
			"backup_path": backupInfo.BackupPath,
			"created_at":  backupInfo.CreatedAt,
			"size":        backupInfo.Size,
		}
	}

	response := types.NewAPIResponse(true, "File storage test completed successfully")
	response.SetData(result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
