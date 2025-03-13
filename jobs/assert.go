package jobs

import (
	"log"

	"vbnecro/config"
	"vbnecro/vboxOperations"
)

func Assert(pipeline map[string]string, op config.Operation) {
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
	valueType := "string"
	if vt, ok := op.Params["type"].(string); ok && vt != "" {
		valueType = vt
	}
	if err := vboxOperations.RunAssert(pipeline, varName, operator, expectedVal, valueType); err != nil {
		log.Fatalf("Job failed: assertion error for variable '%s': %v", varName, err)
	}
	log.Printf("Assertion passed for variable '%s'", varName)
}
