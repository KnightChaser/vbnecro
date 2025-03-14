package jobs

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"vnecro/config"
	"vnecro/vboxOperations"
)

// Assert retrieves the variable from the pipeline and verifies it using the given operator.
// Instead of halting execution immediately, it returns an error so that the caller (job dispatcher)
// can decide to trigger a rollback or other recovery measures.
func Assert(pipeline map[string]string, op config.Operation) error {
	// Retrieve the variable name.
	varName, ok := op.Params["variable"].(string)
	if !ok || varName == "" {
		return fmt.Errorf("missing 'variable' parameter for Assert operation")
	}

	// Retrieve the operator.
	operator, ok := op.Params["operator"].(string)
	if !ok || operator == "" {
		return fmt.Errorf("missing 'operator' parameter for Assert operation")
	}

	// Retrieve the expected value.
	expectedVal, ok := op.Params["expected"].(string)
	if !ok {
		return fmt.Errorf("missing 'expected' parameter for Assert operation")
	}

	// Optional: Retrieve the value type; default is "string".
	valueType := "string"
	if vt, ok := op.Params["type"].(string); ok && vt != "" {
		valueType = vt
	}

	// Run the assertion using the pipeline value.
	if err := vboxOperations.RunAssert(pipeline, varName, operator, expectedVal, valueType); err != nil {
		return fmt.Errorf("assertion error for variable '%s': %v", varName, err)
	}

	logrus.Infof("Assertion passed for variable '%s'", varName)
	return nil
}
