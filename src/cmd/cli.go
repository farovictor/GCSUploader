package cmd

import (
	"fmt"
	"os"

	logger "github.com/farovictor/GCSUploader/src/logging"
	utils "github.com/farovictor/GCSUploader/src/utils"
	"github.com/spf13/cobra"
)

var (
	Version   string
	GitCommit string
	BuildTime string

	authType             utils.AuthType = "file"
	authHolder           string
	blobPath             string
	blobPrefixPath       string
	blobPrefixName       string
	bucketName           string
	concurrency          int32
	concurrencyDefault   int32 = 32
	errorFileTracker     bool
	errorFileTrackerPath string
	logLevel             string
	logLevelDefault      string = "info"
	projectID            string
	reuseName            bool
	searchPath           string
	sourceFile           string
	sourceDirectory      string
	sourcePrefix         string
)

// Root Command (does nothing, only prints nice things)
var rootCmd = &cobra.Command{
	Short:   "This project aims to support mongodb loading pipelines",
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("For more info, visit: https://github.com/farovictor/GCSUploader\n")
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Built: %s\n", BuildTime)
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Log-Level: %v\n", logLevel)
	},
}

// Create Bucket Command
var createBucketCmd = &cobra.Command{
	Use:     "create-bucket",
	Version: rootCmd.Version,
	Short:   "Create bucket",
	Run:     createBucket,
}

// Delete Bucket Command
var deleteBucketCmd = &cobra.Command{
	Use:     "delete-bucket",
	Version: rootCmd.Version,
	Short:   "Delete bucket",
	Run:     deleteBucket,
}

// Attrs Bucket Command
var attrsBucketCmd = &cobra.Command{
	Use:     "attrs-bucket",
	Version: rootCmd.Version,
	Short:   "Attributes from bucket",
	Run:     attrsBucket,
}

// Check if bucket Exists Command
var existsBucketCmd = &cobra.Command{
	Use:     "exists-bucket",
	Version: rootCmd.Version,
	Short:   "Check if bucket exists",
	Run:     bucketExists,
}

// Check if bucket Exists Command
var loadFileCmd = &cobra.Command{
	Use:     "load",
	Version: rootCmd.Version,
	Short:   "Load file to bucket",
	Run:     loadFile,
}

// Check if bucket Exists Command
var loadBatchesCmd = &cobra.Command{
	Use:     "load-batch",
	Version: rootCmd.Version,
	Short:   "Load files concurrently using batches",
	Run:     loadBatches,
}

// Executes cli
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.ErrorLogger.Printf("%v %s\n", os.Stderr, err)
		println()
		os.Exit(1)
	}
}

func init() {

	// Root command flags setup
	rootCmd.PersistentFlags().StringVarP(&bucketName, "bucket-name", "b", "", "Bucket name")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", logLevelDefault, "Set a max log level")
	rootCmd.PersistentFlags().StringVarP((*string)(&authType), "auth-type", "a", utils.File, "Set what authentication type to use")
	rootCmd.PersistentFlags().StringVarP(&authHolder, "auth", "c", "", "Pass authentication")
	// Create
	createBucketCmd.PersistentFlags().StringVarP(&projectID, "project-id", "i", "", "Project ID")
	createBucketCmd.MarkPersistentFlagRequired("project-id")
	// Load file flags setup
	loadFileCmd.PersistentFlags().StringVarP(&blobPath, "blob-path", "n", "", "Blob path to be created")
	loadFileCmd.PersistentFlags().StringVarP(&sourceFile, "source-file", "f", "", "File to be uploaded")
	loadFileCmd.MarkFlagsRequiredTogether("blob-path", "source-file")
	// Load Prefix
	loadBatchesCmd.PersistentFlags().StringVarP(&searchPath, "search-path", "s", ".", "Search path")
	loadBatchesCmd.PersistentFlags().StringVarP(&blobPrefixPath, "blob-prefix-path", "p", "", "Blob prefix path (this is regarding path)")
	loadBatchesCmd.PersistentFlags().StringVarP(&blobPrefixName, "blob-prefix-name", "n", "", "Blob prefix name")
	loadBatchesCmd.PersistentFlags().BoolVarP(&reuseName, "reuse-name", "r", true, "Reuse name")
	loadBatchesCmd.PersistentFlags().Int32VarP(&concurrency, "num-concurrency", "x", concurrencyDefault, "Number of concurrent workers")
	loadBatchesCmd.PersistentFlags().BoolVarP(&errorFileTracker, "track-errors", "e", true, "Collect files that failed to upload")
	loadBatchesCmd.PersistentFlags().StringVarP(&errorFileTrackerPath, "error-track-logs", "t", ".", "Path to dump a list with files")

	// Attaching commands to root
	rootCmd.AddCommand(createBucketCmd)
	rootCmd.AddCommand(deleteBucketCmd)
	rootCmd.AddCommand(attrsBucketCmd)
	rootCmd.AddCommand(existsBucketCmd)
	rootCmd.AddCommand(loadFileCmd)
	rootCmd.AddCommand(loadBatchesCmd)
}
