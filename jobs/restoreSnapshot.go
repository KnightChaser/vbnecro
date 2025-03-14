package jobs

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"vnecro/config"
	"vnecro/vmOperations"
)

// RestoreSnapshot restores the given VM to a specified snapshot.
// It first lists the snapshots, then either uses the provided snapshot name or parses the first available one.
// Returns an error if any step fails.
func RestoreSnapshot(vmConfig *config.VMConfig, op config.Operation, operator vmOperations.VMOperator) error {
	logrus.Infof("Listing snapshots for VM '%s'", vmConfig.VMName)
	output, err := operator.ListSnapshots(vmConfig.VMName)
	if err != nil {
		return fmt.Errorf("error listing snapshots for VM '%s': %w", vmConfig.VMName, err)
	}
	logrus.Info("Snapshot list output:")
	logrus.Info(output)

	var snapshotToRestore string
	if val, ok := op.Params["snapshot"].(string); ok && val != "" {
		snapshotToRestore = val
	}
	if snapshotToRestore == "" {
		snapshotToRestore, err = operator.ParseSnapshot(output)
		if err != nil {
			return fmt.Errorf("error parsing snapshot for VM '%s': %w", vmConfig.VMName, err)
		}
	}
	logrus.Infof("Restoring VM '%s' to snapshot '%s'", vmConfig.VMName, snapshotToRestore)
	if err := operator.RestoreSnapshot(vmConfig.VMName, snapshotToRestore); err != nil {
		return fmt.Errorf("error restoring snapshot for VM '%s': %w", vmConfig.VMName, err)
	}
	logrus.Info("Snapshot restored successfully!")
	return nil
}
