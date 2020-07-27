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

func (c *Conf) GetContext() context.Context {
	ctx := context.Background()

	toSubmit := make(chan exploit.Flag, c.Workers*5)

	ctx = context.WithValue(ctx, "exploitDir", c.Directory)
	ctx = context.WithValue(ctx, "tick", c.Tick)
	ctx = context.WithValue(ctx, "teamDir", c.TeamDir)
	ctx = context.WithValue(ctx, "workers", c.Workers)
	ctx = context.WithValue(ctx, "submit", toSubmit)
	ctx = context.WithValue(ctx, "flagRegex", c.FlagRegex)
	ctx = context.WithValue(ctx, "timeout", c.Timeout)
	ctx = context.WithValue(ctx, "threshold", c.Threshold)
	ctx = context.WithValue(ctx, "gameServer", c.GameServer)
	ctx = context.WithValue(ctx, "flagRegex", c.FlagRegex)
	ctx = context.WithValue(ctx, "subType", c.SubmissionType)
	ctx = context.WithValue(ctx, "flagAccepted", c.FlagAccepted)
	ctx = context.WithValue(ctx, "flagDuplicated", c.FlagDuplicated)
	ctx = context.WithValue(ctx, "token", c.Token)
	return ctx
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

	exploitCtx := c.GetContext()

	submitterCtx := c.GetContext()

	go exploit.StartExploiter(exploitCtx, &wg)
	go submit.StartSubmitter(submitterCtx, &wg)
	go promStats.StartStatistics()

	wg.Wait()
}
