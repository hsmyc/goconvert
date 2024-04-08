package process

import (
	"archive/zip"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func converter(file *zip.File, outputFormat string, inputFormat string, outputDir string) error {
	fmt.Printf("Processing %s\n", inputFormat)
	fileReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", file.Name, err)
	}
	defer fileReader.Close()

	// Define output file name
	inputFileName := strings.TrimSuffix(file.Name, filepath.Ext(file.Name))
	outputFileName := fmt.Sprintf("%s.%s", inputFileName, outputFormat)
	outputFilePath := filepath.Join(outputDir, outputFileName)

	// Execute the conversion command
	cmd := exec.Command("pandoc", "-t", outputFormat, "-f", inputFormat, "-o", outputFilePath)
	cmd.Stdin = fileReader
	if _, err = cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error converting file %s: %v", file.Name, err)
	}

	fmt.Printf("Converted file %s to %s successfully.\n", file.Name, outputFilePath)
	return nil
}
