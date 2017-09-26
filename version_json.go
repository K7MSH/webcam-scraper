package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type VersionInfo struct {
	Directory string
	Camera    string
	Latest    string
}

func (vi *VersionInfo) Save() {
	filename := path.Join(vi.Directory, "version.json")
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		httplogger.Warningf("Failed to open version file: %s", err.Error())
		return
	}
	defer f.Close()
	f.Truncate(0)
	data, err := json.Marshal(vi)
	if err != nil {
		httplogger.Warningf("Failed to marshal version json file: %s", err.Error())
		return
	}
	_, err = f.Write(data)
	if err != nil {
		httplogger.Warningf("Failed to write to version file: %s", err.Error())
		return
	}
}

func (vi *VersionInfo) Load(dir string) error {
	filename := path.Join(dir, "version.json")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		httplogger.Warningf("Failed to read version file: %s", err.Error())
		return err
	}
	err = json.Unmarshal(data, vi)
	if err != nil {
		httplogger.Warningf("Failed to unmarshal version json file: %s", err.Error())
		return err
	}
	return nil
}
