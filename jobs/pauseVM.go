package jobs

import (
	"log"

	"vbnecro/config"
	"vbnecro/vboxOperations"
)

func PauseVM(vmConfig *config.VMConfig) {
	log.Printf("Pausing VM '%s'", vmConfig.VMName)
	if err := vboxOperations.PauseVM(vmConfig.VMName); err != nil {
		log.Fatalf("Job failed: error pausing VM '%s': %v", vmConfig.VMName, err)
	}
	log.Println("VM paused successfully!")
}
