package jobs

import (
	"github.com/sirupsen/logrus"
	"vnecro/config"
	"vnecro/vmOperations"
)

func RestoreSnapshot(vmConfig *config.VMConfig, op config.Operation, operator vmOperations.VMOperator) {
	logrus.Printf("Listing snapshots for VM '%s'", vmConfig.VMName)
	output, err := operator.ListSnapshots(vmConfig.VMName)
	if err != nil {
		logrus.Fatalf("Job failed: error listing snapshots for VM '%s': %v", vmConfig.VMName, err)
	}
	logrus.Println("Snapshot list output:")
	logrus.Println(output)

	var snapshotToRestore string
	if val, ok := op.Params["snapshot"].(string); ok && val != "" {
		snapshotToRestore = val
	}
	if snapshotToRestore == "" {
		snapshotToRestore, err = operator.ParseSnapshot(output)
		if err != nil {
			logrus.Fatalf("Job failed: error parsing snapshot for VM '%s': %v", vmConfig.VMName, err)
		}
	}
	logrus.Printf("Restoring VM '%s' to snapshot '%s'", vmConfig.VMName, snapshotToRestore)
	if err := operator.RestoreSnapshot(vmConfig.VMName, snapshotToRestore); err != nil {
		logrus.Fatalf("Job failed: error restoring snapshot for VM '%s': %v", vmConfig.VMName, err)
	}
	logrus.Println("Snapshot restored successfully!")
}
