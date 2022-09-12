package main

import (
	"bytes"
	"encoding/json"
)

func prettify(data []byte) []byte {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return []byte("")
	}
	return prettyJSON.Bytes()
}

func New() interface{} {
	return &Config{}
}
