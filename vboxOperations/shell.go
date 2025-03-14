package vboxOperations

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// WaitForGuestExecReady polls the guest execution service by trying to run a simple echo command.
// It will keep retrying until the command succeeds or the timeout is reached, printing a logrus message each second.
func WaitForGuestExecReady(vmName, username, password string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	// Ensure the dummy command uses an absolute path.
	exe := "echo"
	if !strings.HasPrefix(exe, "/") {
		exe = "/bin/" + exe
	}

	currentSecond := 0
	for {
		// Build the command to check guest readiness.
		cmdArgs := []string{
			"guestcontrol", vmName, "run",
			"--username", username,
			"--password", password,
			"--exe", exe,
			"--", "ready",
		}
		cmd := exec.Command("VBoxManage", cmdArgs...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err == nil {
			// Command succeeded; guest execution service is ready.
			return nil
		}
		// If we've passed the deadline, return a timeout error.
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for guest execution service to be ready: last error: %v, output: %s", err, out.String())
		}
		// Print a waiting message every second.
		logrus.Printf("Waiting for guest execution service to be ready on VM '%s' (%d / %d seconds)", vmName, currentSecond, int(timeout.Seconds()))
		currentSecond++
		time.Sleep(1 * time.Second)
	}
}

// ExecuteShellCommand executes a shell command inside the guest OS.
// It uses VBoxManage guestcontrol run and requires Guest Additions to be installed.
func ExecuteShellCommand(vmName, username, password, command string, args ...string) (string, error) {
	// If the command does not start with "/", assume it's in /bin/ and prepend it.
	if len(command) > 0 && command[0] != '/' {
		command = "/bin/" + command
	}

	cmdArgs := []string{
		"guestcontrol", vmName, "run",
		"--username", username,
		"--password", password,
		"--exe", command,
		"--",
	}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("VBoxManage", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing shell command: %v, output: %s", err, output)
	}
	return string(output), nil
}
