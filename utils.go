package main

import (
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cast"
	"github.com/xeipuuv/gojsonschema"
	"strings"
)

func New() interface{} {
	return &Config{}
}

func getSchema(value interface{}) (schema *openapi3.Schema, err error) {
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
		// TODO: Not sure how to handle null. ex: {"size":3,"err":null}
		schema = openapi3.NewStringSchema()
	default:
		// TODO:
		// I've tested additionalProperties and it seems like properties - we will might have problems in the diff logic
		// openapi3.MapProperty()
		// openapi3.RefProperty()
		// openapi3.RefSchema()
		// openapi3.ComposedSchema() - discriminator?
		return nil, fmt.Errorf("unexpected value type. value=%v, type=%T", value, value)
	}
	return schema, nil
}

func getStringSchema(value interface{}) (schema *openapi3.Schema) {
	return openapi3.NewStringSchema().WithFormat(getStringFormat(value))
}

func getNumberSchema(value interface{}) (schema *openapi3.Schema) {
	// https://swagger.io/docs/specification/data-models/data-types/#numbers
	// It is important to try first convert it to int
	if _, err := value.(json.Number).Int64(); err != nil {
		// if failed to convert to int it's a double
		// TODO: we will set a 'double' and not a 'float' - is that ok?
		schema = openapi3.NewFloat64Schema()
	} else {
		schema = openapi3.NewInt64Schema()
	}
	// TODO: Format
	// openapi3.Int8Property()
	// openapi3.Int16Property()
	// openapi3.Int32Property()
	// openapi3.Float64Property()
	// openapi3.Float32Property()
	return schema /*.WithExample(value)*/
}

func getObjectSchema(value interface{}) (schema *openapi3.Schema, err error) {
	schema = openapi3.NewObjectSchema()
	stringMapE, err := cast.ToStringMapE(value)
	if err != nil {
		return nil, fmt.Errorf("failed to cast to string map. value=%v: %w", value, err)
	}
	for key, val := range stringMapE {
		if s, err := getSchema(val); err != nil {
			return nil, fmt.Errorf("failed to get schema from string map. key=%v, value=%v: %w", key, val, err)
		} else {
			schema = schema.WithProperty(escapeString(key), s)
		}
	}
	return schema, nil
}

func escapeString(key string) string {
	// need to escape double quotes if exists
	if strings.Contains(key, "\"") {
		key = strings.ReplaceAll(key, "\"", "\\\"")
	}
	return key
}

func getArraySchema(value interface{}) (schema *openapi3.Schema, err error) {
	sliceE, err := cast.ToSliceE(value)
	if err != nil {
		return nil, fmt.Errorf("failed to cast to slice. value=%v: %w", value, err)
	}
	// in order to support mixed type array we will map all schemas by schema type
	schemaTypeToSchema := make(map[string]*openapi3.Schema)
	for i := range sliceE {
		item, err := getSchema(sliceE[i])
		if err != nil {
			return nil, fmt.Errorf("failed to get items schema from slice. value=%v: %w", sliceE[i], err)
		}
		if _, ok := schemaTypeToSchema[item.Type]; !ok {
			schemaTypeToSchema[item.Type] = item
		}
	}
	switch len(schemaTypeToSchema) {
	case 0:
		// array is empty, but we can't create an empty array property (Schemas with 'type: array', require a sibling 'items:' field)
		// we will create string type items as a default value
		schema = openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())
	case 1:
		for _, s := range schemaTypeToSchema {
			schema = openapi3.NewArraySchema().WithItems(s)
			break
		}
	default:
		// oneOf
		// https://swagger.io/docs/specification/data-models/oneof-anyof-allof-not/
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
	// "relative-json-pointer", // matched with "1.147.1"
	// "hostname",
	// "regex",
	// "uri",           // can be also iri
	// "uri-reference", // can be also iri-reference
	// "uri-template",
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
