package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/ionrock/configurmaktion/agent"
	"github.com/ionrock/configurmaktion/maketypes"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

var builddate = ""
var gitref = ""

type mkconfig struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func loadConfig(path string) ([]agent.Makefile, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	configs := make([]mkconfig, 0)

	err = yaml.Unmarshal(b, &configs)
	if err != nil {
		return nil, err
	}

	makefiles := make([]agent.Makefile, len(configs))
	for i, conf := range configs {
		log.Infof("Adding config %s: %s", conf.Name, conf.URL)
		makefile, err := url.Parse(conf.URL)
		if err != nil {
			return nil, err
		}

		switch makefile.Scheme {
		default:
			// try a local path if no scheme is present
			makefiles[i] = maketypes.LocalMakefile{Name: conf.Name, Path: makefile.Path}
		}
	}

	return makefiles, nil
}

func runConfmaked(c *cli.Context) error {
	mfs, err := loadConfig(c.String("config"))

	if err != nil {
		return err
	}

	a := agent.New(mfs, "ws")
	return a.Run()
}

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s-%s", gitref, builddate)
	app.Name = "confmaked"
	app.Action = runConfmaked
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "The confmaked YAML config file",
			Value: "confmaked.yml",
		},
	}
	app.Run(os.Args)
}
