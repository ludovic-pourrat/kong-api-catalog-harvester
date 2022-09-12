package factories

import "github.com/getkin/kin-openapi/openapi3"

func BuildSpecification(name string, version string) *openapi3.T {

	info := &openapi3.Info{
		Title:       name,
		Version:     "0.0.0-snapshot",
		Description: "",
	}

	return &openapi3.T{OpenAPI: version, Info: info}
}
