package factories

import (
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/imdario/mergo"
	"github.com/spf13/cast"
	"github.com/xeipuuv/gojsonschema"
	"mime/multipart"
	"net/url"
	"strings"
)

func BuildMultiPart(body string, mediaTypeParams map[string]string) (*openapi3.Schema, error) {
	boundary, ok := mediaTypeParams["boundary"]
	if !ok {
		return nil, fmt.Errorf("no multipart boundary param in Content-Type")
	}

	form, err := multipart.NewReader(strings.NewReader(body), boundary).ReadForm(32 << 20)
	if err != nil {
		return nil, fmt.Errorf("failed to read form: %w", err)
	}

	schema := openapi3.NewObjectSchema()

	for key, fileHeaders := range form.File {
		fileSchema := openapi3.NewStringSchema().WithFormat("binary")
		switch len(fileHeaders) {
		case 0:
		case 1:
			schema.WithProperty(key, fileSchema)
		default:
			schema.WithProperty(key, openapi3.NewArraySchema().WithItems(fileSchema))
		}
	}

	for key, values := range form.Value {
		schema.WithProperty(key, getSchemaFromValues(values, false, ""))
	}

	return schema, nil
}

func BuildForm(body string) (*openapi3.Schema, error) {
	parseQuery, err := url.ParseQuery(body)
	if err != nil {
		return nil, err
	}
	schema := openapi3.NewObjectSchema()
	for key, values := range parseQuery {
		schema.WithProperty(key, getSchemaFromQueryValues(values))
	}
	return schema, nil
}

func IsApplicationJSONMediaType(mediaType string) bool {
	return strings.HasPrefix(mediaType, "application/") &&
		strings.HasSuffix(mediaType, "json")
}

func MergeSchema(value interface{}, schema *openapi3.Schema) (*openapi3.Schema, error) {
	merged, err := BuildSchema(value)
	if err != nil {
		return nil, err
	}
	err = mergo.Merge(schema, merged)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

func MergeSchemas(destination *openapi3.Schema, source *openapi3.Schema) error {
	err := mergo.Merge(destination, source)
	if err != nil {
		return err
	}
	return nil
}

func BuildSchema(value interface{}) (*openapi3.Schema, error) {
	var schema *openapi3.Schema
	var err error

	switch value.(type) {
	case bool:
		schema = openapi3.NewBoolSchema()
	case string:
		schema = getStringSchema(value)
	case json.Number:
		schema = getNumberSchema(value)
	case map[string]interface{}:
		schema, err = getObjectSchema(value)
		if err != nil {
			return nil, err
		}
	case []interface{}:
		schema, err = getArraySchema(value)
		if err != nil {
			return nil, err
		}
	case nil:
		schema = openapi3.NewStringSchema()
	default:
		schema = openapi3.NewStringSchema()
	}
	return schema, nil
}

func getStringSchema(value interface{}) (schema *openapi3.Schema) {
	return openapi3.NewStringSchema().WithFormat(getStringFormat(value))
}

func getNumberSchema(value interface{}) (schema *openapi3.Schema) {
	if _, err := value.(json.Number).Int64(); err != nil {
		schema = openapi3.NewFloat64Schema()
	} else {
		schema = openapi3.NewInt64Schema()
	}
	return schema
}

func getObjectSchema(value interface{}) (*openapi3.Schema, error) {
	schema := openapi3.NewObjectSchema()
	stringMapE, err := cast.ToStringMapE(value)
	if err != nil {
		return nil, err
	}
	for key, val := range stringMapE {
		if s, err := BuildSchema(val); err != nil {
			return nil, err
		} else {
			schema = schema.WithProperty(escapeString(key), s)
		}
	}
	return schema, nil
}

func escapeString(key string) string {
	if strings.Contains(key, "\"") {
		key = strings.ReplaceAll(key, "\"", "\\\"")
	}
	return key
}

func getArraySchema(value interface{}) (*openapi3.Schema, error) {
	var schema *openapi3.Schema
	sliceE, err := cast.ToSliceE(value)
	if err != nil {
		return nil, err
	}
	schemaTypeToSchema := make(map[string]*openapi3.Schema)
	for i := range sliceE {
		item, err := BuildSchema(sliceE[i])
		if err != nil {
			return nil, err
		}
		if _, ok := schemaTypeToSchema[item.Type]; !ok {
			schemaTypeToSchema[item.Type] = item
		}
	}
	switch len(schemaTypeToSchema) {
	case 0:
		schema = openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())
	case 1:
		for _, s := range schemaTypeToSchema {
			schema = openapi3.NewArraySchema().WithItems(s)
			break
		}
	default:
		var schemas []*openapi3.Schema
		for _, s := range schemaTypeToSchema {
			schemas = append(schemas, s)
		}
		schema = openapi3.NewOneOfSchema(schemas...)
	}
	return schema, nil
}

var formats = []string{
	"date",
	"time",
	"date-time",
	"email",
	"ipv4",
	"ipv6",
	"uuid",
	"json-pointer",
}

func getStringFormat(value interface{}) string {
	str, ok := value.(string)
	if !ok || str == "" {
		return ""
	}
	for _, format := range formats {
		if gojsonschema.FormatCheckers.IsFormat(format, value) {
			return format
		}
	}
	return ""
}
