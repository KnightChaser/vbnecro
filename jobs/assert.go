package jobs

import (
	"github.com/sirupsen/logrus"
	"vnecro/config"
	"vnecro/vboxOperations"
)

// Assert retrieves the variable from the pipeline and verifies it using the given operator.
// It halts execution if the assertion fails.
func Assert(pipeline map[string]string, op config.Operation) {
	// Retrieve the variable name.
	varName, ok := op.Params["variable"].(string)
	if !ok || varName == "" {
		logrus.Fatalf("Job failed: missing 'variable' parameter for Assert operation")
	}

	// Retrieve the operator.
	operator, ok := op.Params["operator"].(string)
	if !ok || operator == "" {
		logrus.Fatalf("Job failed: missing 'operator' parameter for Assert operation")
	}

	// Retrieve the expected value.
	expectedVal, ok := op.Params["expected"].(string)
	if !ok {
		logrus.Fatalf("Job failed: missing 'expected' parameter for Assert operation")
	}

	// Optional: Retrieve the value type; default is "string".
	valueType := "string"
	if vt, ok := op.Params["type"].(string); ok && vt != "" {
		valueType = vt
	}

	// Run the assertion using the pipeline value.
	if err := vboxOperations.RunAssert(pipeline, varName, operator, expectedVal, valueType); err != nil {
		logrus.Fatalf("Job failed: assertion error for variable '%s': %v", varName, err)
	}

	logrus.Infof("Assertion passed for variable '%s'", varName)
}
