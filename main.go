package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of our configuration file.
type Config struct {
	VMs  []VMConfig  `yaml:"vms"`
	Jobs []JobConfig `yaml:"jobs"`
}

// VMConfig holds the VirtualBox VM configuration.
type VMConfig struct {
	Alias    string `yaml:"alias"`
	VMName   string `yaml:"vm_name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// JobConfig represents a job to perform on a VM.
type JobConfig struct {
	VMAlias    string      `yaml:"vm_alias"`
	Operations []Operation `yaml:"operations"`
}

// Operation holds the type of operation and its parameters.
type Operation struct {
	Type     string `yaml:"type"`
	Snapshot string `yaml:"snapshot"`
}

// loadConfig reads and parses the YAML configuration.
func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// listSnapshots lists snapshots for a given VM and returns the raw output.
func listSnapshots(vmName string) (string, error) {
	cmd := exec.Command("VBoxManage", "snapshot", vmName, "list", "--details")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error listing snapshots: %w", err)
	}
	return out.String(), nil
}

// parseSnapshot extracts a clean snapshot name from the given raw snapshot listing.
// It expects lines like: "Name: InitialInstallation (UUID: ...)" and returns "InitialInstallation".
func parseSnapshot(snapshotOutput string) (string, error) {
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

// restoreSnapshot restores the given VM to the specified snapshot.
func restoreSnapshot(vmName, snapshot string) error {
	cmd := exec.Command("VBoxManage", "snapshot", vmName, "restore", snapshot)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error restoring snapshot '%s': %w", snapshot, err)
	}
	return nil
}

// getVMConfig finds a VM configuration by its alias.
func getVMConfig(vms []VMConfig, alias string) (*VMConfig, error) {
	for _, vm := range vms {
		if vm.Alias == alias {
			return &vm, nil
		}
	}
	return nil, fmt.Errorf("VM with alias '%s' not found", alias)
}

func main() {
	// Define command-line flag for configuration file path.
	configPath := flag.String("config-path", "", "Path to the YAML configuration file")
	flag.Parse()
	if *configPath == "" {
		log.Fatal("Missing required flag: --config-path")
	}

	// Load configuration from YAML.
	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config from '%s': %v", *configPath, err)
	}

	// Process each job from the config.
	for _, job := range cfg.Jobs {
		vmConfig, err := getVMConfig(cfg.VMs, job.VMAlias)
		if err != nil {
			log.Printf("Skipping job: %v", err)
			continue
		}

		// For each operation defined for this job.
		for _, op := range job.Operations {
			switch op.Type {
			case "RestoreSnapshot":
				// List snapshots to verify the snapshot exists.
				fmt.Printf("Listing snapshots for VM '%s'\n", vmConfig.VMName)
				output, err := listSnapshots(vmConfig.VMName)
				if err != nil {
					log.Printf("Error listing snapshots: %v", err)
					continue
				}
				fmt.Println("Snapshot list output:")
				fmt.Println(output)

				// Use the snapshot from the operation or parse from the list.
				var snapshotToRestore string
				if op.Snapshot != "" {
					snapshotToRestore = op.Snapshot
				} else {
					snapshotToRestore, err = parseSnapshot(output)
					if err != nil {
						log.Printf("Error parsing snapshot: %v", err)
						continue
					}
				}

				fmt.Printf("Restoring VM '%s' to snapshot '%s'\n", vmConfig.VMName, snapshotToRestore)
				if err := restoreSnapshot(vmConfig.VMName, snapshotToRestore); err != nil {
					log.Printf("Error restoring snapshot: %v", err)
					continue
				}
				fmt.Println("Snapshot restored successfully!")
			default:
				log.Printf("Unknown operation type: %s", op.Type)
			}
		}
	}
}
