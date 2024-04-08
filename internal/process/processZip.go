package process

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"sync"
)

// ProcessZipFile processes each file within a zip archive for conversion.
func ProcessZipFile(zipFile io.ReaderAt, size int64, outputFormat string, outputDir string) error {
	zipReader, err := zip.NewReader(zipFile, size)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, file := range zipReader.File {
		wg.Add(1)
		go func(file *zip.File) {
			defer wg.Done()
			fmt.Printf("Processing %s\n", file.Name)
			ext := filepath.Ext(file.Name)
			if len(ext) <= 1 { // Check if the extension is empty or only contains the period
				fmt.Printf("No valid extension for file %s\n", file.Name)
				return
			}
			fileFormat := filepath.Ext(file.Name)[1:]
			err := converter(file, outputFormat, fileFormat, outputDir)
			if err != nil {
				fmt.Printf("Error converting file %s: %v\n", file.Name, err)
				return
			}
		}(file)
	}
	wg.Wait()
	return nil
}
