package main

import (
	"fmt"
	"os"

	"htpatcher/internal/repository"
	"htpatcher/internal/service"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Initialize repositories
	patchRepo := repository.NewPatchRepository()
	storageRepo := repository.NewStorageRepository()

	// Initialize logger (writes to stderr so it doesn't interfere with MCP stdio)
	logger := &MCPLogger{}

	// Initialize services
	collectionService, err := service.NewCollectionService(storageRepo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize collection service: %v\n", err)
		os.Exit(1)
	}

	svc := &Services{
		Game:       service.NewGameService(logger),
		Patch:      service.NewPatchService(patchRepo, logger),
		Backup:     service.NewBackupService(logger),
		Export:     service.NewExportService(logger),
		Collection: collectionService,
		Download:   service.NewDownloadService(patchRepo, logger),
	}

	s := server.NewMCPServer(
		"htpatcher",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	registerTools(s, svc)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

type Services struct {
	Game       *service.GameService
	Patch      *service.PatchService
	Backup     *service.BackupService
	Export     *service.ExportService
	Collection *service.CollectionService
	Download   *service.DownloadService
}

// MCPLogger implements service.Logger, writing to stderr
type MCPLogger struct{}

func (l *MCPLogger) Info(message string)    { fmt.Fprintf(os.Stderr, "[INFO] %s\n", message) }
func (l *MCPLogger) Success(message string) { fmt.Fprintf(os.Stderr, "[OK] %s\n", message) }
func (l *MCPLogger) Error(message string)   { fmt.Fprintf(os.Stderr, "[ERROR] %s\n", message) }
func (l *MCPLogger) Warn(message string)    { fmt.Fprintf(os.Stderr, "[WARN] %s\n", message) }
