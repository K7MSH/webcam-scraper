package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/golang/glog"
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
		glog.Warningf("Failed to open version file: %s", err.Error())
		return
	}
	defer f.Close()
	data, err := json.Marshal(vi)
	if err != nil {
		glog.Warningf("Failed to marshal version json file: %s", err.Error())
		return
	}
	_, err = f.Write(data)
	if err != nil {
		glog.Warningf("Failed to write to version file: %s", err.Error())
		return
	}
}

func (vi *VersionInfo) Load() {
}

var HttpClient = &http.Client{
	Timeout: time.Second * 60,
}

func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		return err
	}
	return nil
}

func getImage(dir string, cam *Camera) error {
	var filename string
	var filepath string
	var auth *CameraAuth
	auth = cam.Auth
	filepath = path.Join(dir, cam.Name)
	filename = fmt.Sprintf("%s.jpg", time.Now().Format(time.RFC3339))
	var version *VersionInfo = &VersionInfo{filepath, cam.Name, filename}
	filename = path.Join(filepath, filename)
	if err := ensureDir(filepath); err != nil {
		return err
	}
	if cam.SaveTo != "" {
		tmpstr, _ := path.Split(cam.SaveTo)
		if err := ensureDir(tmpstr); err != nil {
			return err
		}
	}
	if auth != nil {
		glog.Infoln("Found Auth, Not implemented, Bailing!")
		return nil

	} else {
		glog.Infof("[%s] Initiating request to %s", cam.Name, cam.URL)
		response, err := HttpClient.Get(cam.URL)
		glog.Infof("[%s] got image from %s", cam.Name, cam.URL)
		if err != nil {
			return err
		}
		glog.Infof("[%s] Saving image from %s", cam.Name, cam.URL)
		fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer fp.Close()
		io.Copy(fp, response.Body)
		if cam.SaveTo != "" {
			glog.Infof("[%s] Saving image to %s", cam.Name, cam.SaveTo)
			fp2, err := os.OpenFile(cam.SaveTo, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			defer fp2.Close()
			io.Copy(fp2, response.Body)
		}
		version.Save()
		glog.Infof("[%s] Saved image to %s", cam.Name, filename)
	}
	return nil
}
