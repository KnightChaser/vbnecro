package jobs

import (
	"log"

	"vbnecro/config"
	"vbnecro/vboxOperations"
)

func StartVM(vmConfig *config.VMConfig) {
	log.Printf("Starting VM '%s'", vmConfig.VMName)
	if err := vboxOperations.StartVM(vmConfig.VMName); err != nil {
		log.Fatalf("Job failed: error starting VM '%s': %v", vmConfig.VMName, err)
	}
	log.Println("VM started successfully!")
}
