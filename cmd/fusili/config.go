package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type config struct {
	Hosts  map[string][]int
	Output map[string]interface{}
}

func loadConfig(path string) (*config, error) {
	var c config

	if path == "" {
		return nil, fmt.Errorf("no configuration file path provided")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read configuration file: %s", err)
	}

	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("unable to parse JSON: %s", err)
	}

	return &c, nil
}
