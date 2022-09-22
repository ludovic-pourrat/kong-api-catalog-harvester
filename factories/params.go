package factories

import (
	"fmt"
	"github.com/Kong/go-pdk/log"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/types"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/utils"
	"strings"
)

func BuildParams(url string, path string, log types.Log, logger log.Log) []*openapi3.ParameterRef {
	var params []*openapi3.ParameterRef
	for k, v := range log.Request.Querystring {
		var schema *openapi3.Schema
		switch v.(type) {
		case string:
			schema = openapi3.NewStringSchema()
		case []interface{}:
			schema = openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())
		default:
			_ = logger.Err("unknown type for querystring ", fmt.Sprintf("%T", v))
		}
		param := openapi3.ParameterRef{
			Value: openapi3.NewQueryParameter(k).WithSchema(schema),
		}
		params = append(params, &param)
	}

	parts := strings.Split(url, "/")
	values := strings.Split(path, "/")

	for index, part := range parts {
		if utils.IsPathParam(part) {
			part = strings.TrimPrefix(part, "{")
			part = strings.TrimSuffix(part, "}")
			param := openapi3.ParameterRef{
				Value: openapi3.NewPathParameter(part).WithSchema(getParamSchema(values[index])),
			}
			params = append(params, &param)
		}
	}

	return params
}

func MergeParams(operation *openapi3.Operation, url string, path string, log types.Log, logger log.Log) bool {
	var params []*openapi3.ParameterRef
	var updated = false
	params = BuildParams(url, path, log, logger)
	for _, param := range params {
		found := false
		for _, op := range operation.Parameters {
			if op.Value.Name == param.Value.Name {
				found = true
				break
			}
		}
		if !found {
			operation.Parameters = append(operation.Parameters, param)
			updated = true
		}
	}
	return updated
}
