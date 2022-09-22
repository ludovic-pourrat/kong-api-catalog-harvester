package factories

import (
	"github.com/getkin/kin-openapi/openapi3"
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

func CloneSpecification(specification *openapi3.T, paths PathTrie) *openapi3.T {

	clone := BuildSpecification(specification.Info.Title, specification.OpenAPI)
	UpdateSpecification(clone, paths)
	return clone
}

func UpdateSpecification(specification *openapi3.T, paths PathTrie) {

	for _, path := range paths.Nodes() {
		// TODO merge params
		//params := BuildParamsPath(path.URL, path.Path)
		//for _, param := range params {
		//	if path.Operation.Parameters.GetByInAndName("path", param.Value.Name) == nil {
		//		path.Operation.Parameters = append(path.Operation.Parameters, param)
		//	}
		//}
		for method, operation := range path.Operations {
			specification.AddOperation(path.Path, method, operation)
			AddPath(specification.Paths, path.Path, method, operation)
		}

	}

}
