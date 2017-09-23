package main

import (
	"flag"
	"fmt"
	"os"
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
		if config.StoragePath[len(config.StoragePath)-1] != '/' {
			logger.Tracef("Storage path lacking /, adding to %s", config.StoragePath)
			config.StoragePath += "/"
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
			err = getImage(config.StoragePath, c)
			if err != nil {
				logger.Errorf("[%s] Unable to retrieve image: %s", c.Name, err.Error())
			}
		}(c)
	}
	wg.Wait()
	logger.Tracef("Finished camera capture")
}
