package utils

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/factories"
	"strings"
)

func GetOperationName(operation *openapi3.Operation) string {
	return operation.OperationID
}

func GetName(method string, url string) string {
	pathParts := strings.Split(url, "/")
	name := method
	for _, path := range pathParts {
		if len(path) > 0 {
			if !factories.IsPathParam(path) {
				name += "-" + path
			} else {
				name += "-by-x"
			}
		}
	}
	return strings.ToLower(name)
}
