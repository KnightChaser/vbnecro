package jobs

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"vnecro/config"
	"vnecro/vmOperations"
)

// ExecuteShellCommand executes a shell command on the given VM using the provided operator.
// It waits for the guest execution service to be ready, retrieves the command and arguments,
// executes the command, and optionally prints and stores the output in the pipeline.
// Returns an error if any step fails.
func ExecuteShellCommand(vmConfig *config.VMConfig, op config.Operation, pipeline map[string]string, operator vmOperations.VMOperator) error {
	// Determine which role to use (default to "user" if not specified).
	role := op.Role
	if role == "" {
		role = "user"
	}
	credentials, err := config.GetUserByRole(vmConfig, role)
	if err != nil {
		return fmt.Errorf("error retrieving user for role '%s': %w", role, err)
	}

	// Wait until the guest execution service is ready.
	if err := operator.WaitForGuestExecReady(vmConfig.VMName, credentials.Username, credentials.Password, 60*time.Second); err != nil {
		return fmt.Errorf("guest execution service not ready on VM '%s': %w", vmConfig.VMName, err)
	}
	logrus.Infof("Guest execution service is ready on VM '%s'. Executing shell command...", vmConfig.VMName)

	// Retrieve command from parameters.
	cmdStr, ok := op.Params["command"].(string)
	if !ok || cmdStr == "" {
		return fmt.Errorf("missing command parameter for ExecuteShellCommand")
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

	// Execute the shell command.
	output, err := operator.ExecuteShellCommand(vmConfig.VMName, credentials.Username, credentials.Password, cmdStr, args...)
	if err != nil {
		return fmt.Errorf("error executing shell command: %w", err)
	}

	// If PrintOutput is true, print a structured log.
	if op.PrintOutput {
		fullCommand := cmdStr
		if len(args) > 0 {
			fullCommand += " " + strings.Join(args, " ")
		}
		logrus.Infof("Shell command executed on VM '%s':", vmConfig.VMName)
		logrus.Infof(" - Executed command: %s", fullCommand)
		logrus.Infof(" - Executed command result:\n%s", boxOutput(output))
	}

	// If "store_as" is specified, store the output in the pipeline.
	if op.StoreAs != "" {
		pipeline[op.StoreAs] = output
		logrus.Infof("Stored output in variable '%s'", op.StoreAs)
	}

	return nil
}

// boxOutput returns the given text surrounded by an ASCII box.
func boxOutput(output string) string {
	lines := strings.Split(output, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Create a top and bottom border.
	border := "+" + strings.Repeat("-", maxLen+2) + "+"
	var sb strings.Builder
	sb.WriteString(border + "\n")
	for _, line := range lines {
		sb.WriteString(fmt.Sprintf("| %-*s |\n", maxLen, line))
	}
	sb.WriteString(border)
	return sb.String()
}
