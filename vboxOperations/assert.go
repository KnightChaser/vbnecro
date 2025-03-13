package vboxOperations

import (
	"fmt"
	"strconv"
	"strings"
)

// Assert compares the given variableValue with the expected value according to the operator and type.
// valueType can be "int", "float", or "string" (default).
func Assert(variableValue string, operator string, expected string, valueType string) error {
	switch operator {
	case "equal":
		switch valueType {
		case "int":
			actualInt, err := strconv.Atoi(variableValue)
			if err != nil {
				return fmt.Errorf("failed to convert actual value '%s' to int: %v", variableValue, err)
			}
			expectedInt, err := strconv.Atoi(expected)
			if err != nil {
				return fmt.Errorf("failed to convert expected value '%s' to int: %v", expected, err)
			}
			if actualInt != expectedInt {
				return fmt.Errorf("assertion failed: expected %d, got %d", expectedInt, actualInt)
			}

		case "float":
			actualFloat, err := strconv.ParseFloat(variableValue, 64)
			if err != nil {
				return fmt.Errorf("failed to convert actual value '%s' to float: %v", variableValue, err)
			}
			expectedFloat, err := strconv.ParseFloat(expected, 64)
			if err != nil {
				return fmt.Errorf("failed to convert expected value '%s' to float: %v", expected, err)
			}
			if actualFloat != expectedFloat {
				return fmt.Errorf("assertion failed: expected %f, got %f", expectedFloat, actualFloat)
			}

		case "string":
			if variableValue != expected {
				return fmt.Errorf("assertion failed: expected '%s', got '%s'", expected, variableValue)
			}

		default:
			return fmt.Errorf("unknown value type: %s, only support 'int', 'float', 'string'", valueType)
		}

	case "includes":
		if !strings.Contains(variableValue, expected) {
			return fmt.Errorf("assertion failed: expected '%s' to include '%s'", variableValue, expected)
		}

	case "greater":
		switch valueType {
		case "int":
			actualInt, err := strconv.Atoi(variableValue)
			if err != nil {
				return fmt.Errorf("failed to convert actual value '%s' to int: %v", variableValue, err)
			}
			expectedInt, err := strconv.Atoi(expected)
			if err != nil {
				return fmt.Errorf("failed to convert expected value '%s' to int: %v", expected, err)
			}
			if actualInt <= expectedInt {
				return fmt.Errorf("assertion failed: expected greater than %d, got %d", expectedInt, actualInt)
			}

		case "float":
			actualFloat, err := strconv.ParseFloat(variableValue, 64)
			if err != nil {
				return fmt.Errorf("failed to convert actual value '%s' to float: %v", variableValue, err)
			}
			expectedFloat, err := strconv.ParseFloat(expected, 64)
			if err != nil {
				return fmt.Errorf("failed to convert expected value '%s' to float: %v", expected, err)
			}
			if actualFloat <= expectedFloat {
				return fmt.Errorf("assertion failed: expected greater than %f, got %f", expectedFloat, actualFloat)
			}

		default:
			return fmt.Errorf("assertion 'greater' requires numeric(int, float) type")
		}

	case "smaller":
		switch valueType {
		case "int":
			actualInt, err := strconv.Atoi(variableValue)
			if err != nil {
				return fmt.Errorf("failed to convert actual value '%s' to int: %v", variableValue, err)
			}
			expectedInt, err := strconv.Atoi(expected)
			if err != nil {
				return fmt.Errorf("failed to convert expected value '%s' to int: %v", expected, err)
			}
			if actualInt >= expectedInt {
				return fmt.Errorf("assertion failed: expected smaller than %d, got %d", expectedInt, actualInt)
			}

		case "float":
			actualFloat, err := strconv.ParseFloat(variableValue, 64)
			if err != nil {
				return fmt.Errorf("failed to convert actual value '%s' to float: %v", variableValue, err)
			}
			expectedFloat, err := strconv.ParseFloat(expected, 64)
			if err != nil {
				return fmt.Errorf("failed to convert expected value '%s' to float: %v", expected, err)
			}
			if actualFloat >= expectedFloat {
				return fmt.Errorf("assertion failed: expected smaller than %f, got %f", expectedFloat, actualFloat)
			}

		default:
			return fmt.Errorf("assertion 'smaller' requires numeric(int, float) type")
		}
	default:
		return fmt.Errorf("unknown operator: %s", operator)
	}
	return nil
}

// RunAssert retrieves the stored variable from the pipeline and performs the assertion.
// pipeline is a map from variable names to their stored string outputs.
func RunAssert(pipeline map[string]string, variable, operator, expected, valueType string) error {
	value, ok := pipeline[variable]

	if !ok {
		return fmt.Errorf("variable '%s' not found in pipeline", variable)
	}
	return Assert(value, operator, expected, valueType)
}
