package cmd

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"

	"cloud.google.com/go/storage"
	logger "github.com/farovictor/GCSUploader/src/logging"
	utils "github.com/farovictor/GCSUploader/src/utils"
	"github.com/spf13/cobra"
)

const emulator string = "localhost:9023"

// This serves development process. It should not be called in prod
func setEmulator() {
	// Set STORAGE_EMULATOR_HOST environment variable.
	err := os.Setenv("STORAGE_EMULATOR_HOST", emulator)
	if err != nil {
		// TODO: Handle error.
		logger.ErrorLogger.Fatalln(err)
	}
	logger.InfoLogger.Println("Setting and Using Emulator", emulator)
}

// Create bucket
func createBucket(cmd *cobra.Command, args []string) {
	logger.Initialize(logLevel)

	if bucketName == "" {
		logger.ErrorLogger.Fatalln("Specify a bucket name")
	}
	ctx := context.Background()

	// Create client as usual.
	client, err := storage.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
		logger.ErrorLogger.Fatalln(err)
	}
	defer client.Close()

	// This request is now directed to http://localhost:9023/storage/v1/b
	// instead of https://storage.googleapis.com/storage/v1/b
	if err := client.Bucket(bucketName).Create(ctx, projectID, nil); err != nil {
		// TODO: Handle error.
		logger.ErrorLogger.Fatalln(err)
	}
	logger.InfoLogger.Printf("Bucket %s created successfully\n", bucketName)
}

// Delete Bucket
func deleteBucket(cmd *cobra.Command, args []string) {
	logger.Initialize(logLevel)

	if bucketName == "" {
		logger.ErrorLogger.Fatalln("Specify a bucket name")
	}
	ctx := context.Background()

	// Create client as usual.
	client, err := storage.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
		logger.ErrorLogger.Fatalln(err)
	}
	defer client.Close()

	if err = client.Bucket(bucketName).Delete(ctx); err != nil {
		logger.ErrorLogger.Fatalln(err)
	}

	logger.InfoLogger.Printf("Bucket %s deleted successfully\n", bucketName)
}

// Bucket Exists
func bucketExists(cmd *cobra.Command, args []string) {
	logger.Initialize(logLevel)

	if bucketName == "" {
		logger.ErrorLogger.Fatalln("Specify a bucket name")
	}
	ctx := context.Background()

	// Create client as usual.
	client, err := storage.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
		logger.ErrorLogger.Fatalln(err)
	}
	defer client.Close()

	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}
	if attrs == nil {
		logger.ErrorLogger.Fatalf("Bucket %s doesn't exists", bucketName)
	}

	logger.InfoLogger.Printf("Bucket %s exists", bucketName)
}

// Get Bucket Attribute
func attrsBucket(cmd *cobra.Command, args []string) {
	logger.Initialize(logLevel)

	if bucketName == "" {
		logger.ErrorLogger.Fatalln("Specify a bucket name")
	}
	ctx := context.Background()

	// Create client as usual.
	client, err := storage.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
		logger.ErrorLogger.Fatalln(err)
	}
	defer client.Close()

	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}

	jsonBytes, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}

	logger.InfoLogger.Printf("Attributes for bucket %s\n%s\n", bucketName, string(jsonBytes))
}

// Load File
func loadFile(cmd *cobra.Command, args []string) {
	logger.Initialize(logLevel)

	if bucketName == "" {
		logger.ErrorLogger.Fatalln("Specify a bucket name")
	}

	if sourceFile == "" {
		logger.ErrorLogger.Fatalln("Specify a source file")
	}

	file, err := os.Open(sourceFile)
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}
	defer file.Close()

	ctx := context.Background()

	// Create client as usual.
	client, err := storage.NewClient(ctx)
	defer client.Close()
	if err != nil {
		// TODO: Handle error.
		logger.ErrorLogger.Fatalln(err)
	}

	// Creating blob
	obj := client.Bucket(bucketName).Object(blobPath)
	writer := obj.NewWriter(ctx)

	if _, err := io.Copy(writer, file); err != nil {
		logger.ErrorLogger.Fatalln(err)
	}

	if err := writer.Close(); err != nil {
		logger.ErrorLogger.Fatalln(err)
	}

	logger.InfoLogger.Println("Blob created successfully:", blobPath)
}

// Load files that match prefix
func loadPrefix(cmd *cobra.Command, args []string) {
	logger.Initialize(logLevel)

	ctx := context.Background()

	// Check if bucket name is set
	if bucketName == "" {
		logger.ErrorLogger.Fatalln("Specify a bucket name")
	}

	// Check if search path is not empty
	if searchPath == "" {
		logger.ErrorLogger.Fatalln("Specify a search path")
	}

	// Initializing a Error Collector
	bc := utils.BatchCollector{}

	// Creating workers
	var wg sync.WaitGroup
	var pipe chan string = make(chan string)

	// Create client as usual.
	client, err := storage.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
		logger.ErrorLogger.Println(err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	// Distributin workers
	for i := 0; i < int(concurrency); i++ {
		wg.Add(1)
		go DispatchThis(&ctx, &wg, pipe, bucket, &bc)
	}

	// Read files and send to channel
	total, err := utils.EmitFilesToChannel(sourcePrefix, searchPath, pipe)
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}

	// Closing pipe
	close(pipe)

	logger.InfoLogger.Printf("Sent %d files to %d workers\n", total, concurrency)

	// Wait until channel is drainned and all workers are done
	wg.Wait()

	amountErrors := len(bc.FilesNotProcessed)
	perc := float32(amountErrors) / float32(total) * 100
	logger.InfoLogger.Printf("Done loading %d files (%.2f %%) - total %d\n", total-amountErrors, perc, total)
}

// dispatcher function
func DispatchThis(ctx *context.Context, wg *sync.WaitGroup, pipe <-chan string, bucket *storage.BucketHandle, coll *utils.BatchCollector) {
	defer wg.Done()

	for source := range pipe {

		// Open file
		filePath := filepath.Join(searchPath, source)
		file, err := os.Open(filePath)
		if err != nil {
			// logger.ErrorLogger.Println(err)
			coll.AddError(filePath)
			continue
		}

		// Read file content
		content, err := io.ReadAll(file)
		if err != nil {
			// logger.ErrorLogger.Println(err)
			coll.AddError(filePath)
			continue
		}

		if err := file.Close(); err != nil {
			// logger.ErrorLogger.Println(err)
			coll.AddError(filePath)
			continue
		}

		// Assemble blob name
		aBlobPath := utils.BlobNameAssemble(blobPrefixPath, blobPrefixName, source)

		// Creating blob
		obj := bucket.Object(aBlobPath)
		writer := obj.NewWriter(*ctx)

		// Writing file to writer
		_, err = writer.Write(content)
		if err != nil {
			// logger.ErrorLogger.Println(err)
			coll.AddError(filePath)
			continue
		}

		// Closing writer
		// READ: https://github.com/googleapis/google-cloud-go/issues/7090
		if err := writer.Close(); err != nil {
			// logger.ErrorLogger.Println(err)
			coll.AddError(filePath)
			continue
		}
	}
}
