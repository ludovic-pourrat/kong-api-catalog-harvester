package factories

import "github.com/getkin/kin-openapi/openapi3"

func BuildSpecification(name string, version string) *openapi3.T {

	info := &openapi3.Info{
		Title:       name,
		Version:     "0.0.0-snapshot",
		Description: "API harvested from live traffic",
	}

	return &openapi3.T{
		OpenAPI: version,
		Info:    info,
		Paths:   openapi3.Paths{},
	}
}

func AggregateSpecification(specification *openapi3.T,
	paths map[string]string,
	methods map[string]string,
	operations map[string]*openapi3.Operation) *openapi3.T {

	clone := openapi3.T{}
	clone = *specification
	for name, path := range paths {
		operation := operations[name]
		clone.AddOperation(path, methods[name], operation)
		AddPath(clone.Paths, path, methods[name], operation)
	}

	return &clone
}
