package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/log"
	"github.com/Kong/go-pdk/server"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/invopop/yaml"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/factories"
	"github.com/patrickmn/go-cache"
	"github.com/xeipuuv/gojsonschema"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

var Version = "0.0.1"
var Priority = 1
var requests = cache.New(5*time.Minute, 10*time.Minute)
var responses = cache.New(5*time.Minute, 10*time.Minute)
var specs = make(map[string]*openapi3.T) // FIXME add mutex

func (conf Config) Active() bool {
	if conf.PluginActive != nil {
		return *conf.PluginActive
	}
	return false
}

func (conf Config) Response(kong *pdk.PDK) {
	logger := kong.Log
	id, err := kong.Request.GetHeader("Kong-Request-ID")
	if err != nil {
		kong.Log.Err("Error getting unique request id: ", err.Error())
		return
	}
	request, err := kong.Request.GetRawBody()
	if err != nil {
		logger.Err(err)
		return
	}
	requests.SetDefault(id, &request)
	response, err := kong.ServiceResponse.GetRawBody()
	if err != nil {
		logger.Err(err)
		return
	}
	responses.SetDefault(id, &response)
	//os.WriteFile("/logs/request", request, 0644)
	//os.WriteFile("/logs/response", []byte(response), 0644)
}

func (conf Config) Log(kong *pdk.PDK) {
	if !conf.Active() {
		return
	}
	// get and parse log message
	rawLog, err := kong.Log.Serialize()
	if err != nil {
		kong.Log.Err("Error getting log message: ", err.Error())
		return
	}
	id, err := kong.Request.GetHeader("Kong-Request-ID")
	if err != nil {
		kong.Log.Err("Error getting unique request id: ", err.Error())
		return
	}
	// retrieve request and response from cache
	rawRequest, _ := requests.Get(id)
	defer requests.Delete(id)
	rawResponse, _ := responses.Get(id)
	defer responses.Delete(id)
	process(&rawLog, rawRequest.(*[]byte), rawResponse.(*string), kong.Log)
	//os.WriteFile("/logs/log", []byte(rawLog), 0644)
}

func process(rawLog *string, rawRequest *[]byte, rawResponse *string, logger log.Log) {
	var log Log
	if err := json.Unmarshal([]byte(*rawLog), &log); err != nil {
		logger.Err("Error unmarshalling log message: ", err.Error())
		return
	}
	// URL
	u, err := url.Parse(log.UpstreamURI)
	if err != nil {
		logger.Err(err)
		return
	}
	// Build specification
	if _, found := specs[log.Service.Name]; !found {
		specs[log.Service.Name] = factories.BuildSpecification(log.Service.Name, "3.0.0")
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
			Value: openapi3.NewQueryParameter(k).WithSchema(schema),
		}
		params = append(params, &param)
	}
	var contentType string
	var requestRef *openapi3.RequestBodyRef
	// Request
	if log.Request.Method != "GET" && log.Request.Method != "DELETE" {
		request := openapi3.NewRequestBody()
		if _, found := log.Request.Headers["content-type"]; found {
			contentType = log.Request.Headers["content-type"].(string)
		} else {
			contentType = "application/json"
		}
		reqBodyJSON, _ := gojsonschema.NewBytesLoader(*rawRequest).LoadJSON()
		requestSchema, _ := getSchema(reqBodyJSON)
		requestContent := openapi3.Content{
			contentType: openapi3.NewMediaType().WithSchema(requestSchema),
		}
		request.Content = requestContent
		request.Description = ""
		requestRef = &openapi3.RequestBodyRef{
			Value: request,
		}
	} else {
		requestRef = nil
	}
	// Response
	responses := openapi3.NewResponses()
	if _, found := log.Response.Headers["content-type"]; found {
		contentType = log.Response.Headers["content-type"].(string)
	} else {
		contentType = "application/json"
	}
	respBodyJSON, _ := gojsonschema.NewStringLoader(*rawResponse).LoadJSON()
	responseSchema, _ := getSchema(respBodyJSON)
	responseContent := openapi3.Content{
		contentType: openapi3.NewMediaType().WithSchema(responseSchema),
	}
	response := openapi3.NewResponse()
	response.WithContent(responseContent)
	response.WithDescription("")
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
	// Marshal to yaml
	data, err := yaml.Marshal(specs[log.Service.Name])
	if err != nil {
		logger.Err(err)
		return
	}
	// Write to file
	os.WriteFile(fmt.Sprintf("/logs/%s.yaml", log.Service.Name), data, 0644)
}

func main() {
	err := server.StartServer(New, Version, Priority)
	if err != nil {
		fmt.Println("Error starting embedded plugin server:", err.Error())
		panic(err)
	}
}
