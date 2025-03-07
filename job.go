package main

import (
	"fmt"
	"log"
)

func ProcessJobs(cfg *Config) {
	for _, job := range cfg.Jobs {
		vmConfig, err := GetVMConfig(cfg.VMs, job.VMAlias)
		if err != nil {
			log.Printf("Skipping job: %v", err)
			continue
		}

		// For each operation defined for this job.
		for _, op := range job.Operations {
			switch op.Type {
			case "RestoreSnapshot":
				log.Printf("Listing snapshots for VM '%s'", vmConfig.VMName)
				output, err := ListSnapshots(vmConfig.VMName)
				if err != nil {
					log.Printf("Error listing snapshots: %v", err)
					continue
				}
				fmt.Println("Snapshot list output:")
				fmt.Println(output)

				// Determine which snapshot to restore from the Params map.
				var snapshotToRestore string
				if val, ok := op.Params["snapshot"]; ok {
					if s, ok := val.(string); ok && s != "" {
						snapshotToRestore = s
					}
				}

				// If not provided in Params, parse the first available snapshot.
				if snapshotToRestore == "" {
					snapshotToRestore, err = ParseSnapshot(output)
					if err != nil {
						log.Printf("Error parsing snapshot: %v", err)
						continue
					}
				}

				log.Printf("Restoring VM '%s' to snapshot '%s'", vmConfig.VMName, snapshotToRestore)
				if err := RestoreSnapshot(vmConfig.VMName, snapshotToRestore); err != nil {
					log.Printf("Error restoring snapshot: %v", err)
					continue
				}
				log.Println("Snapshot restored successfully!")
			default:
				log.Printf("Unknown operation type: %s", op.Type)
			}
		}
	}
}
