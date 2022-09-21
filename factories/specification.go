package factories

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/utils/pathtrie"
)

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

func CloneSpecification(specification *openapi3.T, paths pathtrie.PathTrie) *openapi3.T {

	clone := BuildSpecification(specification.Info.Title, specification.OpenAPI)
	UpdateSpecification(clone, paths)
	return clone
}

func UpdateSpecification(specification *openapi3.T, paths pathtrie.PathTrie) {

	for _, path := range paths.Nodes() {
		if len(path.Children) == 0 {
			params := BuildParamsPath(path.URL, path.Path)
			for _, param := range params {
				if path.Operation.Parameters.GetByInAndName("path", param.Value.Name) == nil {
					path.Operation.Parameters = append(path.Operation.Parameters, param)
				}
			}
			specification.AddOperation(path.URL, path.Method, path.Operation)
			AddPath(specification.Paths, path.URL, path.Method, path.Operation)
		}
	}

}
