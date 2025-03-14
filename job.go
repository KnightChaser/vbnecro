package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"vnecro/config"
	"vnecro/jobs"
	"vnecro/vmOperations"
)

// ProcessJobs iterates over each job in the configuration, executing operations.
// If an operation fails or if the user interrupts (CTRL+C), the current job is
// considered failed, and if a rollback snapshot is specified, the VM is rolled back.
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

	// Variables to track the current job and VM.
	var currentJob *config.JobConfig
	var currentVM *config.VMConfig

	// Set up a channel to listen for CTRL+C (SIGINT).
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		sig := <-sigChan
		logrus.Infof("Received signal: %v. Treating current job as failed.", sig)
		// If a current job is in progress and a rollback is specified, trigger rollback.
		if currentJob != nil && currentJob.RollbackOnFailure != "" && currentVM != nil {
			err := jobs.RollbackVM(currentVM, currentJob.RollbackOnFailure, operator)
			if err != nil {
				logrus.Errorf("Rollback failed on VM '%s': %v", currentVM.VMName, err)
			} else {
				logrus.Infof("Rollback successful on VM '%s'", currentVM.VMName)
			}
		}
		logrus.Warn("Program interrupted. Exiting now.")
		os.Exit(1)
	}()

	// Process each job.
	for _, job := range cfg.Jobs {
		currentJob = &job

		vmConfig, err := config.GetVMConfig(cfg.VMs, job.VMAlias)
		if err != nil {
			logrus.Errorf("Job for VM alias '%s' failed: %v", job.VMAlias, err)
			continue
		}
		currentVM = vmConfig

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
				// Stop processing further operations in this job.
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
