package jobs

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"vnecro/config"
	"vnecro/vmOperations"
)

// PauseVM pauses the virtual machine specified in vmConfig using the provided operator.
// It returns an error if the operation fails.
func PauseVM(vmConfig *config.VMConfig, operator vmOperations.VMOperator) error {
	logrus.Infof("Pausing VM '%s'", vmConfig.VMName)
	if err := operator.Pause(vmConfig.VMName); err != nil {
		return fmt.Errorf("error pausing VM '%s': %w", vmConfig.VMName, err)
	}
	logrus.Info("VM paused successfully!")
	return nil
}
