package vboxOperations

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// StartVM starts a VirtualBox VM in headless mode.
func StartVM(vmName string) error {
	cmd := exec.Command("VBoxManage", "startvm", vmName, "--type", "headless")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error starting VM '%s': %w", vmName, err)
	}
	return nil
}

// PauseVM pauses a running VirtualBox VM.
func PauseVM(vmName string) error {
	cmd := exec.Command("VBoxManage", "controlvm", vmName, "pause")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error pausing VM '%s': %w", vmName, err)
	}
	return nil
}

// ResumeVM resumes a paused VirtualBox VM.
func ResumeVM(vmName string) error {
	cmd := exec.Command("VBoxManage", "controlvm", vmName, "resume")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error resuming VM '%s': %w", vmName, err)
	}
	return nil
}

// ShutdownVM attempts to shut down a VM.
// If the VM is not running or is aborted, it treats that as success.
func ShutdownVM(vmName string) error {
	err := shutdown(vmName)
	if err != nil {
		errMsg := strings.ToLower(err.Error())

		// Treat "not currently running" or "aborted" as success.
		if strings.Contains(errMsg, "not currently running") || strings.Contains(errMsg, "aborted") {
			return nil
		}

		// If the error indicates that the VM is paused, try to resume and retry shutdown.
		if strings.Contains(errMsg, "paused") {
			if resumeErr := ResumeVM(vmName); resumeErr != nil {
				return fmt.Errorf("failed to resume paused VM '%s': %w", vmName, resumeErr)
			}
			// Retry shutdown after resuming.
			err = shutdown(vmName)
			errMsg = strings.ToLower(err.Error())
			if err != nil && (!strings.Contains(errMsg, "not currently running") && !strings.Contains(errMsg, "aborted")) {
				return fmt.Errorf("error shutting down VM '%s' after resuming: %w", vmName, err)
			}
			return nil
		}
		return fmt.Errorf("error shutting down VM '%s': %w", vmName, err)
	}
	return nil
}

// shutdown is a helper that issues the poweroff command and captures error output.
func shutdown(vmName string) error {
	var out bytes.Buffer
	cmd := exec.Command("VBoxManage", "controlvm", vmName, "poweroff")
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v: %s", err, out.String())
	}
	return nil
}
