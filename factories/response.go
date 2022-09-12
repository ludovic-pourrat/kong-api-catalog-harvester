package factories

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/types"
	"github.com/xeipuuv/gojsonschema"
	"strconv"
)

func BuildResponse(raw string, log types.Log) *openapi3.ResponseRef {
	var contentType string
	if _, found := log.Response.Headers["content-type"]; found {
		contentType = log.Response.Headers["content-type"].(string)
	} else {
		contentType = "application/json"
	}
	respBodyJSON, _ := gojsonschema.NewStringLoader(raw).LoadJSON()
	responseSchema, _ := BuildSchema(respBodyJSON)
	responseContent := openapi3.Content{
		contentType: openapi3.NewMediaType().WithSchema(responseSchema),
	}
	response := openapi3.NewResponse()
	response.WithContent(responseContent)
	response.WithDescription("")
	responseRef := &openapi3.ResponseRef{
		Value: response,
	}
	return responseRef
}

func BuildResponses(raw string, log types.Log) openapi3.Responses {
	responses := openapi3.NewResponses()
	responses[strconv.Itoa(log.Response.Status)] = BuildResponse(raw, log)
	return responses
}

func MergeResponses(operation *openapi3.Operation, raw string, log types.Log) bool {
	responses := operation.Responses
	if operation.Responses.Get(log.Response.Status) == nil {
		responses[strconv.Itoa(log.Response.Status)] = BuildResponse(raw, log)
	}
	return true
}
