package jobs

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"vnecro/config"
	"vnecro/vmOperations"
)

// StartVM starts the VM specified in vmConfig using the provided operator.
// Returns an error if starting the VM fails.
func StartVM(vmConfig *config.VMConfig, operator vmOperations.VMOperator) error {
	logrus.Infof("Starting VM '%s'", vmConfig.VMName)
	if err := operator.Start(vmConfig.VMName); err != nil {
		return fmt.Errorf("error starting VM '%s': %w", vmConfig.VMName, err)
	}
	logrus.Info("VM started successfully!")
	return nil
}
