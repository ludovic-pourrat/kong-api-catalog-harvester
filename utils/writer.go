package utils

import (
	"context"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/invopop/yaml"
	"os"
)

func Write(name string, specification *openapi3.T) error {
	// validate
	err := specification.Validate(context.Background())
	if err != nil {
		return err
	}
	// marshal to yaml
	data, err := yaml.Marshal(specification)
	if err != nil {
		return err
	}
	// write to file
	os.WriteFile(fmt.Sprintf("/logs/%s.yaml", name), data, 0644)

	return nil
}
