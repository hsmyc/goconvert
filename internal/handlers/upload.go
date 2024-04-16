package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"goconvert/internal/process"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	outputDir, err := os.MkdirTemp("", "converted")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(outputDir) // Clean up

	err = process.ProcessZipFile(file, r.ContentLength, outputDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	zipFilePath := filepath.Join(os.TempDir(), "output.zip")
	err = process.ZipOutputDirectory(outputDir, zipFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(zipFilePath)
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=output.zip")
	http.ServeFile(w, r, zipFilePath)
}
