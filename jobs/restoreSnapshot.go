package jobs

import (
	"log"

	"vbnecro/config"
	"vbnecro/vboxOperations"
)

func RestoreSnapshot(vmConfig *config.VMConfig, op config.Operation) {
	log.Printf("Listing snapshots for VM '%s'", vmConfig.VMName)
	output, err := vboxOperations.ListSnapshots(vmConfig.VMName)
	if err != nil {
		log.Fatalf("Job failed: error listing snapshots for VM '%s': %v", vmConfig.VMName, err)
	}
	log.Println("Snapshot list output:")
	log.Println(output)

	var snapshotToRestore string
	if val, ok := op.Params["snapshot"].(string); ok && val != "" {
		snapshotToRestore = val
	}
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
}
