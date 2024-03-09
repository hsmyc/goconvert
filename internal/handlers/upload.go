package handlers

import (
	"fmt"
	"net/http"

	"goconvert/internal/process"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UploadHandler")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	inputFormat := r.FormValue("inputFormat")
	outputFormat := r.FormValue("outputFormat")

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

	err = process.ProcessZipFile(file, r.ContentLength, inputFormat, outputFormat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded and processed successfully.")
}
