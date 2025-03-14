package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"vnecro/config"
	"vnecro/jobs"
	"vnecro/vmOperations"
)

// ProcessJobs iterates over each job in the configuration, executing operations
// and, in case of any error, rolls back the VM to the specified snapshot if configured.
func ProcessJobs(cfg *config.Config) {
	// Create an instance of the VM operator.
	var operator vmOperations.VMOperator
	if cfg.VMManager == "virtualbox" {
		operator = vmOperations.NewVirtualBoxOperator()
	} else {
		logrus.Fatalf("Unsupported VM manager: %s (only virtualbox is supported)", cfg.VMManager)
	}

	// Pipeline to hold outputs from shell commands.
	pipeline := make(map[string]string)

	// Process each job.
	for _, job := range cfg.Jobs {
		vmConfig, err := config.GetVMConfig(cfg.VMs, job.VMAlias)
		if err != nil {
			logrus.Errorf("Job for VM alias '%s' failed: %v", job.VMAlias, err)
			continue
		}

		// If ensure_off is true, shut down the VM before processing operations.
		if job.EnsureOff {
			logrus.Infof("Ensuring VM '%s' is off", vmConfig.VMName)
			if err := jobs.ShutdownVM(vmConfig, operator); err != nil {
				logrus.Errorf("Failed to shut down VM '%s': %v", vmConfig.VMName, err)
				continue
			}
			logrus.Infof("VM '%s' shut down successfully.", vmConfig.VMName)
		}

		// Process each operation; if one fails, mark the job as failed.
		jobFailed := false
		for _, op := range job.Operations {
			var opErr error
			switch op.Type {
			case "RestoreSnapshot":
				opErr = jobs.RestoreSnapshot(vmConfig, op, operator)
			case "StartVM":
				opErr = jobs.StartVM(vmConfig, operator)
			case "PauseVM":
				opErr = jobs.PauseVM(vmConfig, operator)
			case "ShutdownVM":
				opErr = jobs.ShutdownVM(vmConfig, operator)
			case "ExecuteShellCommand":
				opErr = jobs.ExecuteShellCommand(vmConfig, op, pipeline, operator)
			case "Assert":
				opErr = jobs.Assert(pipeline, op)
			default:
				opErr = fmt.Errorf("unknown operation type: %s", op.Type)
			}

			if opErr != nil {
				logrus.Errorf("Operation %s failed: %v", op.Type, opErr)
				jobFailed = true
				// Break out of the operations loop on first error.
				break
			}
		}

		// If any operation failed and a rollback snapshot is specified, perform rollback.
		if jobFailed && job.RollbackOnFailure != "" {
			logrus.Infof("Job failed; initiating rollback on VM '%s' to snapshot '%s'",
				vmConfig.VMName, job.RollbackOnFailure)
			if err := jobs.RollbackVM(vmConfig, job.RollbackOnFailure, operator); err != nil {
				logrus.Errorf("Rollback failed on VM '%s': %v", vmConfig.VMName, err)
			} else {
				logrus.Infof("Rollback successful on VM '%s'", vmConfig.VMName)
			}
		}
	}
}
