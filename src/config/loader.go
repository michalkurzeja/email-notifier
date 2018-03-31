package config

import (
	"encoding/json"
	"io/ioutil"
)

const path = "config/config.json"

func Load() Config {
	raw, err := ioutil.ReadFile(path)
    if err != nil {
        panic(err)
    }

	var config Config
	json.Unmarshal(raw, &config)

	return config
}