package factories

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/shared"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/types"
	"github.com/xeipuuv/gojsonschema"
)

func BuildRequest(raw []byte, contentType string, log types.Log) *openapi3.RequestBodyRef {
	var requestBodyRef *openapi3.RequestBodyRef
	if log.Request.Method != "GET" && log.Request.Method != "DELETE" {
		request := openapi3.NewRequestBody()
		reqBodyJSON, _ := gojsonschema.NewBytesLoader(raw).LoadJSON()
		requestSchema, _ := shared.GetSchema(reqBodyJSON)
		requestContent := openapi3.Content{
			contentType: openapi3.NewMediaType().WithSchema(requestSchema),
		}
		request.Content = requestContent
		request.Description = ""
		requestBodyRef = &openapi3.RequestBodyRef{
			Value: request,
		}
	} else {
		requestBodyRef = nil
	}
	return requestBodyRef
}

func MergeRequest(operation *openapi3.Operation, raw []byte, contentType string, log types.Log) bool {
	var requestBodyRef *openapi3.RequestBodyRef
	operation.RequestBody = requestBodyRef
	return true
}
