package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	StoragePath string
	Cameras     Cameras
}

func (c *Config) DecodeJson(raw []byte) {
	logger.Tracef("Parsing json response")
	err := json.Unmarshal(raw, &c)
	if err != nil {
		logger.Criticalf("Failed to parse file: %v", err)
	}
}

func (c *Config) LoadJson(file string) {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Criticalf("Failed to read file: %v", err)
	}
	c.DecodeJson(raw)
}
