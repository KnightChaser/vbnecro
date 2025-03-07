package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	// Replace these with your actual VM details.
	vmName := "vbnecro_ubuntu2204"
	// username := "vbnecro"  // Currently not used but available for further authentication logic if needed.
	// password := "pass12##" // Currently not used but available for further authentication logic if needed.

	// Execute VBoxManage command to list snapshots with details.
	cmd := exec.Command("VBoxManage", "snapshot", vmName, "list", "--details")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error listing snapshots: %v", err)
	}

	output := out.String()
	fmt.Println("Snapshot list output:")
	fmt.Println(output)

	// Parse the output to get the first snapshot name.
	// Expected format: "Name: <snapshot_name> (UUID: <uuid>) [*]"
	var firstSnapshot string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Name:") {
			// Remove the "Name:" prefix
			nameLine := strings.TrimPrefix(line, "Name:")
			nameLine = strings.TrimSpace(nameLine)
			// Split on " (" to separate the snapshot name from the rest of the details.
			if idx := strings.Index(nameLine, " ("); idx != -1 {
				firstSnapshot = nameLine[:idx]
			} else {
				firstSnapshot = nameLine
			}
			break
		}
	}

	if firstSnapshot == "" {
		log.Fatalf("No snapshot found for VM '%s'", vmName)
	}

	fmt.Printf("Restoring to the first snapshot: %s\n", firstSnapshot)

	// Restore the VM to the first snapshot found.
	restoreCmd := exec.Command("VBoxManage", "snapshot", vmName, "restore", firstSnapshot)
	if err := restoreCmd.Run(); err != nil {
		log.Fatalf("Error restoring snapshot '%s': %v", firstSnapshot, err)
	}

	fmt.Println("Snapshot restored successfully!")
}
