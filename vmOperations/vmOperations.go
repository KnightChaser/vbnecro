package vmOperations

import (
	"time"

	"vnecro/vboxOperations"
)

// VMOperator defines the interface for performing operations on virtual machines.
// This abstraction allows for different backends (e.g., VirtualBox, Hyper-V, etc.).
type VMOperator interface {
	// Start launches the virtual machine identified by vmName.
	Start(vmName string) error

	// Pause suspends the virtual machine identified by vmName.
	Pause(vmName string) error

	// Shutdown turns off the virtual machine identified by vmName.
	Shutdown(vmName string) error

	// RestoreSnapshot reverts the virtual machine to the specified snapshot.
	RestoreSnapshot(vmName, snapshot string) error

	// ListSnapshots returns a string with the list of snapshots for the virtual machine.
	ListSnapshots(vmName string) (string, error)

	// ParseSnapshot extracts a clean snapshot name from the given snapshot output.
	ParseSnapshot(snapshotOutput string) (string, error)

	// WaitForGuestExecReady waits until the guest execution service is ready,
	// given the VM name, credentials, and a timeout duration.
	WaitForGuestExecReady(vmName, username, password string, timeout time.Duration) error

	// ExecuteShellCommand executes a command inside the guest OS with the provided arguments.
	ExecuteShellCommand(vmName, username, password, command string, args ...string) (string, error)
}

// VirtualBoxOperator is a concrete implementation of VMOperator using VirtualBox's VBoxManage tool.
type VirtualBoxOperator struct{}

// NewVirtualBoxOperator returns a new instance of VirtualBoxOperator.
func NewVirtualBoxOperator() VMOperator {
	return &VirtualBoxOperator{}
}

// Start launches the virtual machine using VBoxManage.
func (v *VirtualBoxOperator) Start(vmName string) error {
	return vboxOperations.StartVM(vmName)
}

// Pause suspends the virtual machine using VBoxManage.
func (v *VirtualBoxOperator) Pause(vmName string) error {
	return vboxOperations.PauseVM(vmName)
}

// Shutdown turns off the virtual machine using VBoxManage.
// It handles cases where the VM is already off or aborted.
func (v *VirtualBoxOperator) Shutdown(vmName string) error {
	return vboxOperations.ShutdownVM(vmName)
}

// RestoreSnapshot reverts the virtual machine to a specified snapshot.
func (v *VirtualBoxOperator) RestoreSnapshot(vmName, snapshot string) error {
	return vboxOperations.RestoreSnapshot(vmName, snapshot)
}

// ListSnapshots returns the snapshot list of a virtual machine.
func (v *VirtualBoxOperator) ListSnapshots(vmName string) (string, error) {
	return vboxOperations.ListSnapshots(vmName)
}

// ParseSnapshot extracts a clean snapshot name from the snapshot list output.
func (v *VirtualBoxOperator) ParseSnapshot(snapshotOutput string) (string, error) {
	return vboxOperations.ParseSnapshot(snapshotOutput)
}

// WaitForGuestExecReady polls until the guest execution service is ready, or the timeout expires.
func (v *VirtualBoxOperator) WaitForGuestExecReady(vmName, username, password string, timeout time.Duration) error {
	return vboxOperations.WaitForGuestExecReady(vmName, username, password, timeout)
}

// ExecuteShellCommand runs a shell command inside the guest OS and returns its output.
func (v *VirtualBoxOperator) ExecuteShellCommand(vmName, username, password, command string, args ...string) (string, error) {
	return vboxOperations.ExecuteShellCommand(vmName, username, password, command, args...)
}
