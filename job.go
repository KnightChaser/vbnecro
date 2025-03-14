package main

import (
	"log"

	"vnecro/config"
	"vnecro/jobs"
	"vnecro/vmOperations"
)

func ProcessJobs(cfg *config.Config) {
	// Create an instance of the VM operator.

	var operator vmOperations.VMOperator

	// Currently only VirtualBox is supported.
	if cfg.VMManager == "virtualbox" {
		operator = vmOperations.NewVirtualBoxOperator()
	} else {
		log.Fatalf("Unsupported VM manager: %s (only virtualbox is supported)", cfg.VMManager)
	}

	// Pipeline to hold outputs from shell commands.
	pipeline := make(map[string]string)

	for _, job := range cfg.Jobs {
		vmConfig, err := config.GetVMConfig(cfg.VMs, job.VMAlias)
		if err != nil {
			log.Fatalf("Job for VM alias '%s' failed: %v", job.VMAlias, err)
		}

		// If ensure_off is true, shut down the VM before processing operations.
		if job.EnsureOff {
			log.Printf("Ensuring VM '%s' is off", vmConfig.VMName)
			jobs.ShutdownVM(vmConfig, operator)
			log.Printf("VM '%s' shut down successfully.", vmConfig.VMName)
		}

		// Process each operation.
		for _, op := range job.Operations {
			switch op.Type {
			case "RestoreSnapshot":
				jobs.RestoreSnapshot(vmConfig, op, operator)
			case "StartVM":
				jobs.StartVM(vmConfig, operator)
			case "PauseVM":
				jobs.PauseVM(vmConfig, operator)
			case "ShutdownVM":
				jobs.ShutdownVM(vmConfig, operator)
			case "ExecuteShellCommand":
				jobs.ExecuteShellCommand(vmConfig, op, pipeline, operator)
			case "Assert":
				jobs.Assert(pipeline, op)
			default:
				log.Fatalf("Job failed: unknown operation type: %s", op.Type)
			}
		}
	}
}
