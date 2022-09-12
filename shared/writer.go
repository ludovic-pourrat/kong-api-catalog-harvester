package shared

import (
	"context"
	"fmt"
	"github.com/Kong/go-pdk/log"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/invopop/yaml"
	"os"
)

func Write(name string, specification *openapi3.T, logger log.Log) {
	// validate
	err := specification.Validate(context.Background())
	if err != nil {
		logger.Warn(err)
	}
	// marshal to yaml
	data, err := yaml.Marshal(specification)
	if err != nil {
		logger.Err(err)
		return
	}
	// write to file
	os.WriteFile(fmt.Sprintf("/logs/%s.yaml", name), data, 0644)
}
