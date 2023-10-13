package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	logger "github.com/farovictor/GCSUploader/src/logging"
)

// This Struct will collect all failing requests (file path reference)
type BatchCollector struct {
	mu                sync.Mutex
	FilesNotProcessed []string
}

// Mutex to track files (prevent data race)
func (b *BatchCollector) AddError(file string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.FilesNotProcessed = append(b.FilesNotProcessed, file)
}

// Dump the list
func (b *BatchCollector) DumpNotProcessedList(path string) error {
	if len(b.FilesNotProcessed) > 0 {
		filePath := fmt.Sprintf("%s/errors.log", path)

		content := strings.Join(b.FilesNotProcessed, "\n")

		err := ioutil.WriteFile(filePath, []byte(content), 0644)

		if err == nil {
			logger.InfoLogger.Println("List of not processed files is available in:", filePath)
		}
		return err
	}

	logger.InfoLogger.Println("No files presented error")

	return nil

}

// Emits filtered files to a channel
func EmitFilesToChannel(filePrefix string, searchPath string, emit chan<- string) (int, error) {
	total := 0

	// Walking through directory
	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.ErrorLogger.Println(err)
			return err
		}

		// To check if file does not have a regular mode
		if !info.Mode().IsRegular() {
			return nil
		}

		// Emit only files that match prefix
		if !info.IsDir() && strings.HasPrefix(info.Name(), filePrefix) {
			emit <- info.Name()
			total += 1
		}
		return nil
	})
	return total, err
}

func BlobNameAssemble(pathPrefix string, namePrefix string, filePath string) string {
	var blobName string

	if namePrefix != "" {
		blobName += fmt.Sprintf("%s-", namePrefix)
	}

	s := strings.Split(filePath, string(os.PathSeparator))
	fileName := s[len(s)-1]

	blobName += fileName

	return fmt.Sprintf("%s/%s", pathPrefix, blobName)
}
