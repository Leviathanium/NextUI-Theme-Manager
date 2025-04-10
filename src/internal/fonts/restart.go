// src/internal/fonts/restart.go
// Functions for restarting the application

package fonts

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"nextui-themes/internal/logging"
)

// RestartApp restarts the application to apply font changes
func RestartApp() error {
	logging.LogDebug("Restarting application to apply font changes")

	// Get current directory and executable path
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return err
	}

	// Get the path to the application executable
	exePath := filepath.Join(cwd, "theme-manager")

	// Create a new process
	cmd := exec.Command(exePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = cwd

	// Set the new process to replace this one
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Start the new process
	err = cmd.Start()
	if err != nil {
		logging.LogDebug("Error starting new process: %v", err)
		return err
	}

	logging.LogDebug("Successfully started new process, exiting current process")

	// Give the new process a moment to start
	time.Sleep(500 * time.Millisecond)

	// Exit this process
	os.Exit(0)

	return nil // This line is never reached
}