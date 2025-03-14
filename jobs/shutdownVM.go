package jobs

import (
	"github.com/sirupsen/logrus"

	"vnecro/config"
	"vnecro/vmOperations"
)

func ShutdownVM(vmConfig *config.VMConfig, operator vmOperations.VMOperator) {
	logrus.Printf("Shutting down VM '%s'", vmConfig.VMName)
	if err := operator.Shutdown(vmConfig.VMName); err != nil {
		logrus.Fatalf("Job failed: error shutting down VM '%s': %v", vmConfig.VMName, err)
	}
	logrus.Println("VM shut down successfully!")
}
