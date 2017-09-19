package main

import (
	"encoding/json"
	"io/ioutil"
)

type Camera struct {
	Name   string
	URL    string
	SaveTo string
	Auth   *CameraAuth
}
type CameraAuth struct {
	CameraName string
	User       string
	Password   string
}

type Cameras []*Camera

func (c *Cameras) Decode(raw []byte) {
	logger.Tracef("Parsing json response")
	err := json.Unmarshal(raw, &c)
	if err != nil {
		logger.Criticalf("Failed to parse file: %v", err)
	}
}

func (c *Cameras) Load(file string) {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Criticalf("Failed to read file: %v", err)
	}
	c.Decode(raw)
}
