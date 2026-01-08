package main

import (
	"embed"
	"fmt"
	"htpatcher/internal/service"
	"htpatcher/internal/util"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Check for command line arguments before starting GUI
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--update":
			if len(os.Args) < 3 {
				fmt.Println("Usage: htpatcher --update <target_path>")
				os.Exit(1)
			}
			if err := performSelfUpdate(os.Args[2]); err != nil {
				fmt.Printf("Update failed: %v\n", err)
				os.Exit(1)
			}
			os.Exit(0)

		case "--version":
			fmt.Printf("htpatcher v%d\n", service.Version)
			os.Exit(0)
		}
	}

	// Check if just updated
	justUpdated := len(os.Args) > 1 && os.Args[1] == "--updated"

	// Create an instance of the app structure
	app := NewApp()
	app.justUpdated = justUpdated

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "HTTranslations Patcher",
		Width:  1280,
		Height: 720,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// performSelfUpdate handles the --update flag logic
// It waits for the original process to exit, copies itself to the target path,
// then launches the updated executable and deletes itself from cache
func performSelfUpdate(targetPath string) error {
	currentExePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	// Wait for original process to exit by polling file lock
	// Try to open target file for writing to check if it's released
	maxRetries := 30 // 30 seconds max
	for i := 0; i < maxRetries; i++ {
		time.Sleep(1 * time.Second)

		// Try to open target for writing
		f, err := os.OpenFile(targetPath, os.O_WRONLY, 0)
		if err == nil {
			f.Close()
			break // File is no longer locked
		}

		if i == maxRetries-1 {
			return fmt.Errorf("timeout waiting for process to exit")
		}
	}

	// Copy ourselves to target location (overwrite)
	if err := copyFile(currentExePath, targetPath); err != nil {
		return fmt.Errorf("failed to replace executable: %w", err)
	}

	// Launch the updated version with --updated flag
	cmd := exec.Command(targetPath, "--updated")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch updated app: %w", err)
	}

	// Brief delay to ensure process starts
	time.Sleep(500 * time.Millisecond)

	// Delete ourselves from cache (best effort)
	os.Remove(currentExePath)

	// Clean up the entire update cache
	util.CleanUpdateCache()

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err = io.Copy(destination, source); err != nil {
		return err
	}

	return destination.Sync()
}
