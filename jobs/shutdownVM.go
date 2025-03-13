package jobs

import (
	"log"

	"vbnecro/config"
	"vbnecro/vboxOperations"
)

func ShutdownVM(vmConfig *config.VMConfig) {
	log.Printf("Shutting down VM '%s'", vmConfig.VMName)
	if err := vboxOperations.ShutdownVM(vmConfig.VMName); err != nil {
		log.Fatalf("Job failed: error shutting down VM '%s': %v", vmConfig.VMName, err)
	}
	log.Println("VM shut down successfully!")
}
