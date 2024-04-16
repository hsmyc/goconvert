package process

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"sync"
)

func ProcessZipFile(zipFile io.ReaderAt, size int64, outputDir string) error {
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
			if len(ext) <= 1 {
				fmt.Printf("No valid extension for file %s\n", file.Name)
				return
			}
			fileFormat := filepath.Ext(file.Name)[1:]
			err := converter(file, fileFormat, outputDir)
			if err != nil {
				fmt.Printf("Error converting file %s: %v\n", file.Name, err)
				return
			}
		}(file)
	}
	wg.Wait()
	return nil
}
