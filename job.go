package main

import (
	"log"
	"time"

	"vbnecro/vboxOperations"
)

func ProcessJobs(cfg *Config) {
	for _, job := range cfg.Jobs {
		vmConfig, err := GetVMConfig(cfg.VMs, job.VMAlias)
		if err != nil {
			log.Fatalf("Job for VM alias '%s' failed: %v", job.VMAlias, err)
		}

		// Ensure the VM is turned off if specified.
		if job.EnsureOff {
			log.Printf("Ensuring VM '%s' is off before executing operations", vmConfig.VMName)
			if err := vboxOperations.ShutdownVM(vmConfig.VMName); err != nil {
				log.Fatalf("Job failed: error shutting down VM '%s': %v", vmConfig.VMName, err)
			}
			log.Printf("VM '%s' shut down successfully.", vmConfig.VMName)
		}

		// Process each operation.
		for _, op := range job.Operations {
			switch op.Type {
			case "RestoreSnapshot":
				// (RestoreSnapshot logic remains the same)
				// ...

			case "StartVM":
				log.Printf("Starting VM '%s'", vmConfig.VMName)
				if err := vboxOperations.StartVM(vmConfig.VMName); err != nil {
					log.Fatalf("Job failed: error starting VM '%s': %v", vmConfig.VMName, err)
				}
				log.Println("VM started successfully!")

			case "PauseVM":
				log.Printf("Pausing VM '%s'", vmConfig.VMName)
				if err := vboxOperations.PauseVM(vmConfig.VMName); err != nil {
					log.Fatalf("Job failed: error pausing VM '%s': %v", vmConfig.VMName, err)
				}
				log.Println("VM paused successfully!")

			case "ShutdownVM":
				log.Printf("Shutting down VM '%s'", vmConfig.VMName)
				if err := vboxOperations.ShutdownVM(vmConfig.VMName); err != nil {
					log.Fatalf("Job failed: error shutting down VM '%s': %v", vmConfig.VMName, err)
				}
				log.Println("VM shut down successfully!")

			case "ExecuteShellCommand":
				log.Printf("Preparing to execute shell command on VM '%s'", vmConfig.VMName)
				// Wait until the guest execution service is ready.
				if err := vboxOperations.WaitForGuestExecReady(vmConfig.VMName, vmConfig.Username, vmConfig.Password, 60*time.Second); err != nil {
					log.Fatalf("Job failed: guest execution service not ready on VM '%s': %v", vmConfig.VMName, err)
				}
				log.Printf("Guest execution service is ready on VM '%s'. Executing shell command...", vmConfig.VMName)
				// Retrieve command.
				cmdStr, ok := op.Params["command"].(string)
				if !ok || cmdStr == "" {
					log.Fatalf("Job failed: missing command parameter for ExecuteShellCommand")
				}
				// Retrieve optional arguments.
				var args []string
				if rawArgs, exists := op.Params["args"]; exists {
					if slice, ok := rawArgs.([]interface{}); ok {
						for _, item := range slice {
							if str, ok := item.(string); ok {
								args = append(args, str)
							}
						}
					}
				}
				output, err := vboxOperations.ExecuteShellCommand(vmConfig.VMName, vmConfig.Username, vmConfig.Password, cmdStr, args...)
				if err != nil {
					log.Fatalf("Job failed: error executing shell command: %v", err)
				}
				log.Printf("Shell command output: %s", output)

			default:
				log.Fatalf("Job failed: unknown operation type: %s", op.Type)
			}
		}
	}
}
