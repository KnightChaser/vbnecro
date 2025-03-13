package jobs

import (
	"log"

	"vbnecro/config"
	"vbnecro/vmOperations"
)

func PauseVM(vmConfig *config.VMConfig, operator vmOperations.VMOperator) {
	log.Printf("Pausing VM '%s'", vmConfig.VMName)
	if err := operator.Pause(vmConfig.VMName); err != nil {
		log.Fatalf("Job failed: error pausing VM '%s': %v", vmConfig.VMName, err)
	}
	log.Println("VM paused successfully!")
}
