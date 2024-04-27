package process

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"goconvert/internal/utils"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func converter(file *zip.File, inputFormat string, outputDir string) error {
	if strings.Contains(file.Name, "__MACOSX") {
		return nil
	}
	fileReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", file.Name, err)
	}
	defer fileReader.Close()

	// Define output file name
	inputFileName := strings.TrimSuffix(file.Name, filepath.Ext(file.Name))
	outputFileName := fmt.Sprintf("%s.%s", inputFileName, "md")
	docxFilePath := filepath.Join(outputDir, fmt.Sprintf("%s.docx", inputFileName))
	outputFilePath := filepath.Join(outputDir, outputFileName)

	if inputFormat == "pdf" {
		// Create a temporary file
		tempFile, err := os.CreateTemp("", "goconvert-*.pdf") // Customize the pattern based on expected file types
		if err != nil {
			return fmt.Errorf("error creating temporary file for %s: %v", file.Name, err)
		}
		defer tempFile.Close()
		// Copy the contents of the zip entry to the temporary file
		if _, err := io.Copy(tempFile, fileReader); err != nil {
			return fmt.Errorf("error writing to temporary file for %s: %v", file.Name, err)
		}
		// Ensure file pointer is at the start for any further operations
		_, err = tempFile.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("error seeking in temporary file for %s: %v", file.Name, err)
		}
		pdfToDocxCmd := exec.Command("python3", "../internal/process/pdftodocx/main.py", tempFile.Name(), docxFilePath)
		if output, err := pdfToDocxCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("error converting PDF to DOCX: %v, output: %s", err, output)
		}
		docxToMdCmd := exec.Command("pandoc", "-f", "docx", "-t", "markdown", "-o", outputFilePath, docxFilePath)
		if output, err := docxToMdCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("error converting DOCX to Markdown: %v, output: %s", err, output)
		}
	} else if inputFormat == "doc" {
		tempFile, err := os.CreateTemp("", fmt.Sprintf("%s.doc", inputFileName))
		if err != nil {
			return fmt.Errorf("error creating temporary file for %s: %v", file.Name, err)
		}
		defer tempFile.Close()
		if _, err := io.Copy(tempFile, fileReader); err != nil {
			return fmt.Errorf("error writing to temporary file for %s: %v", file.Name, err)
		}
		_, err = tempFile.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("error seeking in temporary file for %s: %v", file.Name, err)
		}
		docToDocxCmd := exec.Command("soffice", "--headless", "--convert-to", "docx", tempFile.Name(), "--outdir", outputDir)
		if output, err := docToDocxCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("error converting DOC to DOCX: %v, output: %s", err, output)
		}
		docxToMdCmd := exec.Command("pandoc", "-f", "docx", "-t", "markdown", "-o", outputFilePath, docxFilePath)
		if output, err := docxToMdCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("error converting DOCX to Markdown: %v, output: %s", err, output)
		}

	} else {
		convertCmd := exec.Command("pandoc", "-f", inputFormat, "-t", "markdown", "-o", outputFilePath)
		convertCmd.Stdin = fileReader
		if output, err := convertCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("error converting file %s to Markdown: %v, output: %s", file.Name, err, output)
		}
	}
	fmt.Printf("Converted file %s to %s successfully.\n", file.Name, outputFilePath)

	// Read the Markdown file content
	mdContent, err := os.ReadFile(outputFilePath)
	if err != nil {
		return fmt.Errorf("error reading converted file %s: %v", outputFileName, err)
	}
	// Sanitize the Markdown content
	content := strings.ReplaceAll(string(mdContent), "\r\n", "\n")
	replacer := strings.NewReplacer(">", "", "<", "", "^", "")
	sanitizedContent := replacer.Replace(content)

	// Extract title from the Markdown content
	var title string
	lines := strings.Split(sanitizedContent, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			title = strings.TrimPrefix(line, "# ")
			break
		}
	}
	if title == "" {
		for _, line := range lines {
			if strings.HasPrefix(line, "## ") {
				title = strings.TrimPrefix(line, "## ")
				break
			} else if strings.HasPrefix(line, "### ") {
				title = strings.TrimPrefix(line, "### ")
				break
			} else if strings.HasPrefix(line, "**") {
				title = strings.TrimPrefix(line, "**")
				break
			}
		}
	}
	if title == "" {
		title = inputFileName
	}

	// Create JSON payload
	payload := struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}{
		Title: title,
		Body:  string(sanitizedContent),
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
