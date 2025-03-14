package jobs

import (
	"github.com/sirupsen/logrus"

	"vnecro/config"
	"vnecro/vmOperations"
)

func StartVM(vmConfig *config.VMConfig, operator vmOperations.VMOperator) {
	logrus.Printf("Starting VM '%s'", vmConfig.VMName)
	if err := operator.Start(vmConfig.VMName); err != nil {
		logrus.Fatalf("Job failed: error starting VM '%s': %v", vmConfig.VMName, err)
	}
	logrus.Println("VM started successfully!")
}
