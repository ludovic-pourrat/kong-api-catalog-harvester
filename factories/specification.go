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

func CloneSpecification(specification *openapi3.T, paths pathtrie.PathTrie, methods map[string]string, operations map[string]*openapi3.Operation) *openapi3.T {

	clone := BuildSpecification(specification.Info.Title, specification.OpenAPI)
	UpdateSpecification(clone, paths, methods, operations)

	return clone
}

func UpdateSpecification(specification *openapi3.T, paths pathtrie.PathTrie, methods map[string]string, operations map[string]*openapi3.Operation) {

	//for name, path := range paths.Trie {
	//	if name == "" {
	//		continue
	//	}
	//	operation := operations[name]
	//	specification.AddOperation(path.FullPath, methods[name], operation)
	//	AddPath(specification.Paths, path.FullPath, methods[name], operation)
	//}

}
