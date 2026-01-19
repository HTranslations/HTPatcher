package util

import (
	"errors"
	"os/exec"
	goruntime "runtime"
)

// OpenFolder opens a folder in the system's file explorer
func OpenFolder(folderPath string) error {
	var cmd *exec.Cmd
	switch goruntime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", folderPath)
	case "darwin":
		cmd = exec.Command("open", folderPath)
	case "linux":
		cmd = exec.Command("xdg-open", folderPath)
	default:
		return errors.New("unsupported operating system")
	}
	return cmd.Start()
}

// LaunchExecutable launches an executable file
func LaunchExecutable(exePath, workingDir string) error {
	cmd := exec.Command(exePath)
	cmd.Dir = workingDir
	return cmd.Start()
}




