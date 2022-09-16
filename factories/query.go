package factories

import (
	"strconv"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

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
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return openapi3.NewInt64Schema()
	}
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return openapi3.NewFloat64Schema()
	}
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
			schemaType = schema.Type
		} else if schemaType != schema.Type {
			return openapi3.NewStringSchema()
		}
	}
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
