package factories

import (
	spec "github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

type PathParam struct {
	*spec.Parameter
}

func generateParamName() string {

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

var digitCheck = regexp.MustCompile(`^[0-9]+$`)

func CreateParameterizedPath(path string) string {
	var ParameterizedPathParts []string
	pathParts := strings.Split(path, "/")

	for _, part := range pathParts {
		// if part is a suspect param, replace it with a param name, otherwise do nothing
		if isSuspectPathParam(part) {
			paramName := generateParamName()
			ParameterizedPathParts = append(ParameterizedPathParts, "{"+paramName+"}")
		} else {
			ParameterizedPathParts = append(ParameterizedPathParts, part)
		}
	}

	parameterizedPath := strings.Join(ParameterizedPathParts, "/")

	return parameterizedPath
}

func getParamSchema(pathPart string) *spec.Schema {
	if isNumber(pathPart) {
		return spec.NewInt64Schema()
	} else if isUUID(pathPart) {
		return spec.NewUUIDSchema()
	} else if isMixed(pathPart) {
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
	if isMixed(pathPart) {
		return true
	}
	return false
}

func isNumber(pathPart string) bool {
	return digitCheck.MatchString(pathPart)
}

func isUUID(pathPart string) bool {
	_, err := uuid.Parse(pathPart)
	return err == nil
}

func isMixed(pathPart string) bool {
	const maxLen = 256
	const minDigitsLen = 2
	if len(pathPart) < maxLen {
		return false
	}
	return countDigitsInString(pathPart) > minDigitsLen
}

func countDigitsInString(s string) int {
	count := 0
	for _, c := range s {
		if unicode.IsNumber(c) || unicode.IsLetter(c) {
			count++
		}
	}
	return count
}
