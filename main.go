package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/golang/glog"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: example -stderrthreshold=[INFO|WARN|FATAL] -log_dir=[string]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	// NOTE: This next line is key you have to call flag.Parse() for the command line
	// options or "flags" that are defined in the glog module to be picked up.
	flag.Parse()
}

func main() {
	var err error
	//c := &Camera{
	//	"LakeMtn-North",
	//	"http://10.33.130.118/snap.jpeg",
	//}
	var cameras []*Camera
	raw, err := ioutil.ReadFile("cameras.json")
	if err != nil {
		glog.Fatal(err)
	}
	err = json.Unmarshal(raw, &cameras)
	if err != nil {
		glog.Fatal(err)
	}
	var wg sync.WaitGroup
	for _, c := range cameras {
		wg.Add(1)
		go func(c *Camera) {
			defer wg.Done()
			err = getImage("./tmp/", c)
			if err != nil {
				glog.Errorf("[%s] Unable to retrieve image: %s", c.Name, err.Error())
			}
		}(c)
	}
	wg.Wait()
}
