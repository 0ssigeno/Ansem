package main

import (
	"./internal/exploit"
	"./internal/promStats"
	"./internal/submit"
	"context"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"text/tabwriter"
)

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
	FlagDuplicated string `yaml:"flag_duplicated"`
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

func main() {
	var c Conf
	c.GetConf()

	//Fix path if the last char is not "/"
	if c.Directory[len(c.Directory)-1] != '/' {
		c.Directory = fmt.Sprintf("%s/", c.Directory)
	}

	//Aligned print
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 8, 0, '\t', 0)
	_, _ = fmt.Fprintf(writer, "Hi, I'm starting with these settings:\n\n"+
		"Exploits Dir:\t%s\n"+
		"Gameserver:\t%s\n"+
		"Threshold:\t%d\n"+
		"TeamDir:\t%s\n"+
		"SubmissionType:\t%s\n"+
		"Flag Regex:\t%s\n"+
		"Tick:\t%d\n"+
		"Workers:\t%d\n"+
		"Token:\t%s\n"+
		"Timeout:\t%d\n"+
		"FlagAccepted:\t%s\n"+
		"FlagDuplicated:\t%s\n",

		c.Directory, c.GameServer, c.Threshold, c.TeamDir, c.SubmissionType, c.FlagRegex, c.Tick,
		c.Workers, c.Token, c.Timeout, c.FlagAccepted, c.FlagDuplicated)
	_ = writer.Flush()

	wg := sync.WaitGroup{}
	wg.Add(2)

	toSubmit := make(chan exploit.Flag, c.Workers*5)

	exploitCtx := context.Background()
	exploitCtx = context.WithValue(exploitCtx, "exploitDir", c.Directory)
	exploitCtx = context.WithValue(exploitCtx, "tick", c.Tick)
	exploitCtx = context.WithValue(exploitCtx, "teamDir", c.TeamDir)
	exploitCtx = context.WithValue(exploitCtx, "workers", c.Workers)
	exploitCtx = context.WithValue(exploitCtx, "submit", toSubmit)
	exploitCtx = context.WithValue(exploitCtx, "flagRegex", c.FlagRegex)
	exploitCtx = context.WithValue(exploitCtx, "timeout", c.Timeout)
	exploitCtx = context.WithValue(exploitCtx, "threshold", c.Threshold)

	submitterCtx := context.Background()
	submitterCtx = context.WithValue(submitterCtx, "gameServer", c.GameServer)
	submitterCtx = context.WithValue(submitterCtx, "submit", toSubmit)
	submitterCtx = context.WithValue(submitterCtx, "flagRegex", c.FlagRegex)
	submitterCtx = context.WithValue(submitterCtx, "subType", c.SubmissionType)
	submitterCtx = context.WithValue(submitterCtx, "flagAccepted", c.FlagAccepted)
	submitterCtx = context.WithValue(submitterCtx, "flagDuplicated", c.FlagDuplicated)
	submitterCtx = context.WithValue(submitterCtx, "token", c.Token)
	submitterCtx = context.WithValue(submitterCtx, "workers", c.Workers)

	go exploit.StartExploiter(exploitCtx, &wg)
	go submit.StartSubmitter(submitterCtx, &wg)
	go promStats.StartStatistics()

	wg.Wait()
}
