package validation

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func ValidateData(data interface{}) error {
	schema := gojsonschema.NewStringLoader(`{
		"type": "object",
		"required": ["platform", "timestamp"],
		"properties": {
			"platform": {"type": "string", "enum": ["linkedin", "youtube", "instagram"]},
			"timestamp": {"type": "string"}
		}
	}`)
	dataLoader := gojsonschema.NewGoLoader(data)
	result, err := gojsonschema.Validate(schema, dataLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return fmt.Errorf("validation failed: %v", result.Errors())
	}
	return nil
}