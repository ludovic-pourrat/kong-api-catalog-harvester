package utils

import (
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/factories"
	"strings"
)

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
