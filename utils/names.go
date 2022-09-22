package utils

import (
	"fmt"
	"strings"
)

func GetName(method string, paths []string) string {
	name := method
	for _, path := range paths {
		if len(path) > 0 {
			if !IsPathParam(path) {
				name += "-" + path
			} else {
				name += "-by-" + GetParam(path)
			}
		}
	}
	return strings.ToLower(name)
}

func IsPathParam(segment string) bool {
	return strings.HasPrefix(segment, "{") &&
		strings.HasSuffix(segment, "}")
}

func GetParam(segment string) string {
	return strings.TrimPrefix(strings.TrimSuffix(segment, "}"), "{")
}

func GenerateParamName(idx int, parts []string) string {
	var name string
	if idx < 1 {
		name = "nane"
	} else {
		if IsPathParam(parts[idx-1]) {
			name = parts[idx-1] + "-name"
		} else {
			name = GetParam(parts[idx-1]) + "-name"
		}
	}
	return fmt.Sprintf("{%s}", name)
}
