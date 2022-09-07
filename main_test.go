package main

import (
	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/Kong/go-pdk/log"
	"testing"
)

func Test_processLog(t *testing.T) {
	tests := []struct {
		name string
		log  string
	}{
		{
			name: "Find pet by ID",
			log:  "{\"request\":{\"uri\":\"/openapi/pet/%3Clong%3E\",\"url\":\"http://localhost:8000/openapi/pet/%3Clong%3E\",\"method\":\"GET\",\"headers\":{\"host\":\"localhost:8000\",\"accept\":\"application/xml\",\"user-agent\":\"PostmanRuntime/7.29.2\",\"accept-encoding\":\"gzip, deflate, br\",\"connection\":\"keep-alive\",\"postman-token\":\"51d3ba25-ac9f-4eee-b2ea-2b3c54bcd57b\"},\"querystring\":{},\"size\":235},\"tries\":[{\"balancer_start\":1662549178041,\"ip\":\"172.20.0.2\",\"balancer_latency\":0,\"port\":4010}],\"client_ip\":\"172.20.0.1\",\"route\":{\"regex_priority\":0,\"name\":\"openapi\",\"id\":\"4208d053-d722-5d24-9535-fe1fefa26161\",\"protocols\":[\"http\",\"https\"],\"strip_path\":true,\"ws_id\":\"0dc6f45b-8f8d-40d2-a504-473544ee190b\",\"paths\":[\"/openapi\"],\"response_buffering\":true,\"service\":{\"id\":\"51a54110-f000-5535-98dd-b30846dc75b0\"},\"updated_at\":1662549171,\"preserve_host\":false,\"https_redirect_status_code\":426,\"request_buffering\":true,\"created_at\":1662549171,\"path_handling\":\"v0\"},\"workspace\":\"0dc6f45b-8f8d-40d2-a504-473544ee190b\",\"service\":{\"host\":\"prism\",\"name\":\"openapi\",\"port\":4010,\"updated_at\":1662549171,\"enabled\":true,\"id\":\"51a54110-f000-5535-98dd-b30846dc75b0\",\"write_timeout\":60000,\"retries\":5,\"created_at\":1662549171,\"path\":\"/\",\"connect_timeout\":60000,\"protocol\":\"http\",\"read_timeout\":60000,\"ws_id\":\"0dc6f45b-8f8d-40d2-a504-473544ee190b\"},\"upstream_uri\":\"/pet/%3Clong%3E\",\"response\":{\"status\":400,\"size\":460,\"headers\":{\"connection\":\"close\",\"sl-violations\":\"[{\\\"location\\\":[\\\"request\\\",\\\"path\\\",\\\"petid\\\"],\\\"severity\\\":\\\"Error\\\",\\\"code\\\":\\\"type\\\",\\\"message\\\":\\\"must be integer\\\"}]\",\"via\":\"kong/2.8.1.4-enterprise-edition\",\"x-kong-proxy-latency\":\"103\",\"content-length\":\"0\",\"x-kong-upstream-latency\":\"90\",\"date\":\"Wed, 07 Sep 2022 11:12:58 GMT\",\"access-control-expose-headers\":\"*\",\"access-control-allow-origin\":\"*\",\"access-control-allow-headers\":\"*\",\"access-control-allow-credentials\":\"true\"}},\"started_at\":1662549177938,\"latencies\":{\"proxy\":90,\"kong\":103,\"request\":193}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processLog(tt.log, log.Log{bridge.New(bridgetest.Mock(t, nil))})
		})
	}
}
