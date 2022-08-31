package main

import (
	"fmt"
	"github.com/tedmo/cav/clamav"
	"github.com/tedmo/cav/http"
	"github.com/tedmo/cav/logger"
	"os"
)

const tempDir = "/temp/scan"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	log := logger.Init()
	defer log.Sync()

	server := &http.Server{
		Scanner: &clamav.FileScanner{TempDirectory: tempDir},
	}

	return server.ListenAndServe(":8080")
}
