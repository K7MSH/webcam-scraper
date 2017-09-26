package main

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/juju/loggo"
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

var camlogger = loggo.GetLogger("main.camera")

var HttpClient = &http.Client{
	//Timeout: time.Second * 60,
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 60 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   60 * time.Second,
		ResponseHeaderTimeout: 20 * time.Second,
		ExpectContinueTimeout: 2 * time.Second,
	},
}

func (cam *Camera) getFilename() string {
	// "%Y%m%d-%H%M%S"
	format := "20060102-150405MST"
	filename := fmt.Sprintf("%s.jpg", time.Now().Format(format))
	return filename
}
func (cam *Camera) ensureDirectories(filepath string) error {
	if err := ensureDir(filepath); err != nil {
		return err
	}
	if cam.SaveTo != "" {
		tmpstr, _ := path.Split(cam.SaveTo)
		if err := ensureDir(tmpstr); err != nil {
			return err
		}
	}
	return nil
}
func (cam *Camera) requestImage() (*http.Response, error) {
	camlogger.Tracef("[%s] Initiating request to %s", cam.Name, cam.URL)
	return HttpClient.Get(cam.URL)
}

func (cam *Camera) saveImage(writer io.WriterTo, filename string) error {
	camlogger.Tracef("[%s] Saving image from %s", cam.Name, cam.URL)
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = writer.WriteTo(fp)
	if err != nil {
		return err
	}
	return nil
}
func (cam *Camera) copyImage(filename string) error {
	camlogger.Tracef("[%s] opening image for copy", cam.Name)
	fp, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()
	camlogger.Tracef("[%s] Saving image copy to %s", cam.Name, cam.SaveTo)
	fp2, err := os.OpenFile(cam.SaveTo, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fp2.Close()
	_, err = io.Copy(fp2, fp)
	if err != nil {
		return err
	}
	camlogger.Infof("[%s] Saved image copy to %s", cam.Name, cam.SaveTo)
	return nil
}

func (cam *Camera) bufferImage(r *http.Response) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	count, err := buf.ReadFrom(r.Body)
	if err != nil {
		return nil, err
	}
	if count != r.ContentLength {
		return nil, errors.New("Data size downloaded is not equal to content length")
	}
	camlogger.Tracef("[%s] got image from %s", cam.Name, cam.URL)
	return buf, nil
}

func verifyImageIntegrity(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = jpeg.Decode(f)
	if err != nil {
		return err
	}
	return nil
}

func (cam *Camera) GetImage(dir string) error {
	var filename string
	var filepath string
	var err error
	filepath = path.Join(dir, cam.Name)

	err = cam.ensureDirectories(filepath)
	if err != nil {
		return err
	}

	filename = cam.getFilename()
	var version *VersionInfo = &VersionInfo{filepath, cam.Name, filename}
	filename = path.Join(filepath, filename)

	if cam.Auth != nil {
		camlogger.Warningf("[%s] Found Auth, Not implemented, Bailing!", cam.Name)
		return nil
	}

	response, err := cam.requestImage()
	if err != nil {
		return err
	}
	defer response.Body.Close()

	buf, err := cam.bufferImage(response)
	if err != nil {
		return err
	}

	err = cam.saveImage(buf, filename)
	if err != nil {
		return err
	}

	err = verifyImageIntegrity(filename)
	if err != nil {
		return err
	}
	if cam.SaveTo != "" {
		err = cam.copyImage(filename)
		if err != nil {
			return err
		}
	}
	version.Save()
	camlogger.Infof("[%s] Saved image to %s", cam.Name, filename)
	return nil
}
