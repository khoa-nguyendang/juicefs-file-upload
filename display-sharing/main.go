package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// JuiceFS mount path - this will be mounted in container
const JUICEFS_MOUNT = "/jfs"

type FileInfo struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	IsDir    bool      `json:"is_dir"`
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
	Mode     string    `json:"mode"`
}

type DirectoryContent struct {
	Path  string     `json:"path"`
	Files []FileInfo `json:"files"`
	Total int        `json:"total"`
}

type FileContent struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Size    int64  `json:"size"`
	Binary  bool   `json:"binary"`
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// List files using POSIX filesystem operations
func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}

	fullPath := filepath.Join(JUICEFS_MOUNT, path)

	// Ensure path doesn't escape mount point
	if !strings.HasPrefix(fullPath, JUICEFS_MOUNT) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	log.Printf("Listing files in JuiceFS path: %s", fullPath)

	// Read directory using standard Go filesystem operations
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		log.Printf("Error reading directory: %v", err)
		http.Error(w, fmt.Sprintf("Error reading directory: %v", err), http.StatusInternalServerError)
		return
	}

	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			log.Printf("Error getting file info for %s: %v", entry.Name(), err)
			continue
		}

		relativePath := filepath.Join(path, entry.Name())
		if !strings.HasPrefix(relativePath, "/") {
			relativePath = "/" + relativePath
		}

		files = append(files, FileInfo{
			Name:     entry.Name(),
			Path:     relativePath,
			IsDir:    entry.IsDir(),
			Size:     info.Size(),
			Modified: info.ModTime(),
			Mode:     info.Mode().String(),
		})
	}

	response := DirectoryContent{
		Path:  path,
		Files: files,
		Total: len(files),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// View file content using POSIX read
func viewFileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(JUICEFS_MOUNT, filePath)

	// Security: Ensure path doesn't escape mount point
	if !strings.HasPrefix(fullPath, JUICEFS_MOUNT) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	log.Printf("Viewing file from JuiceFS: %s", fullPath)

	// Check if file exists and get info
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error accessing file: %v", err), http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		http.Error(w, "Path is a directory, not a file", http.StatusBadRequest)
		return
	}

	// Read file content
	content, err := os.ReadFile(fullPath)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if file is binary
	isBinary := isBinaryContent(content)

	response := FileContent{
		Path:   filePath,
		Size:   info.Size(),
		Binary: isBinary,
	}

	// Only include content for text files under 10MB
	if !isBinary && info.Size() < 10*1024*1024 {
		response.Content = string(content)
	} else if isBinary {
		response.Content = fmt.Sprintf("[Binary file, size: %d bytes]", info.Size())
	} else {
		response.Content = fmt.Sprintf("[File too large to display: %d bytes]", info.Size())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Delete file using POSIX operations
func deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(JUICEFS_MOUNT, filePath)

	// Security: Ensure path doesn't escape mount point
	if !strings.HasPrefix(fullPath, JUICEFS_MOUNT) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	log.Printf("Deleting from JuiceFS: %s", fullPath)

	// Check if path exists
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error accessing file: %v", err), http.StatusInternalServerError)
		return
	}

	// Delete file or directory
	if info.IsDir() {
		err = os.RemoveAll(fullPath)
	} else {
		err = os.Remove(fullPath)
	}

	if err != nil {
		log.Printf("Error deleting: %v", err)
		http.Error(w, fmt.Sprintf("Error deleting: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Deleted: %s", filePath),
		"path":    filePath,
	})
}

// Download file using POSIX read
func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(JUICEFS_MOUNT, filePath)

	if !strings.HasPrefix(fullPath, JUICEFS_MOUNT) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Open file using standard Go file operations
	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error opening file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get file info
	info, err := file.Stat()
	if err != nil {
		http.Error(w, "Error getting file info", http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		http.Error(w, "Cannot download directory", http.StatusBadRequest)
		return
	}

	// Set headers for download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filePath)))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))

	// Stream file content
	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("Error streaming file: %v", err)
	}
}

// Health check that verifies JuiceFS mount
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Check if JuiceFS mount point exists and is accessible
	mountInfo, err := os.Stat(JUICEFS_MOUNT)

	status := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"mount":     JUICEFS_MOUNT,
	}

	if err != nil {
		status["status"] = "unhealthy"
		status["error"] = fmt.Sprintf("Mount point not accessible: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		status["mount_exists"] = true
		status["mount_is_dir"] = mountInfo.IsDir()

		// Try to list files to verify mount is working
		entries, err := os.ReadDir(JUICEFS_MOUNT)
		if err != nil {
			status["status"] = "degraded"
			status["error"] = fmt.Sprintf("Cannot read mount: %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			status["file_count"] = len(entries)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// Stats handler to show JuiceFS mount statistics
func statsHandler(w http.ResponseWriter, r *http.Request) {
	var totalSize int64
	var fileCount int
	var dirCount int

	// Walk through JuiceFS mount to calculate stats
	err := filepath.WalkDir(JUICEFS_MOUNT, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if d.IsDir() {
			dirCount++
		} else {
			fileCount++
			info, err := d.Info()
			if err == nil {
				totalSize += info.Size()
			}
		}
		return nil
	})

	stats := map[string]interface{}{
		"mount_point":  JUICEFS_MOUNT,
		"total_files":  fileCount,
		"total_dirs":   dirCount,
		"total_size":   totalSize,
		"size_formatted": formatBytes(totalSize),
		"timestamp":    time.Now().Format(time.RFC3339),
	}

	if err != nil {
		stats["error"] = fmt.Sprintf("Error calculating stats: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Helper functions
func isBinaryContent(content []byte) bool {
	// Simple check for binary content
	for _, b := range content[:min(len(content), 512)] {
		if b == 0 {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func main() {
	// Check if JuiceFS mount exists
	if _, err := os.Stat(JUICEFS_MOUNT); err != nil {
		log.Printf("WARNING: JuiceFS mount point %s not found: %v", JUICEFS_MOUNT, err)
		log.Printf("Make sure JuiceFS is mounted at %s", JUICEFS_MOUNT)
	} else {
		log.Printf("JuiceFS mount point found at %s", JUICEFS_MOUNT)
	}

	// Routes for POSIX file operations
	http.HandleFunc("/api/list", corsMiddleware(listFilesHandler))
	http.HandleFunc("/api/view", corsMiddleware(viewFileHandler))
	http.HandleFunc("/api/delete", corsMiddleware(deleteFileHandler))
	http.HandleFunc("/api/download", corsMiddleware(downloadFileHandler))
	http.HandleFunc("/api/health", corsMiddleware(healthHandler))
	http.HandleFunc("/api/stats", corsMiddleware(statsHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	log.Printf("Display Sharing Service (JuiceFS POSIX) starting on port %s", port)
	log.Printf("Using JuiceFS mount at: %s", JUICEFS_MOUNT)
	log.Printf("This service demonstrates POSIX filesystem operations via JuiceFS")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}