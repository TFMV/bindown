package bindown

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/qri-io/jsonschema"
)

//go:embed bindown.schema.json
var jsonSchemaText string

// validateConfig checks whether cfg meets the json schema.
func validateConfig(cfg []byte) error {
	cfgJSON, err := yaml.YAMLToJSON(cfg)
	if err != nil {
		return fmt.Errorf("config is not valid yaml (or json)")
	}
	var schema jsonschema.Schema
	err = json.Unmarshal([]byte(jsonSchemaText), &schema)
	if err != nil {
		return err
	}
	validationErrs, err := schema.ValidateBytes(context.TODO(), cfgJSON)
	if err != nil {
		return fmt.Errorf("unexpected error running jsonSchema.ValidateBytes: %v", err)
	}
	if len(validationErrs) == 0 {
		return nil
	}
	sort.Slice(validationErrs, func(i, j int) bool {
		return validationErrs[i].Error() < validationErrs[j].Error()
	})
	msgs := make([]string, len(validationErrs))
	for i, validationErr := range validationErrs {
		msgs[i] = validationErr.Error()
	}
	return fmt.Errorf("invalid config:\n%s", strings.Join(msgs, "\n"))
}
