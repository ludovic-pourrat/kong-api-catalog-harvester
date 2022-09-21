package utils

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/invopop/yaml"
	"os"
)

func Init() error {
	err := os.Mkdir("/repository", 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	err = os.Mkdir("/snapshots", 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func Write(name string, specification *openapi3.T) error {
	// validate
	err := specification.Validate(context.Background())
	if err != nil {
		//return err
	}
	// marshal to yaml
	data, err := yaml.Marshal(specification)
	if err != nil {
		return err
	}
	digest := md5.New()
	digest.Write(data)
	sum := digest.Sum(nil)
	_, err = os.Stat(fmt.Sprintf("/snapshots/%s-%x.yaml", name, sum))
	if os.IsNotExist(err) {
		// write to file
		err = os.WriteFile(fmt.Sprintf("/repository/%s.yaml", name), data, 0644)
		if err != nil {
			return err
		}
		// write to file
		err = os.WriteFile(fmt.Sprintf("/snapshots/%s-%x.yaml", name, sum), data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func Read(name string) (*openapi3.T, error) {
	_, err := os.Stat(fmt.Sprintf("/repository/%s.yaml", name))
	if err == nil {
		var data []byte
		// write to file
		data, err = os.ReadFile(fmt.Sprintf("/repository/%s.yaml", name))
		if err != nil {
			return nil, err
		}
		var specification *openapi3.T
		// marshal to yaml
		err = yaml.Unmarshal(data, specification)
		if err != nil {
			return nil, err
		}
		return specification, nil
	}
	return nil, nil
}
