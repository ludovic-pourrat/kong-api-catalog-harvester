package factories

import (
	"fmt"
	spec "github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/utils"
	"strings"
	"unicode"
)

type PathParam struct {
	*spec.Parameter
}

func CreateParameterizedPath(path string) string {
	var ParameterizedPathParts []string
	pathParts := strings.Split(path, "/")

	for idx, part := range pathParts {
		// if part is a suspect param, replace it with a param name, otherwise do nothing
		if isSuspectPathParam(part) && !utils.IsPathParam(part) {
			var name string
			if idx < 1 {
				name = getParamName(part)
			} else {
				if utils.IsPathParam(pathParts[idx-1]) {
					name = "name-" + getParamName(part)
				} else {
					name = pathParts[idx-1] + "-" + getParamName(part)
				}
			}
			ParameterizedPathParts = append(ParameterizedPathParts, fmt.Sprintf("{%s}", name))
		} else {
			ParameterizedPathParts = append(ParameterizedPathParts, part)
		}
	}
	parameterizedPath := strings.Join(ParameterizedPathParts, "/")
	return parameterizedPath
}

func getParamName(pathPart string) string {
	if isNumber(pathPart) {
		return "id"
	} else if isUUID(pathPart) {
		return "uuid"
	} else if isToken(pathPart) {
		return "token"
	}
	return pathPart
}

func getParamSchema(pathPart string) *spec.Schema {
	if isNumber(pathPart) {
		return spec.NewInt64Schema()
	} else if isUUID(pathPart) {
		return spec.NewUUIDSchema()
	} else if isToken(pathPart) {
		return spec.NewStringSchema()
	}
	return spec.NewStringSchema()
}

func isSuspectPathParam(pathPart string) bool {
	if isNumber(pathPart) {
		return true
	}
	if isUUID(pathPart) {
		return true
	}
	if isToken(pathPart) {
		return true
	}
	return false
}

func isNumber(pathPart string) bool {
	return countDigitsInString(pathPart) > 0
}

func isUUID(pathPart string) bool {
	_, err := uuid.Parse(pathPart)
	return err == nil
}

func isToken(pathPart string) bool {
	const maxLen = 32
	const minDigitsLen = 1
	if len(pathPart) <= maxLen {
		return false
	}
	return countDigitsAndLetterInString(pathPart) > minDigitsLen
}

func countDigitsAndLetterInString(s string) int {
	count := 0
	for _, c := range s {
		if unicode.IsNumber(c) || unicode.IsLetter(c) {
			count++
		}
	}
	return count
}

func countDigitsInString(s string) int {
	count := 0
	for _, c := range s {
		if unicode.IsNumber(c) {
			count++
		}
	}
	return count
}
