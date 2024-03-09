package process

import (
	"archive/zip"
	"fmt"
	"io"
	"strings"
	"sync"
)

func ProcessZipFile(zipFile io.ReaderAt, size int64, i string, o string) error {
	zipReader, err := zip.NewReader(zipFile, size)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, zipFile := range zipReader.File {
		wg.Add(1)
		go func(file *zip.File) {
			defer wg.Done()
			fmt.Printf("Processing %s\n", file.Name)
			parts := strings.Split(file.Name, ".")
			fileFormat := parts[len(parts)-1]
			err := converter(file, o, fileFormat)
			if err != nil {
				fmt.Printf("Error converting file %s: %v\n", file.Name, err)
				return
			}
		}(zipFile)
	}
	wg.Wait()
	return nil
}
