package vmOperations

import (
	"time"

	"vbnecro/vboxOperations"
)

// VMOperator defines the interface for VM operations.
type VMOperator interface {
	Start(vmName string) error
	Pause(vmName string) error
	Shutdown(vmName string) error
	RestoreSnapshot(vmName, snapshot string) error
	ListSnapshots(vmName string) (string, error)
	ParseSnapshot(snapshotOutput string) (string, error)
	WaitForGuestExecReady(vmName, username, password string, timeout time.Duration) error
	ExecuteShellCommand(vmName, username, password, command string, args ...string) (string, error)
}

// VirtualBoxOperator is a concrete implementation of VMOperator using VBoxManage.
type VirtualBoxOperator struct{}

// NewVirtualBoxOperator returns a new instance of VirtualBoxOperator.
func NewVirtualBoxOperator() VMOperator {
	return &VirtualBoxOperator{}
}

func (v *VirtualBoxOperator) Start(vmName string) error {
	return vboxOperations.StartVM(vmName)
}

func (v *VirtualBoxOperator) Pause(vmName string) error {
	return vboxOperations.PauseVM(vmName)
}

func (v *VirtualBoxOperator) Shutdown(vmName string) error {
	return vboxOperations.ShutdownVM(vmName)
}

func (v *VirtualBoxOperator) RestoreSnapshot(vmName, snapshot string) error {
	return vboxOperations.RestoreSnapshot(vmName, snapshot)
}

func (v *VirtualBoxOperator) ListSnapshots(vmName string) (string, error) {
	return vboxOperations.ListSnapshots(vmName)
}

func (v *VirtualBoxOperator) ParseSnapshot(snapshotOutput string) (string, error) {
	return vboxOperations.ParseSnapshot(snapshotOutput)
}

func (v *VirtualBoxOperator) WaitForGuestExecReady(vmName, username, password string, timeout time.Duration) error {
	return vboxOperations.WaitForGuestExecReady(vmName, username, password, timeout)
}

func (v *VirtualBoxOperator) ExecuteShellCommand(vmName, username, password, command string, args ...string) (string, error) {
	return vboxOperations.ExecuteShellCommand(vmName, username, password, command, args...)
}
