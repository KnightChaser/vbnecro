package main

import (
	"fmt"
	"log"

	"vbnecro/vboxOperations"
)

func ProcessJobs(cfg *Config) {
	for _, job := range cfg.Jobs {
		vmConfig, err := GetVMConfig(cfg.VMs, job.VMAlias)
		if err != nil {
			log.Fatalf("Job for VM alias '%s' failed: %v", job.VMAlias, err)
		}

		// Process each operation defined for this job.
		for _, op := range job.Operations {
			switch op.Type {
			case "RestoreSnapshot":
				log.Printf("Listing snapshots for VM '%s'", vmConfig.VMName)
				output, err := vboxOperations.ListSnapshots(vmConfig.VMName)
				if err != nil {
					log.Fatalf("Job failed: error listing snapshots for VM '%s': %v", vmConfig.VMName, err)
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
					snapshotToRestore, err = vboxOperations.ParseSnapshot(output)
					if err != nil {
						log.Fatalf("Job failed: error parsing snapshot for VM '%s': %v", vmConfig.VMName, err)
					}
				}

				log.Printf("Restoring VM '%s' to snapshot '%s'", vmConfig.VMName, snapshotToRestore)
				if err := vboxOperations.RestoreSnapshot(vmConfig.VMName, snapshotToRestore); err != nil {
					log.Fatalf("Job failed: error restoring snapshot for VM '%s': %v", vmConfig.VMName, err)
				}
				log.Println("Snapshot restored successfully!")

			case "StartVM":
				log.Printf("Starting VM '%s'", vmConfig.VMName)
				if err := vboxOperations.StartVM(vmConfig.VMName); err != nil {
					log.Fatalf("Job failed: error starting VM '%s': %v", vmConfig.VMName, err)
				}
				log.Println("VM started successfully!")

			case "PauseVM":
				log.Printf("Pausing VM '%s'", vmConfig.VMName)
				if err := vboxOperations.PauseVM(vmConfig.VMName); err != nil {
					log.Fatalf("Job failed: error pausing VM '%s': %v", vmConfig.VMName, err)
				}
				log.Println("VM paused successfully!")

			default:
				log.Fatalf("Job failed: unknown operation type: %s", op.Type)
			}
		}
	}
}
