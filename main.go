package main

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"os/exec"
)

const fileDir = "/temp/scan"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

var logger *zap.Logger

func NewLogger(ctx context.Context) *zap.SugaredLogger {
	return logger.Sugar()
}

func run() error {
	logger = zap.Must(zap.NewDevelopment())
	defer logger.Sync()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := handleFile(w, r); err != nil {
			NewLogger(r.Context()).Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return http.ListenAndServe(":8080", nil)
}

func handleFile(w http.ResponseWriter, r *http.Request) error {
	log := NewLogger(r.Context())
	if r.Method != http.MethodPost {
		return errors.New("unsupported method")
	}
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return err
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return err
	}

	filename := header.Filename
	if filename == "" {
		return errors.New("filename must not be empty")
	}

	fullFilePath := fmt.Sprintf("%s/%s", fileDir, filename)

	log.Debugw("creating temporary file", zap.String("file", fullFilePath))
	tempFile, err := os.Create(fullFilePath)
	if err != nil {
		return err
	}
	defer tempFile.Close()
	defer os.Remove(fullFilePath)

	log.Debug("writing to file...")
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return err
	}
	log.Debug("wrote to file")

	status := http.StatusOK
	output, err := scan(r.Context(), fullFilePath)
	if err != nil {
		status = http.StatusInternalServerError
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write(output)

	return nil
}

func scan(ctx context.Context, file string) ([]byte, error) {
	// freshclam
	cmd := exec.CommandContext(ctx, "freshclam")
	_, err := executeCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// clamscan -V
	cmd = exec.CommandContext(ctx, "clamscan", "-V")
	_, err = executeCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// clamscan {file}
	cmd = exec.CommandContext(ctx, "clamscan", file)
	return executeCommand(ctx, cmd)
}

func executeCommand(ctx context.Context, cmd *exec.Cmd) ([]byte, error) {
	log := NewLogger(ctx).With(zap.String("command", cmd.String()))

	log.Info("executing command")
	output, err := cmd.CombinedOutput()
	log.Infow("command completed", zap.String("output", string(output)))
	return output, err
}
