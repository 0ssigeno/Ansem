package main

import (
	"./internal"
	"context"
	"fmt"
	"os"
	"sync"
	"text/tabwriter"
)

func main() {
	var c internal.Conf
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
		"Timeout:\t%d\n",

		c.Directory, c.GameServer, c.Threshold, c.TeamDir, c.SubmissionType, c.FlagRegex, c.Tick, c.Workers, c.Token, c.Timeout)
	writer.Flush()
	toSubmit := make(chan internal.Flag, c.Workers*5)

	wg := sync.WaitGroup{}
	wg.Add(2)

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
	submitterCtx = context.WithValue(submitterCtx, "token", c.Token)
	submitterCtx = context.WithValue(submitterCtx, "workers", c.Workers)

	go internal.StartExploiter(exploitCtx, &wg)
	go internal.StartSubmitter(submitterCtx, &wg)

	wg.Wait()
}
