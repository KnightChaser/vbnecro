package main

import (
	"log"
	"time"

	"vbnecro/vboxOperations"
)

func ProcessJobs(cfg *Config) {
	// Pipeline to hold variable outputs from ExecuteShellCommand operations.
	pipeline := make(map[string]string)

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
				log.Printf("Listing snapshots for VM '%s'", vmConfig.VMName)
				output, err := vboxOperations.ListSnapshots(vmConfig.VMName)
				if err != nil {
					log.Fatalf("Job failed: error listing snapshots for VM '%s': %v", vmConfig.VMName, err)
				}
				log.Println("Snapshot list output:")
				log.Println(output)

				var snapshotToRestore string
				if val, ok := op.Params["snapshot"].(string); ok && val != "" {
					snapshotToRestore = val
				}
				if snapshotToRestore == "" {
					snapshotToRestore, err = vboxOperations.ParseSnapshot(output)
					if err != nil {
						log.Fatalf("Job failed: error parsing snapshot for VM '%s': %v", vmConfig.VMName, err)
					}
				}
				log.Printf("Restoring VM '%s' to snapshot '%s'", vmConfig.VMName, snapshotToRestore)
				if err := vboxOperations.RestoreSnapshot(vmConfig.VMName, snapshotToRestore); err != nil {
					log.Fatalf("Job failed: error restoring snapshot for VM '%s': %v", vmConfig.VMName, err)
				}
				log.Println("Snapshot restored successfully!")

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
				// Determine which role to use (default to "user" if not specified).
				role := op.Role
				if role == "" {
					role = "user"
				}
				credentials, err := GetUserByRole(vmConfig, role)
				if err != nil {
					log.Fatalf("Job failed: %v", err)
				}
				// Wait until the guest execution service is ready.
				if err := vboxOperations.WaitForGuestExecReady(vmConfig.VMName, credentials.Username, credentials.Password, 60*time.Second); err != nil {
					log.Fatalf("Job failed: guest execution service not ready on VM '%s': %v", vmConfig.VMName, err)
				}
				log.Printf("Guest execution service is ready on VM '%s'. Executing shell command...", vmConfig.VMName)

				// Retrieve command from parameters.
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
				output, err := vboxOperations.ExecuteShellCommand(vmConfig.VMName, credentials.Username, credentials.Password, cmdStr, args...)
				if err != nil {
					log.Fatalf("Job failed: error executing shell command: %v", err)
				}
				log.Printf("Shell command output: %s", output)
				// Store the output in the pipeline if "store_as" is specified.
				if op.StoreAs != "" {
					pipeline[op.StoreAs] = output
					log.Printf("Stored output in variable '%s'", op.StoreAs)
				}

			case "Assert":
				// Retrieve assertion parameters.
				varName, ok := op.Params["variable"].(string)
				if !ok || varName == "" {
					log.Fatalf("Job failed: missing 'variable' parameter for Assert operation")
				}
				operator, ok := op.Params["operator"].(string)
				if !ok || operator == "" {
					log.Fatalf("Job failed: missing 'operator' parameter for Assert operation")
				}
				expectedVal, ok := op.Params["expected"].(string)
				if !ok {
					log.Fatalf("Job failed: missing 'expected' parameter for Assert operation")
				}

				// Optional: value type conversion (default "string").
				valueType := "string"
				if vt, ok := op.Params["type"].(string); ok && vt != "" {
					valueType = vt
				}

				// Run the assert using the pipeline map.
				if err := vboxOperations.RunAssert(pipeline, varName, operator, expectedVal, valueType); err != nil {
					log.Fatalf("Job failed: assertion error for variable '%s': %v", varName, err)
				}

				// Success message.
				log.Printf("Assertion passed for variable '%s'", varName)

			default:
				log.Fatalf("Job failed: unknown operation type: %s", op.Type)
			}
		}
	}
}
