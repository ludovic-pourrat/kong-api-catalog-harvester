package main

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"net/http"
)

func match(method string, path string, contentType string, specification *openapi3.T) (bool, error) {
	router, err := gorillamux.NewRouter(specification)
	if err != nil {
		return false, err
	}
	search, err := http.NewRequest(method, path, http.NoBody)
	if err != nil {
		return false, err
	}
	if len(contentType) > 0 {
		search.Header.Set("Content-Type", contentType)
	}
	route, _, err := router.FindRoute(search)
	if err != nil {
		return false, err
	}
	if route == nil {
		return false, nil
	} else {
		return true, nil
	}

}
