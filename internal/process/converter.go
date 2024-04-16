package process

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"goconvert/internal/utils"
)

func converter(file *zip.File, inputFormat string, outputDir string) error {
	fileReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", file.Name, err)
	}
	defer fileReader.Close()

	// Define output file name
	inputFileName := strings.TrimSuffix(file.Name, filepath.Ext(file.Name))
	outputFileName := fmt.Sprintf("%s.%s", inputFileName, "md")
	outputFilePath := filepath.Join(outputDir, outputFileName)

	// Execute the conversion command
	cmd := exec.Command("pandoc", "-t", "markdown", "-f", inputFormat, "-o", outputFilePath)
	cmd.Stdin = fileReader
	if _, err = cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error converting file %s: %v", file.Name, err)
	}

	fmt.Printf("Converted file %s to %s successfully.\n", file.Name, outputFilePath)

	// Read the Markdown file content
	mdContent, err := os.ReadFile(outputFilePath)
	if err != nil {
		return fmt.Errorf("error reading converted file %s: %v", outputFileName, err)
	}

	// Create JSON payload
	payload := struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}{
		Title: inputFileName,
		Body:  string(mdContent),
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Make the POST request
	url, token := utils.GetStrapi()

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token) // Set your token if required

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(response.Body)
		return fmt.Errorf("received non-OK response status: %d, body: %s", response.StatusCode, responseBody)
	}

	fmt.Println("POST request to Strapi successful")
	return nil
}
