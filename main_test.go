package main

import (
	"fmt"
	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/Kong/go-pdk/log"
	"testing"
)

func Test_process(t *testing.T) {
	logPattern := "{\"request\":{\"uri\":\"/openapi/%[1]s\",\"url\":\"http://localhost:8000/openapi/%[1]s\",\"method\":\"GET\",\"headers\":{\"host\":\"localhost:8000\",\"accept\":\"application/xml\",\"user-agent\":\"PostmanRuntime/7.29.2\",\"accept-encoding\":\"gzip, deflate, br\",\"connection\":\"keep-alive\",\"postman-token\":\"51d3ba25-ac9f-4eee-b2ea-2b3c54bcd57b\"},\"querystring\":{},\"size\":235},\"tries\":[{\"balancer_start\":1662549178041,\"ip\":\"172.20.0.2\",\"balancer_latency\":0,\"port\":4010}],\"client_ip\":\"172.20.0.1\",\"route\":{\"regex_priority\":0,\"name\":\"openapi\",\"id\":\"4208d053-d722-5d24-9535-fe1fefa26161\",\"protocols\":[\"http\",\"https\"],\"strip_path\":true,\"ws_id\":\"0dc6f45b-8f8d-40d2-a504-473544ee190b\",\"paths\":[\"/openapi\"],\"response_buffering\":true,\"service\":{\"id\":\"51a54110-f000-5535-98dd-b30846dc75b0\"},\"updated_at\":1662549171,\"preserve_host\":false,\"https_redirect_status_code\":426,\"request_buffering\":true,\"created_at\":1662549171,\"path_handling\":\"v0\"},\"workspace\":\"0dc6f45b-8f8d-40d2-a504-473544ee190b\",\"service\":{\"host\":\"prism\",\"name\":\"openapi\",\"port\":4010,\"updated_at\":1662549171,\"enabled\":true,\"id\":\"51a54110-f000-5535-98dd-b30846dc75b0\",\"write_timeout\":60000,\"retries\":5,\"created_at\":1662549171,\"path\":\"/\",\"connect_timeout\":60000,\"protocol\":\"http\",\"read_timeout\":60000,\"ws_id\":\"0dc6f45b-8f8d-40d2-a504-473544ee190b\"},\"upstream_uri\":\"/%[1]s\",\"response\":{\"status\":400,\"size\":460,\"headers\":{\"connection\":\"close\",\"sl-violations\":\"[{\\\"location\\\":[\\\"request\\\",\\\"path\\\",\\\"petid\\\"],\\\"severity\\\":\\\"Error\\\",\\\"code\\\":\\\"type\\\",\\\"message\\\":\\\"must be integer\\\"}]\",\"via\":\"kong/2.8.1.4-enterprise-edition\",\"x-kong-proxy-latency\":\"103\",\"content-length\":\"0\",\"x-kong-upstream-latency\":\"90\",\"date\":\"Wed, 07 Sep 2022 11:12:58 GMT\",\"access-control-expose-headers\":\"*\",\"access-control-allow-origin\":\"*\",\"access-control-allow-headers\":\"*\",\"access-control-allow-credentials\":\"true\"}},\"started_at\":1662549177938,\"latencies\":{\"proxy\":90,\"kong\":103,\"request\":193}}"
	tests := []struct {
		name     string
		log      string
		request  string
		response string
	}{
		{
			name: "Find pet by ID",
			log:  fmt.Sprintf(logPattern, "pet/123/store/yellow"),
		},
		{
			name: "Find pet by ID",
			log:  fmt.Sprintf(logPattern, "pet/123/store/blue"),
		},
		{
			name: "Find pet by ID",
			log:  fmt.Sprintf(logPattern, "pet/123/store/red"),
		},
		{
			name: "Find pet by ID",
			log:  fmt.Sprintf(logPattern, "pet/123/store/brown"),
		},
		{
			name: "Find pet by ID",
			log:  fmt.Sprintf(logPattern, "pet/123/store/yyyy"),
		},
		{
			name: "Find pet by ID",
			log:  fmt.Sprintf(logPattern, "pet/123/store/dsdsdsd"),
		},
		{
			name: "Find pet by ID",
			log:  fmt.Sprintf(logPattern, "pet/tutu/store/green"),
		},
		{
			name:     "Add a new pet to the store",
			log:      "{\"route\":{\"ws_id\":\"0dc6f45b-8f8d-40d2-a504-473544ee190b\",\"request_buffering\":true,\"response_buffering\":true,\"strip_path\":true,\"protocols\":[\"http\",\"https\"],\"service\":{\"id\":\"51a54110-f000-5535-98dd-b30846dc75b0\"},\"preserve_host\":false,\"https_redirect_status_code\":426,\"path_handling\":\"v0\",\"paths\":[\"/openapi\"],\"id\":\"4208d053-d722-5d24-9535-fe1fefa26161\",\"created_at\":1662970647,\"name\":\"openapi\",\"regex_priority\":0,\"updated_at\":1662970647},\"response\":{\"headers\":{\"x-kong-response-latency\":\"171\",\"connection\":\"close\",\"server\":\"kong/2.8.1.4-enterprise-edition\",\"content-length\":\"36\",\"content-type\":\"application/json; charset=utf-8\"},\"status\":503,\"size\":283},\"request\":{\"method\":\"POST\",\"uri\":\"/openapi/pet\",\"querystring\":{},\"url\":\"http://localhost:8000/openapi/pet\",\"size\":279,\"headers\":{\"accept-encoding\":\"gzip, deflate, br\",\"accept\":\"application/json\",\"content-length\":\"295\",\"content-type\":\"application/json\",\"connection\":\"keep-alive\",\"kong-request-id\":\"52bef2d4-951a-4e6d-ae89-90be3105b062\",\"host\":\"localhost:8000\",\"user-agent\":\"PostmanRuntime/7.29.2\",\"postman-token\":\"f0d4d15f-1d9c-4f13-8fc3-2ae37466a909\"}},\"latencies\":{\"request\":172,\"kong\":171,\"proxy\":-1},\"tries\":[],\"upstream_uri\":\"/pet\",\"service\":{\"read_timeout\":60000,\"ws_id\":\"0dc6f45b-8f8d-40d2-a504-473544ee190b\",\"write_timeout\":60000,\"path\":\"/\",\"id\":\"51a54110-f000-5535-98dd-b30846dc75b0\",\"created_at\":1662970647,\"host\":\"prism\",\"protocol\":\"http\",\"name\":\"openapi\",\"enabled\":true,\"port\":4010,\"updated_at\":1662970647,\"connect_timeout\":60000,\"retries\":5},\"client_ip\":\"172.21.0.1\",\"started_at\":1662970651100,\"workspace\":\"0dc6f45b-8f8d-40d2-a504-473544ee190b\"}",
			request:  "{\"id\":850,\"name\":\"firewall\",\"category\":{\"id\":442,\"name\":\"Dogs\"},\"photoUrls\":[\"http:\\/\\/placeimg.com\\/640\\/480\"],\"tags\":[{\"id\":633,\"name\":\"Vision-oriented bi-directional groupware\"}]}",
			response: "{\"category\":{\"id\":-1902395453419917300,\"name\":\"laborum est qui\"},\"id\":5249524181381050000,\"name\":\"sint consectetur occaecat do\",\"photoUrls\":[\"qui\",\"velit deserunt enim\",\"ut id eiusmod ipsum irure\",\"incididunt in\",\"of\",\"aliqua id officia occaecat tempor\",\"in\",\"cillum minim\",\"irure ipsum quis consequat in\",\"dolor commodo aliquip\",\"elit voluptate velit\",\"ipsum\",\"mollit aliqua dolor Duis\",\"reprehender\",\"esse reprehenderit in non\",\"aute veniam ea Duis nostrud\",\"qui eu in fugiat\",\"officia occaecat consequ\",\"est\",\"ullamco irure qui dolor\"],\"status\":\"available\",\"tags\":[{\"name\":\"sunt\",\"id\":-1326673790229442600},{\"name\":\"reprehenderit in deserunt esse ad\",\"id\":-3812182748876263400},{\"name\":\"esse et\",\"id\":1591662340194840600},{\"id\":-532117618352377860,\"name\":\"non tempor d\"},{\"id\":-3394805364899614700,\"name\":\"ullamco\"},{\"name\":\"et in\",\"id\":-8819235237071553000},{\"id\":3892791470753624000,\"name\":\"Lorem mollit magna incididunt\"},{\"id\":-8471013567195009000,\"name\":\"voluptate enim\"},{\"id\":8747122112311742000,\"name\":\"in Lorem dolor aute qui\"},{\"id\":-1461618769312350200,\"name\":\"nisi deserunt esse incididunt fugiat\"},{\"name\":\"exercitation ad\",\"id\":-1032383715320475600},{\"name\":\"sit\",\"id\":-268545135143677950},{\"name\":\"proident consectetur nulla in si\",\"id\":-4035620059527610400},{\"id\":5217205923659395000,\"name\":\"reprehenderit veniam dolore sit voluptate\"},{\"id\":6439031202294186000,\"name\":\"mollit labore nulla in eu\"},{\"name\":\"do\",\"id\":-5849073907085595000},{\"name\":\"laboris incididunt et\",\"id\":8728198017034813000},{\"name\":\"sunt occaecat quis laboris\",\"id\":7369639887330996000},{\"id\":2919760693988819000,\"name\":\"Ex\"},{\"id\":-7937880147096228000,\"name\":\"occaecat Ut culpa\"}]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := []byte(tt.request)
			process(&(tt.log), &r, &(tt.response), log.Log{bridge.New(bridgetest.Mock(t, nil))})
		})
	}
}
