package vboxOperations

import (
	"fmt"
	"os/exec"
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
// Note: VirtualBox uses the "controlvm" command for power management.
func PauseVM(vmName string) error {
	cmd := exec.Command("VBoxManage", "controlvm", vmName, "pause")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error pausing VM '%s': %w", vmName, err)
	}
	return nil
}
