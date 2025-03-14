package jobs

import (
	"log"
	"time"

	"vnecro/config"
	"vnecro/vmOperations"
)

func ExecuteShellCommand(vmConfig *config.VMConfig, op config.Operation, pipeline map[string]string, operator vmOperations.VMOperator) {
	// Determine which role to use (default to "user" if not specified).
	role := op.Role
	if role == "" {
		role = "user"
	}
	credentials, err := config.GetUserByRole(vmConfig, role)
	if err != nil {
		log.Fatalf("Job failed: %v", err)
	}

	// Wait until the guest execution service is ready.

	if err := operator.WaitForGuestExecReady(vmConfig.VMName, credentials.Username, credentials.Password, 60*time.Second); err != nil {
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

	output, err := operator.ExecuteShellCommand(vmConfig.VMName, credentials.Username, credentials.Password, cmdStr, args...)
	if err != nil {
		log.Fatalf("Job failed: error executing shell command: %v", err)
	}
	log.Printf("Shell command output: %s", output)

	// If "store_as" is specified, store the output in the pipeline.
	if op.StoreAs != "" {
		pipeline[op.StoreAs] = output
		log.Printf("Stored output in variable '%s'", op.StoreAs)
	}
}
