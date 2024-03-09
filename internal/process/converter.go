package process

import (
	"archive/zip"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func converter(file *zip.File, o string, i string) error {
	outDir := "outdir"
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %v", outDir, err)
	}

	fileReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", file.Name, err)
	}
	defer fileReader.Close()

	inputFileName := strings.Split(file.Name, ".")[0]

	outputFileName := fmt.Sprintf("%s/%s.%s", outDir, inputFileName, o)

	cmd := exec.Command("pandoc", "-t", o, "-f", i, "-o", outputFileName)
	cmd.Stdin = fileReader

	if _, err = cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error converting file %s: %v", file.Name, err)
	}

	fmt.Printf("Converted file %s to %s successfully.\n", file.Name, outputFileName)
	return nil
}
