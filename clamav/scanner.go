package clamav

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/tedmo/cav/logger"
	"go.uber.org/zap"
	"io"
	"os"
	"os/exec"
)

type FileScanner struct {
	TempDirectory string
}

func (s *FileScanner) Scan(ctx context.Context, contents io.Reader) ([]byte, error) {
	log := logger.New(ctx)
	// Create temporary file with unique name
	filepath := fmt.Sprintf("%s/%s", s.TempDirectory, uuid.NewString())
	tempFile, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()
	defer os.Remove(filepath)

	log.Debug("writing to file...")
	_, err = io.Copy(tempFile, contents)
	if err != nil {
		return nil, err
	}
	log.Debug("wrote to file")

	return s.scanFilepath(ctx, filepath)
}

func (s *FileScanner) scanFilepath(ctx context.Context, file string) ([]byte, error) {
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
	log := logger.New(ctx).With(zap.String("command", cmd.String()))

	log.Info("executing command")
	output, err := cmd.CombinedOutput()
	log.Infow("command completed", zap.String("output", string(output)))
	return output, err
}
