package jobs

import (
	"github.com/sirupsen/logrus"

	"vnecro/config"
	"vnecro/vmOperations"
)

func PauseVM(vmConfig *config.VMConfig, operator vmOperations.VMOperator) {
	logrus.Printf("Pausing VM '%s'", vmConfig.VMName)
	if err := operator.Pause(vmConfig.VMName); err != nil {
		logrus.Fatalf("Job failed: error pausing VM '%s': %v", vmConfig.VMName, err)
	}
	logrus.Println("VM paused successfully!")
}
