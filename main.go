package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/log"
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

type Config struct {
	PluginActive *bool `json:"active"`
}

func (conf Config) Active() bool {
	if conf.PluginActive != nil {
		return *conf.PluginActive
	}
	return false
}

func New() interface{} {
	return &Config{}
}

type LogMsg struct {
	Latencies struct {
		Request int `json:"request"`
		Kong    int `json:"kong"`
		Proxy   int `json:"proxy"`
	} `json:"latencies"`
	Service struct {
		Host           string `json:"host"`
		CreatedAt      int    `json:"created_at"`
		ConnectTimeout int    `json:"connect_timeout"`
		ID             string `json:"id"`
		Name           string `json:"name"`
		Protocol       string `json:"protocol"`
		ReadTimeout    int    `json:"read_timeout"`
		Port           int    `json:"port"`
		Path           string `json:"path"`
		UpdatedAt      int    `json:"updated_at"`
		WriteTimeout   int    `json:"write_timeout"`
		Retries        int    `json:"retries"`
		WsID           string `json:"ws_id"`
	} `json:"service"`
	Request struct {
		Querystring map[string]interface{} `json:"querystring"`
		Size        int                    `json:"size"`
		URI         string                 `json:"uri"`
		URL         string                 `json:"url"`
		Headers     map[string]interface{} `json:"headers"`
		Method      string                 `json:"method"`
	} `json:"request"`
	Tries []struct {
		BalancerLatency int    `json:"balancer_latency"`
		Port            int    `json:"port"`
		BalancerStart   int64  `json:"balancer_start"`
		IP              string `json:"ip"`
	} `json:"tries"`
	ClientIP    string `json:"client_ip"`
	Workspace   string `json:"workspace"`
	UpstreamURI string `json:"upstream_uri"`
	Response    struct {
		Headers map[string]interface{} `json:"headers"`
		Status  int                    `json:"status"`
		Size    int                    `json:"size"`
	} `json:"response"`
	Route struct {
		ID                      string   `json:"id"`
		Paths                   []string `json:"paths"`
		Protocols               []string `json:"protocols"`
		StripPath               bool     `json:"strip_path"`
		CreatedAt               int      `json:"created_at"`
		WsID                    string   `json:"ws_id"`
		RequestBuffering        bool     `json:"request_buffering"`
		UpdatedAt               int      `json:"updated_at"`
		PreserveHost            bool     `json:"preserve_host"`
		RegexPriority           int      `json:"regex_priority"`
		ResponseBuffering       bool     `json:"response_buffering"`
		HTTPSRedirectStatusCode int      `json:"https_redirect_status_code"`
		PathHandling            string   `json:"path_handling"`
		Service                 struct {
			ID string `json:"id"`
		} `json:"service"`
	} `json:"route"`
	StartedAt int64 `json:"started_at"`
}

func (conf Config) Log(kong *pdk.PDK) {
	if !conf.Active() {
		return
	}
	// get and parse log message
	s, err := kong.Log.Serialize()
	if err != nil {
		kong.Log.Err("Error getting log message: ", err.Error())
		return
	}
	processLog(s, kong.Log)
}

func processLog(s string, logger log.Log) {
	var msg LogMsg
	if err := json.Unmarshal([]byte(s), &msg); err != nil {
		_ = logger.Err("Error unmarshalling log message: ", err.Error())
		return
	}
	u, err := url.Parse(msg.UpstreamURI)
	if err != nil {
		logger.Err(err)
		return
	}
	if _, found := specs[msg.Service.Name]; !found {
		info := &openapi3.Info{
			Title:   msg.Service.Name,
			Version: "0.1",
		}
		specs[msg.Service.Name] = &openapi3.T{OpenAPI: "3.0.0", Info: info}
	}
	// Parameters
	var params []*openapi3.ParameterRef
	for k, v := range msg.Request.Querystring {
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
	// Responses
	responses := openapi3.NewResponses()
	content := openapi3.Content{}
	if _, found := msg.Response.Headers["content-type"]; found {
		content[msg.Response.Headers["content-type"].(string)] = openapi3.NewMediaType()
	}
	responses[strconv.Itoa(msg.Response.Status)] = &openapi3.ResponseRef{
		Value: openapi3.NewResponse().WithContent(content),
	}
	// Operation
	op := &openapi3.Operation{
		OperationID: path.Base(u.Path),
		Parameters:  params,
		Responses:   responses,
	}
	specs[msg.Service.Name].AddOperation(u.Path, msg.Request.Method, op)
	// Validate
	err = specs[msg.Service.Name].Validate(context.Background())
	if err != nil {
		logger.Warn(err)
	}
	// marshal to json
	data, err := specs[msg.Service.Name].MarshalJSON()
	if err != nil {
		logger.Err(err)
		return
	}
	// write to file
	os.WriteFile(fmt.Sprintf("/tmp/%s.json", msg.Service.Name), prettify(data), 0644)
}

func prettify(data []byte) []byte {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return []byte("")
	}
	return prettyJSON.Bytes()
}

func main() {
	err := server.StartServer(New, Version, Priority)
	if err != nil {
		fmt.Println("Error starting embedded plugin server:", err.Error())
		panic(err)
	}
}
