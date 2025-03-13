package jobs

import (
	"log"
	"time"

	"vbnecro/config"
	"vbnecro/vboxOperations"
)

func ExecuteShellCommand(vmConfig *config.VMConfig, op config.Operation, pipeline map[string]string) {
	// Determine which role to use (default "user").
	role := op.Role
	if role == "" {
		role = "user"
	}
	credentials, err := config.GetUserByRole(vmConfig, role)
	if err != nil {
		log.Fatalf("Job failed: %v", err)
	}
	// Wait for guest execution service to be ready.
	if err := vboxOperations.WaitForGuestExecReady(vmConfig.VMName, credentials.Username, credentials.Password, 60*time.Second); err != nil {
		log.Fatalf("Job failed: guest execution service not ready on VM '%s': %v", vmConfig.VMName, err)
	}
	log.Printf("Guest execution service is ready on VM '%s'. Executing shell command...", vmConfig.VMName)

	// Retrieve command and arguments.
	cmdStr, ok := op.Params["command"].(string)
	if !ok || cmdStr == "" {
		log.Fatalf("Job failed: missing command parameter for ExecuteShellCommand")
	}
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
	// If a store_as value is provided, save the output in the pipeline.
	if op.StoreAs != "" {
		pipeline[op.StoreAs] = output
		log.Printf("Stored output in variable '%s'", op.StoreAs)
	}
}
