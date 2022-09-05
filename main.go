package main

import (
	"encoding/json"
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"log"
	"strconv"
)

var Version = "0.0.1"
var Priority = 1

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
		Querystring struct {
		} `json:"querystring"`
		Size    int                    `json:"size"`
		URI     string                 `json:"uri"`
		URL     string                 `json:"url"`
		Headers map[string]interface{} `json:"headers"`
		Method  string                 `json:"method"`
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
		_ = kong.Log.Err("Error getting log message: ", err.Error())
		return
	}
	var msg LogMsg
	if err := json.Unmarshal([]byte(s), &msg); err != nil {
		_ = kong.Log.Err("Error unmarshalling log message: ", err.Error())
		return
	}
	log.Printf(strconv.FormatInt(msg.StartedAt, 10))
}

func main() {
	err := server.StartServer(New, Version, Priority)
	if err != nil {
		log.Printf("Error starting embedded plugin server: %s", err.Error())
		panic(err)
	}
}
