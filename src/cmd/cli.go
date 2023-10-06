package cmd

import (
	"fmt"
	"os"

	logger "github.com/farovictor/GCSUploader/src/logging"
	"github.com/spf13/cobra"
)

var (
	Version   string
	GitCommit string
	BuildTime string

	blobPath           string
	blobPrefixPath     string
	blobPrefixName     string
	bucketName         string
	concurrency        int32
	concurrencyDefault int32 = 32
	logLevel           string
	projectID          string
	reuseName          bool
	searchPath         string
	sourceFile         string
	sourceDirectory    string
	sourcePrefix       string
)

// Root Command (does nothing, only prints nice things)
var rootCmd = &cobra.Command{
	Short:   "This project aims to support mongodb loading pipelines",
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("For more info, visit: https://github.com/farovictor/GCSLoader\n")
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
var loadPrefixCmd = &cobra.Command{
	Use:     "load-prefix",
	Version: rootCmd.Version,
	Short:   "Load files that match prefix",
	Run:     loadPrefix,
}

// Executes cli
func Execute() {
	// TODO: Development dependency. Remove on prod.
	setEmulator()
	if err := rootCmd.Execute(); err != nil {
		logger.ErrorLogger.Printf("%v %s\n", os.Stderr, err)
		println()
		os.Exit(1)
	}
}

func init() {

	// Root command flags setup
	rootCmd.PersistentFlags().StringVarP(&bucketName, "bucket-name", "b", "", "Bucket name")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Set a max log level")
	// Create
	createBucketCmd.PersistentFlags().StringVarP(&projectID, "project-id", "i", "", "Project ID")
	createBucketCmd.MarkPersistentFlagRequired("project-id")
	// Load file flags setup
	loadFileCmd.PersistentFlags().StringVarP(&blobPath, "blob-path", "n", "", "Blob path to be created")
	loadFileCmd.PersistentFlags().StringVarP(&sourceFile, "source-file", "f", "", "File to be uploaded")
	loadFileCmd.MarkFlagsRequiredTogether("blob-path", "source-file")
	// Load Prefix
	loadPrefixCmd.PersistentFlags().StringVarP(&searchPath, "search-path", "s", ".", "Search path")
	loadPrefixCmd.PersistentFlags().StringVarP(&blobPrefixPath, "blob-prefix-path", "p", "", "Blob prefix path (this is regarding path)")
	loadPrefixCmd.PersistentFlags().StringVarP(&blobPrefixName, "blob-prefix-name", "a", "", "Blob prefix name")
	loadPrefixCmd.PersistentFlags().BoolVarP(&reuseName, "reuse-name", "r", true, "Reuse name")
	loadPrefixCmd.PersistentFlags().Int32VarP(&concurrency, "num-concurrent-files", "n", concurrencyDefault, "Search path")
	// loadFileCmd.MarkFlagsRequiredTogether("blob-prefix-path", "blob-prefix-name")

	// Attaching commands to root
	rootCmd.AddCommand(createBucketCmd)
	rootCmd.AddCommand(deleteBucketCmd)
	rootCmd.AddCommand(attrsBucketCmd)
	rootCmd.AddCommand(existsBucketCmd)
	rootCmd.AddCommand(loadFileCmd)
	rootCmd.AddCommand(loadPrefixCmd)
}
