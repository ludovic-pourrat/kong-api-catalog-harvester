package factories

import (
	"strconv"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

func addQueryParam(operation *openapi3.Operation, key string, values []string) *openapi3.Operation {
	operation.AddParameter(openapi3.NewQueryParameter(key).WithSchema(getSchemaFromQueryValues(values)))
	return operation
}

func getSchemaFromQueryValues(values []string) *openapi3.Schema {
	var schema *openapi3.Schema
	if len(values) == 0 || values[0] == "" {
		schema = openapi3.NewBoolSchema()
		schema.AllowEmptyValue = true
	} else {
		schema = getSchemaFromValues(values, true, openapi3.ParameterInQuery)
	}
	return schema
}

func getSchemaFromValues(values []string, shouldTryArraySchema bool, paramInType string) *openapi3.Schema {
	valuesLen := len(values)

	if valuesLen == 0 {
		return nil
	}

	if valuesLen == 1 {
		return getSchemaFromValue(values[0], shouldTryArraySchema, paramInType)
	}

	// find the most common schema for the items type
	return openapi3.NewArraySchema().WithItems(getCommonSchema(values, paramInType))
}

func getSchemaFromValue(value string, shouldTryArraySchema bool, paramInType string) *openapi3.Schema {
	if isDateFormat(value) {
		return openapi3.NewStringSchema()
	}

	if shouldTryArraySchema {
		schema, _ := getNewArraySchema(value, paramInType)
		if schema != nil {
			return schema
		}
	}

	// nolint:gomnd
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return openapi3.NewInt64Schema()
	}

	// nolint:gomnd
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return openapi3.NewFloat64Schema()
	}

	// TODO: not sure that `strconv.ParseBool` will do the job, it depends what is considers as boolean string
	// The Go implementation for example uses `strconv.FormatBool(value)` ==> true/false
	// But if we look at swag.ConvertBool - `checked` is evaluated as true so `unchecked` will be false?
	// Also when using `strconv.ParseBool` 1 is considered as true so we must check for int before running it
	if _, err := strconv.ParseBool(value); err == nil {
		return openapi3.NewBoolSchema()
	}

	return openapi3.NewStringSchema().WithFormat(getStringFormat(value))
}

func getNewArraySchema(value string, paramInType string) (schema *openapi3.Schema, style string) {
	var supportedSerializationStyles []string

	switch paramInType {
	case openapi3.ParameterInHeader:
		supportedSerializationStyles = []string{
			openapi3.SerializationSimple,
		}
	case openapi3.ParameterInQuery:
		supportedSerializationStyles = []string{
			openapi3.SerializationForm,
			openapi3.SerializationSpaceDelimited,
			openapi3.SerializationPipeDelimited,
		}
	case openapi3.ParameterInCookie:
		supportedSerializationStyles = []string{
			openapi3.SerializationForm,
		}
	default:
		return nil, ""
	}

	for _, style = range supportedSerializationStyles {
		byStyle := splitByStyle(value, style)
		// Will create an array only if more than a single element exists
		if len(byStyle) > 1 {
			return getSchemaFromValues(byStyle, false, paramInType), style
		}
	}

	return nil, ""
}

func getCommonSchema(values []string, paramInType string) *openapi3.Schema {
	var schemaType string
	var schema *openapi3.Schema

	for _, value := range values {
		schema = getSchemaFromValue(value, false, paramInType)
		if schemaType == "" {
			// first value, save schema type
			schemaType = schema.Type
		} else if schemaType != schema.Type {
			// different schema type found, defaults to string schema
			return openapi3.NewStringSchema()
		}
	}

	// identical schema type found
	return schema
}

func splitByStyle(data, style string) []string {
	if data == "" {
		return nil
	}
	var sep string
	switch style {
	case openapi3.SerializationForm, openapi3.SerializationSimple:
		sep = ","
	case openapi3.SerializationSpaceDelimited:
		sep = " "
	case openapi3.SerializationPipeDelimited:
		sep = "|"
	default:
		return nil
	}
	var result []string
	for _, s := range strings.Split(data, sep) {
		if ts := strings.TrimSpace(s); ts != "" {
			result = append(result, ts)
		}
	}
	return result
}

func isDateFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	if _, err := time.Parse(time.ANSIC, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.UnixDate, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.RubyDate, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.RFC822, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.RFC822Z, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.RFC850, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.RFC1123, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.RFC1123Z, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.Stamp, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.StampMilli, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.StampMicro, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.StampNano, asString); err == nil {
		return true
	}
	return false
}
