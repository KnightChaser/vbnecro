package jobs

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"vnecro/config"
	"vnecro/vmOperations"
)

// ShutdownVM shuts down the VM specified in vmConfig using the provided operator.
// Returns an error if the shutdown fails.
func ShutdownVM(vmConfig *config.VMConfig, operator vmOperations.VMOperator) error {
	logrus.Infof("Shutting down VM '%s'", vmConfig.VMName)
	if err := operator.Shutdown(vmConfig.VMName); err != nil {
		return fmt.Errorf("error shutting down VM '%s': %w", vmConfig.VMName, err)
	}
	logrus.Info("VM shut down successfully!")
	return nil
}
