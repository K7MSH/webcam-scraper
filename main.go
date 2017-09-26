package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: webcam-scraper -loglevel=[TRACE|DEBUG|INFO|WARNING|ERROR|CRITICAL]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var logger = loggo.GetLogger("main")
var httplogger = loggo.GetLogger("main.http")
var rootLogger = loggo.GetLogger("")

var (
	version string = "unknown"
	date    string = "unknown"
	commit  string = "unknown"
)

func init() {
	var log_level string
	var root_log_level string
	flag.StringVar(&log_level, "loglevel", "TRACE", "Set the logger level")
	flag.StringVar(&root_log_level, "rootloglevel", "INFO", "Set the root logger level")

	loggo.RemoveWriter("default")
	loggo.RegisterWriter("default", loggocolor.NewWriter(os.Stderr))

	flag.Usage = usage
	// NOTE: This next line is key you have to call flag.Parse() for the command line
	// options or "flags" that are defined in the glog module to be picked up.
	flag.Parse()

	loggo.ConfigureLoggers(fmt.Sprintf("<root>=%s;main=%s", root_log_level, log_level))
}

func main() {
	logger.Tracef("Beginning main")
	logger.Tracef("Commit: %s", commit)
	logger.Tracef("Version: %s", version)
	logger.Tracef("Build Date: %s", date)
	var err error

	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		rootLogger.Criticalf("%v", err.Error())
		return
	}
	err = os.Chdir(cwd)
	if err != nil {
		rootLogger.Criticalf("%v", err.Error())
		return
	}

	//c := &Camera{
	//	"LakeMtn-North",
	//	"http://10.33.130.118/snap.jpeg",
	//}
	logger.Tracef("Reading list of cameras from config file")
	var config Config
	config.Load("cameras.json")
	var wg sync.WaitGroup
	logger.Tracef("Starting camera capture")
	if config.StoragePath != "" && config.StoragePath != "." && config.StoragePath != "./" {
		if config.StoragePath[len(config.StoragePath)-1] != os.PathSeparator {
			logger.Tracef("Storage path lacking %c, adding to %s", os.PathSeparator, config.StoragePath)
			config.StoragePath += string(os.PathSeparator)
		}
		if err = ensureDir(config.StoragePath); err != nil {
			logger.Criticalf("Unable to ensure storage path, %v: %v", config.StoragePath, err)
		}
	}
	if err = ensureDir("failures/"); err != nil {
		logger.Criticalf("Unable to ensure failure dir, failures: %v", err)
	}
	for _, c := range config.Cameras {
		wg.Add(1)
		go func(c *Camera) {
			logger.Tracef("Starting camera capture: %s", c.Name)
			defer wg.Done()
			err = imageLoop(&config, c)
			if err != nil {
				logger.Warningf("Image retrieval failed, retrying 2 more times")
				err = imageLoop(&config, c)
				if err != nil {
					logger.Warningf("Image retrieval failed, retrying 1 more time")
					imageLoop(&config, c)
				}
			}
		}(c)
	}
	wg.Wait()
	logger.Tracef("Finished camera capture")
}

func imageLoop(config *Config, c *Camera) error {
	err := getImage(config.StoragePath, c)
	if err != nil {
		logger.Errorf("[%s] Unable to retrieve image: %s", c.Name, err.Error())
		var v VersionInfo
		v.Load(path.Join(config.StoragePath, c.Name))
		filename := path.Join(v.Directory, v.Latest)
		cmd := exec.Command("convert", filename, "-fill", "rgba(20,20,20,0.80)", "-draw", "rectangle 0,340 1920,740", "-fill", "white", "-strokewidth", "4", "-stroke", "black", "-gravity", "Center", "-weight", "800", "-pointsize", "90", "-annotate", "0", "K7MSH CAMERA\nTEMPORARILY OFFLINE", path.Join("failures", fmt.Sprintf("%s___%s", c.Name, v.Latest)))
		err := cmd.Run()
		if err != nil {
			logger.Warningf("Failed to create debug offline file: %s", err.Error())
			log, err := cmd.CombinedOutput()
			rootLogger.Debugf("%s -- %v", err.Error(), log)
			return err
		}
		if err == nil && c.SaveTo != "" {
			image, err := ioutil.ReadFile(path.Join("failures", fmt.Sprintf("%s___%s", c.Name, v.Latest)))
			if err != nil {
				logger.Warningf("Failed to read offline file: %s", err.Error())
				return err
			}
			httplogger.Tracef("[%s] Saving offline image to %s", c.Name, c.SaveTo)
			fp2, err := os.OpenFile(c.SaveTo, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				logger.Warningf("Failed to create debug offline file: %s", err.Error())
				return err
			}
			defer fp2.Close()
			fp2.Write(image)
			httplogger.Infof("[%s] Saved offline image to %s", c.Name, c.SaveTo)
		}
	}
	return nil
}
