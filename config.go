package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Devices map[string]string
}

func loadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ERROR - unable to load config file, %s", filename)
	}

	config := &Config{}
	err = json.Unmarshal(data, config)
	return config, err
}
