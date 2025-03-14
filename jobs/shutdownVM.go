package jobs

import (
	"log"

	"vnecro/config"
	"vnecro/vmOperations"
)

func ShutdownVM(vmConfig *config.VMConfig, operator vmOperations.VMOperator) {
	log.Printf("Shutting down VM '%s'", vmConfig.VMName)
	if err := operator.Shutdown(vmConfig.VMName); err != nil {
		log.Fatalf("Job failed: error shutting down VM '%s': %v", vmConfig.VMName, err)
	}
	log.Println("VM shut down successfully!")
}
