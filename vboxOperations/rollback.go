package vboxOperations

import (
	"fmt"
)

// Rollback restores the given VM to the specified snapshot in case of contingencies.
// It first attempts to shut down the VM, then restores the snapshot.
// If the VM is already off (or aborted), it proceeds directly to the snapshot restoration.
func Rollback(vmName, snapshot string) error {
	// Attempt to shut down the VM.
	// ShutdownVM is already implemented to handle cases where the VM is not running.
	if err := ShutdownVM(vmName); err != nil {
		return fmt.Errorf("failed to shutdown VM '%s': %v", vmName, err)
	}
	// Restore the snapshot.
	if err := RestoreSnapshot(vmName, snapshot); err != nil {
		return fmt.Errorf("failed to restore snapshot '%s' on VM '%s': %v", snapshot, vmName, err)
	}
	return nil
}
