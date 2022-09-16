package factories

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/types"
	"github.com/xeipuuv/gojsonschema"
	"net/http"
	"strconv"
)

func BuildResponse(raw string, log types.Log) *openapi3.ResponseRef {
	var contentType string
	response := openapi3.NewResponse()
	if log.Response.Status != 204 {
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
		response.WithContent(responseContent)
	}
	response.WithDescription(http.StatusText(log.Response.Status))
	responseRef := &openapi3.ResponseRef{
		Value: response,
	}
	return responseRef
}

func BuildResponses(raw string, log types.Log) openapi3.Responses {
	responses := make(map[string]*openapi3.ResponseRef)
	responses[strconv.Itoa(log.Response.Status)] = BuildResponse(raw, log)
	return responses
}

func MergeResponse(raw string, log types.Log, responseRef *openapi3.ResponseRef) (bool, error) {
	var contentType string
	var updated = false
	response := responseRef.Value
	if log.Response.Status != 204 {
		if _, found := log.Response.Headers["content-type"]; found {
			contentType = log.Response.Headers["content-type"].(string)
		} else {
			contentType = "application/json"
		}
		respBodyJSON, _ := gojsonschema.NewStringLoader(raw).LoadJSON()
		responseContent := response.Content.Get(contentType)
		schema, err := MergeSchema(respBodyJSON, responseContent.Schema.Value)
		if err != nil {
			return false, err
		}
		responseContent.Schema.Value = schema
		updated = true
	}
	return updated, nil
}

func MergeResponses(operation *openapi3.Operation, raw string, log types.Log) (bool, error) {
	var updated = false
	var err error
	responses := operation.Responses
	if responses[strconv.Itoa(log.Response.Status)] == nil {
		responses[strconv.Itoa(log.Response.Status)] = BuildResponse(raw, log)
		updated = true
	} else {
		updated, err = MergeResponse(raw, log, responses[strconv.Itoa(log.Response.Status)])
		if err != nil {
			return updated, err
		}
	}
	return updated, nil
}
