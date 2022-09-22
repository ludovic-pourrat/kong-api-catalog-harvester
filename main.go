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
	"github.com/ludovic-pourrat/kong-api-catalog-harvester/utils/pathtrie"
	"github.com/patrickmn/go-cache"
	"net/url"
	"strings"
	"time"
)

var Version = "0.0.1"
var Priority = 1
var requests = cache.New(5*time.Minute, 10*time.Minute)
var responses = cache.New(5*time.Minute, 10*time.Minute)
var specs = make(map[string]*openapi3.T) // FIXME add mutex
var registeredPaths = pathtrie.New()     // FIXME add mutex

type Config struct {
	PluginActive *bool `json:"active"`
}

func New() interface{} {
	_ = utils.Init()
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
		_ = logger.Err("Error getting unique request id: ", err.Error())
		return
	}
	request, err := kong.Request.GetRawBody()
	if err != nil {
		_ = logger.Err(err)
		return
	}
	requests.SetDefault(id, &request)
	response, err := kong.ServiceResponse.GetRawBody()
	if err != nil {
		_ = logger.Err(err)
		return
	}
	responses.SetDefault(id, &response)

}

func (conf Config) Log(kong *pdk.PDK) {
	if !conf.Active() {
		return
	}
	logger := kong.Log
	// get and parse log message
	rawLog, err := kong.Log.Serialize()
	if err != nil {
		_ = logger.Err("Error getting log message: ", err.Error())
		return
	}
	id, err := kong.Request.GetHeader("Kong-Request-ID")
	if err != nil {
		_ = logger.Err("Error getting unique request id: ", err.Error())
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
	var captured types.Log
	var updated = false
	// build
	if err := json.Unmarshal([]byte(*rawLog), &captured); err != nil {
		_ = logger.Err("Error unmarshalling log message: ", err.Error())
		return
	}
	// URL
	u, err := url.Parse(captured.UpstreamURI)
	if err != nil {
		_ = logger.Err(err)
		return
	}
	// content Type
	var contentType string
	if _, found := captured.Request.Headers["content-type"]; found {
		contentType = captured.Request.Headers["content-type"].(string)
	} else {
		contentType = "application/json"
	}
	// search specification
	if _, found := specs[captured.Service.Name]; !found {
		specs[captured.Service.Name] = factories.BuildSpecification(captured.Service.Name, "3.0.0")
	} else {
		specs[captured.Service.Name] = factories.CloneSpecification(specs[captured.Service.Name],
			registeredPaths)
	}
	var name string
	// match
	matched, route, err := match(captured.Request.Method, u.Path, contentType, specs[captured.Service.Name])
	if !matched {
		var computed string
		// computed
		if strings.EqualFold(err.Error(), "method not allowed") {
			computed = factories.CreateParameterizedPath(route)
		} else {
			computed = factories.CreateParameterizedPath(u.Path)
		}
		// parameters
		params := factories.BuildParams(computed, u.Path, captured, logger)
		// request
		operationRequest := factories.BuildRequest(*rawRequest, contentType, captured)
		// response
		operationResponse := factories.BuildResponses(*rawResponse, captured)
		// operation
		name = utils.GetName(captured.Request.Method, strings.Split(computed, "/"))
		operation := &openapi3.Operation{
			OperationID: name,
			Parameters:  params,
			RequestBody: operationRequest,
			Responses:   operationResponse,
		}
		registeredPaths.Insert(computed, u.Path, operation, captured.Request.Method, 1)
		updated = true
	} else {
		var updatedRequest, updatedResponses bool
		// merge
		name = utils.GetName(captured.Request.Method, strings.Split(route, "/"))
		for _, path := range registeredPaths.Nodes() {
			for _, operation := range path.Operations {
				if operation.OperationID == name {
					updatedParams := factories.MergeParams(operation, route, u.Path, captured, logger)
					updatedRequest, err = factories.MergeRequest(operation, *rawRequest, contentType, captured)
					if err != nil {
						_ = logger.Err(err)
						return
					}
					updatedResponses, err = factories.MergeResponses(operation, *rawResponse, captured)
					if err != nil {
						_ = logger.Err(err)
						return
					}
					if updatedParams || updatedRequest || updatedResponses {
						updated = true
					}
				}
			}
		}
	}
	if updated {
		specs[captured.Service.Name] = factories.CloneSpecification(specs[captured.Service.Name],
			registeredPaths)
		err = utils.Write(captured.Service.Name, specs[captured.Service.Name])
		if err != nil {
			_ = logger.Err(err)
		}
	}
	return
}

func main() {
	err := server.StartServer(New, Version, Priority)
	if err != nil {
		fmt.Println("Error starting embedded plugin server:", err.Error())
		panic(err)
	}
}
