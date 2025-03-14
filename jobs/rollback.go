package jobs

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"vnecro/config"
	"vnecro/vmOperations"
)

// RollbackVM attempts to roll back the given VM to the specified snapshot.
// It first attempts to shut down the VM (logging a warning if that fails) and then
// restores the snapshot. Returns an error if the rollback fails.
func RollbackVM(vmConfig *config.VMConfig, rollbackSnapshot string, operator vmOperations.VMOperator) error {
	logrus.Infof("Rolling back VM '%s' to snapshot '%s'", vmConfig.VMName, rollbackSnapshot)

	// Attempt to shut down the VM; if shutdown fails, log a warning and continue.
	if err := operator.Shutdown(vmConfig.VMName); err != nil {
		logrus.Warnf("Error shutting down VM '%s' during rollback: %v", vmConfig.VMName, err)
	}

	// Use the VirtualBox-specific rollback function from vmOperations.
	if err := operator.Rollback(vmConfig.VMName, rollbackSnapshot); err != nil {
		return fmt.Errorf("rollback failed: error restoring snapshot '%s' on VM '%s': %w", rollbackSnapshot, vmConfig.VMName, err)
	}

	logrus.Infof("Rollback successful: VM '%s' is now restored to snapshot '%s'", vmConfig.VMName, rollbackSnapshot)
	return nil
}
