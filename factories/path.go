package factories

import (
	"github.com/getkin/kin-openapi/openapi3"
	"net/http"
)

func AddPath(paths openapi3.Paths, path string, method string, operation *openapi3.Operation) {
	pathItem := paths.Find(path)
	if pathItem == nil {
		newPathItem := &openapi3.PathItem{}
		AddOperationToPathItem(newPathItem, method, operation)
		paths[path] = newPathItem
	} else {
		AddOperationToPathItem(pathItem, method, operation)
	}
}

func AddOperationToPathItem(item *openapi3.PathItem, method string, operation *openapi3.Operation) {
	switch method {
	case http.MethodGet:
		item.Get = operation
	case http.MethodDelete:
		item.Delete = operation
	case http.MethodOptions:
		item.Options = operation
	case http.MethodPatch:
		item.Patch = operation
	case http.MethodHead:
		item.Head = operation
	case http.MethodPost:
		item.Post = operation
	case http.MethodPut:
		item.Put = operation
	}
}
