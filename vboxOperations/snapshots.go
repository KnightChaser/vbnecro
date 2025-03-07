package vboxOperations

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ListSnapshots lists snapshots for a given VM and returns the raw output.
func ListSnapshots(vmName string) (string, error) {
	cmd := exec.Command("VBoxManage", "snapshot", vmName, "list", "--details")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error listing snapshots: %w", err)
	}
	return out.String(), nil
}

// ParseSnapshot extracts a clean snapshot name from the given raw snapshot listing.
// It expects lines like: "Name: InitialInstallation (UUID: ...)" and returns "InitialInstallation".
func ParseSnapshot(snapshotOutput string) (string, error) {
	lines := strings.Split(snapshotOutput, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Name:") {
			// Remove the "Name:" prefix and trim spaces.
			nameLine := strings.TrimSpace(strings.TrimPrefix(line, "Name:"))
			// Split on " (" to remove the UUID and any extra markers.
			if idx := strings.Index(nameLine, " ("); idx != -1 {
				return nameLine[:idx], nil
			}
			return nameLine, nil
		}
	}
	return "", fmt.Errorf("no snapshot found")
}

// RestoreSnapshot restores the given VM to the specified snapshot.
func RestoreSnapshot(vmName, snapshot string) error {
	cmd := exec.Command("VBoxManage", "snapshot", vmName, "restore", snapshot)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error restoring snapshot '%s': %w", snapshot, err)
	}
	return nil
}
