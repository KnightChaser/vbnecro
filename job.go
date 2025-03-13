package main

import (
	"log"

	"vbnecro/config"
	"vbnecro/jobs"
)

func ProcessJobs(cfg *config.Config) {
	// Pipeline to hold variable outputs from ExecuteShellCommand operations.
	pipeline := make(map[string]string)

	for _, job := range cfg.Jobs {
		vmConfig, err := config.GetVMConfig(cfg.VMs, job.VMAlias)
		if err != nil {
			log.Fatalf("Job for VM alias '%s' failed: %v", job.VMAlias, err)
		}

		// If ensure_off is true, shut down the VM before processing operations.
		if job.EnsureOff {
			log.Printf("Ensuring VM '%s' is off before executing operations", vmConfig.VMName)
			jobs.ShutdownVM(vmConfig)
			log.Printf("VM '%s' shut down successfully.", vmConfig.VMName)
		}

		// Process each operation in the job.
		for _, op := range job.Operations {
			switch op.Type {
			case "RestoreSnapshot":
				jobs.RestoreSnapshot(vmConfig, op)
			case "StartVM":
				jobs.StartVM(vmConfig)
			case "PauseVM":
				jobs.PauseVM(vmConfig)
			case "ShutdownVM":
				jobs.ShutdownVM(vmConfig)
			case "ExecuteShellCommand":
				jobs.ExecuteShellCommand(vmConfig, op, pipeline)
			case "Assert":
				jobs.Assert(pipeline, op)
			default:
				log.Fatalf("Job failed: unknown operation type: %s", op.Type)
			}
		}
	}
}
