package factories

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/types"
	"github.com/xeipuuv/gojsonschema"
	"mime"
)

func BuildRequest(raw []byte, contentType string, log types.Log) *openapi3.RequestBodyRef {
	var requestBodyRef *openapi3.RequestBodyRef
	if log.Request.Method != "GET" && log.Request.Method != "DELETE" {
		request := openapi3.NewRequestBody()
		var requestSchema *openapi3.Schema
		mediaType, mediaTypeParams, err := mime.ParseMediaType(contentType)
		if err != nil {
			requestSchema = openapi3.NewStringSchema()
		} else {
			switch true {
			case IsApplicationJSONMediaType(mediaType):
				respBodyJSON, _ := gojsonschema.NewBytesLoader(raw).LoadJSON()
				requestSchema, _ = BuildSchema(respBodyJSON)
			case mediaType == "application/x-www-form-urlencoded":
				requestSchema, _ = BuildForm(string(raw))
			case mediaType == "multipart/form-data":
				requestSchema, _ = BuildMultiPart(string(raw), mediaTypeParams)
			default:
				requestSchema = openapi3.NewStringSchema()
			}
		}
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

func MergeRequest(operation *openapi3.Operation, raw []byte, contentType string, log types.Log) (bool, error) {
	if log.Request.Method != "GET" && log.Request.Method != "DELETE" {
		request := operation.RequestBody.Value
		reqBodyJSON, _ := gojsonschema.NewBytesLoader(raw).LoadJSON()
		requestContent := request.Content.Get(contentType)
		schema, err := MergeSchema(reqBodyJSON, requestContent.Schema.Value)
		if err != nil {
			return false, err
		}
		requestContent.Schema.Value = schema
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
