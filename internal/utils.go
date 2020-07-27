package internal

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Flag struct {
	flag    string
	team    string
	exploit string
}

type Conf struct {
	Directory      string `yaml:"exploits_dir"`
	Tick           int    `yaml:"tick"`
	TeamDir        string `yaml:"team_dir"`
	GameServer     string `yaml:"gameserver"`
	Threshold      int    `yaml:"threshold"`
	Workers        int    `yaml:"workers"`
	SubmissionType string `yaml:"submission_type"`
	FlagRegex      string `yaml:"flag_regex"`
	FlagAccepted   string `yaml:"flag_accepted"`
	Token          string `yaml:"token"`
	Timeout        int    `yaml:"timeout"`
}

// Function that initialize the config
func (c *Conf) GetConf() *Conf {

	yamlFile, err := ioutil.ReadFile("configs/conf.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
