package main

import (
	"encoding/json"
	"fmt"
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/log"
	"github.com/Kong/go-pdk/server"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/factories"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/types"
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/utils"
	"github.com/patrickmn/go-cache"
	"net/url"
	"time"
)

var Version = "0.0.1"
var Priority = 1
var requests = cache.New(5*time.Minute, 10*time.Minute)
var responses = cache.New(5*time.Minute, 10*time.Minute)
var specs = make(map[string]*openapi3.T)                         // FIXME add mutex
var operations = make(map[string]map[string]*openapi3.Operation) // FIXME add mutex
var methods = make(map[string]map[string]string)                 // FIXME add mutex
var registeredPaths = make(map[string]map[string]string)         // FIXME add mutex

type Config struct {
	PluginActive *bool `json:"active"`
}

func New() interface{} {
	utils.Init()
	return &Config{}
}

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
	//os.WriteFile("/logs/request.log", request, 0644) TODO remove
	//os.WriteFile("/logs/response.log", []byte(response), 0644) TODO remove
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
	//os.WriteFile("/logs/log.log", []byte(rawLog), 0644) TODO remove
}

func process(rawLog *string, rawRequest *[]byte, rawResponse *string, logger log.Log) {
	var log types.Log
	var updated = false
	// build
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
	// content Type
	var contentType string
	if _, found := log.Request.Headers["content-type"]; found {
		contentType = log.Request.Headers["content-type"].(string)
	} else {
		contentType = "application/json"
	}
	// search specification
	if _, found := specs[log.Service.Name]; !found {
		operations[log.Service.Name] = make(map[string]*openapi3.Operation)
		methods[log.Service.Name] = make(map[string]string)
		registeredPaths[log.Service.Name] = make(map[string]string)
		var read *openapi3.T
		read, err = utils.Read(log.Service.Name)
		if err != nil {
			logger.Err(err)
			return
		}
		if read == nil {
			specs[log.Service.Name] = factories.BuildSpecification(log.Service.Name, "3.0.0")
		} else {
			specs[log.Service.Name] = read
		}
	}
	var name string
	// aggregate
	lookup := factories.AggregateSpecification(specs[log.Service.Name],
		registeredPaths[log.Service.Name],
		methods[log.Service.Name],
		operations[log.Service.Name])
	// match
	matched, route, _ := match(log.Request.Method, u.Path, contentType, lookup)
	if !matched {
		// url
		url := factories.CreateParameterizedPath(u.Path)
		// parameters
		params := factories.BuildParams(url, u.Path, log, logger)
		// request
		operationRequest := factories.BuildRequest(*rawRequest, contentType, log)
		// response
		operationResponse := factories.BuildResponses(*rawResponse, log)
		// operation
		name = utils.GetName(log.Request.Method, url)
		operation := &openapi3.Operation{
			OperationID: name,
			Parameters:  params,
			RequestBody: operationRequest,
			Responses:   operationResponse,
		}
		operations[log.Service.Name][name] = operation
		registeredPaths[log.Service.Name][name] = url
		methods[log.Service.Name][name] = log.Request.Method
		updated = true
	} else {
		// merge
		name = utils.GetName(log.Request.Method, route)
		loaded := operations[log.Service.Name][name]
		if factories.MergeParams(loaded, registeredPaths[log.Service.Name][name], u.Path, log, logger) ||
			factories.MergeRequest(loaded, *rawRequest, contentType, log) ||
			factories.MergeResponses(loaded, *rawResponse, log) {
			updated = true
		}
	}
	if updated {
		specification := factories.AggregateSpecification(specs[log.Service.Name],
			registeredPaths[log.Service.Name],
			methods[log.Service.Name],
			operations[log.Service.Name])
		err = utils.Write(log.Service.Name, specification)
		if err != nil {
			logger.Err(err)
			return
		}
	}
}

func main() {
	err := server.StartServer(New, Version, Priority)
	if err != nil {
		fmt.Println("Error starting embedded plugin server:", err.Error())
		panic(err)
	}
}
