package factories

import (
	"fmt"
	"github.com/Kong/go-pdk/log"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/types"
)

func BuildParams(log types.Log, logger log.Log) []*openapi3.ParameterRef {
	var params []*openapi3.ParameterRef
	for k, v := range log.Request.Querystring {
		var schema *openapi3.Schema
		switch v.(type) {
		case string:
			schema = openapi3.NewStringSchema()
		case []interface{}:
			schema = openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())
		default:
			logger.Err("unknown type for querystring ", fmt.Sprintf("%T", v))
		}
		param := openapi3.ParameterRef{
			Value: openapi3.NewQueryParameter(k).WithSchema(schema),
		}
		params = append(params, &param)
	}
	return params
}

func MergeParams(operation *openapi3.Operation, log types.Log, logger log.Log) bool {
	var params []*openapi3.ParameterRef
	params = BuildParams(log, logger)
	operation.Parameters = append(operation.Parameters, params...)
	return true
}
