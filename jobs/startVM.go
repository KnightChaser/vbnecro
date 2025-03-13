package jobs

import (
	"log"

	"vbnecro/config"
	"vbnecro/vmOperations"
)

func StartVM(vmConfig *config.VMConfig, operator vmOperations.VMOperator) {
	log.Printf("Starting VM '%s'", vmConfig.VMName)
	if err := operator.Start(vmConfig.VMName); err != nil {
		log.Fatalf("Job failed: error starting VM '%s': %v", vmConfig.VMName, err)
	}
	log.Println("VM started successfully!")
}
