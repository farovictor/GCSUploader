package main

import (
	"context"
	"os"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/farovictor/GCSUploader/src/cmd"
	logger "github.com/farovictor/GCSUploader/src/logging"
)

const projectID string = "Casa"
const bucket string = "casaBucket"
const emulator string = "localhost:9023"

func smain() {

	// Creating workers
	var wg sync.WaitGroup
	var pipe chan string = make(chan string)

	// Distributin workers
	// TODO: Implement function
	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup, pipe <-chan string) {
			defer wg.Done()
			for file := range pipe {
				println(file)
			}
		}(&wg, pipe)
		wg.Add(1)
	}

	// Distributing files to pipeline
	// TODO

	ctx := context.Background()

	// Set STORAGE_EMULATOR_HOST environment variable.
	err := os.Setenv("STORAGE_EMULATOR_HOST", emulator)
	if err != nil {
		// TODO: Handle error.
		logger.ErrorLogger.Fatalln(err)
	}

	// Create client as usual.
	client, err := storage.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
		logger.ErrorLogger.Fatalln(err)
	}

	// This request is now directed to http://localhost:9023/storage/v1/b
	// instead of https://storage.googleapis.com/storage/v1/b
	if err := client.Bucket(bucket).Create(ctx, projectID, nil); err != nil {
		// TODO: Handle error.
		logger.ErrorLogger.Fatalln(err)
	}
	logger.InfoLogger.Printf("Bucket %s created successfully\n", bucket)

	// Wait pipe to close and workers to finish
	wg.Wait()
}

func main() {
	cmd.Execute()
}
