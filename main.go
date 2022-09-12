package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/getkin/kin-openapi/openapi3"
	"net/url"
	"os"
	"path"
	"strconv"
)

var Version = "0.0.1"
var Priority = 1
var specs = make(map[string]*openapi3.T) // FIXME add mutex
var rawRequestBody string                // FIXME add mutex
var rawResponseBody string               // FIXME add mutex

func (conf Config) Active() bool {
	if conf.PluginActive != nil {
		return *conf.PluginActive
	}
	return false
}

func (conf Config) Response(kong *pdk.PDK) {
	logger := kong.Log
	bytes, err := kong.Request.GetRawBody()
	if err != nil {
		logger.Err(err)
		return
	}
	rawRequestBody = string(bytes)
	rawResponseBody, err = kong.ServiceResponse.GetRawBody()
	if err != nil {
		logger.Err(err)
		return
	}
}

func (conf Config) Log(kong *pdk.PDK) {
	if !conf.Active() {
		return
	}
	// get and parse log message
	data, err := kong.Log.Serialize()
	if err != nil {
		kong.Log.Err("Error getting log message: ", err.Error())
		return
	}
	var log Log
	if err := json.Unmarshal([]byte(data), &log); err != nil {
		kong.Log.Err("Error unmarshalling log message: ", err.Error())
		return
	}
	process(log, kong)
}

func process(log Log, kong *pdk.PDK) {
	logger := kong.Log
	// URL
	u, err := url.Parse(log.UpstreamURI)
	if err != nil {
		logger.Err(err)
		return
	}
	// Build specification
	if _, found := specs[log.Service.Name]; !found {
		info := &openapi3.Info{
			Title:   log.Service.Name,
			Version: "0.0.0-snapshot",
		}
		specs[log.Service.Name] = &openapi3.T{OpenAPI: "3.0.0", Info: info}
	}
	// Parameters
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
			Value: openapi3.NewPathParameter(k).WithSchema(schema),
		}
		params = append(params, &param)
	}
	// Request
	request := openapi3.NewRequestBody()
	var contentType string
	if _, found := log.Request.Headers["content-type"]; found {
		contentType = log.Request.Headers["content-type"].(string)
	} else {
		contentType = "application/json"
	}
	requestSchema := openapi3.NewSchema()
	requestContent := openapi3.Content{
		contentType: openapi3.NewMediaType().WithSchema(requestSchema),
	}
	request.Content = requestContent
	// Convert request to schema ref TODO
	// https://github.com/openclarity/speculator/blob/c8dcbd330eaf8a6551c5fd5b8fde6becdd06c6b5/pkg/spec/operation.go#L34
	// iterate from json to generate the schema
	requestRef := &openapi3.RequestBodyRef{
		Value: request,
	}
	// Response
	responses := openapi3.NewResponses()
	if _, found := log.Response.Headers["content-type"]; found {
		contentType = log.Response.Headers["content-type"].(string)
	} else {
		contentType = "application/json"
	}
	responseSchema := openapi3.NewSchema()
	responseContent := openapi3.Content{
		contentType: openapi3.NewMediaType().WithSchema(responseSchema),
	}
	response := openapi3.NewResponse()
	response.WithContent(responseContent)
	// Convert response to schema ref TODO
	// https://github.com/openclarity/speculator/blob/c8dcbd330eaf8a6551c5fd5b8fde6becdd06c6b5/pkg/spec/operation.go#L34
	// iterate from json to generate the schema
	responses[strconv.Itoa(log.Response.Status)] = &openapi3.ResponseRef{
		Value: response,
	}
	// Operation
	op := &openapi3.Operation{
		OperationID: path.Base(u.Path),
		Parameters:  params,
		RequestBody: requestRef,
		Responses:   responses,
	}
	specs[log.Service.Name].AddOperation(u.Path, log.Request.Method, op)
	// Validate
	err = specs[log.Service.Name].Validate(context.Background())
	if err != nil {
		logger.Warn(err)
	}
	// Marshal to json
	data, err := specs[log.Service.Name].MarshalJSON()
	if err != nil {
		logger.Err(err)
		return
	}
	// Write to file
	os.WriteFile(fmt.Sprintf("/logs/%s.json", log.Service.Name), prettify(data), 0644)
}

func main() {
	err := server.StartServer(New, Version, Priority)
	if err != nil {
		fmt.Println("Error starting embedded plugin server:", err.Error())
		panic(err)
	}
}
