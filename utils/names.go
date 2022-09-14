package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GetName(method string, url string) string {
	pathParts := strings.Split(url, "/")
	name := method
	for _, path := range pathParts {
		if len(path) > 0 {
			if !IsPathParam(path) {
				name += "-" + path
			} else {
				name += "-by-" + "x"
			}
		}
	}
	return strings.ToLower(name)
}

func IsPathParam(segment string) bool {
	return strings.HasPrefix(segment, "{") &&
		strings.HasSuffix(segment, "}")
}

const charset = "abcdefghijklmnopqrstuvwxyz"

func GenerateParamName() string {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return fmt.Sprintf("{%s}", string(b))
}
