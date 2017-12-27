package main

import (
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"golang.org/x/crypto/ssh"
	"os"
	"sync"
)

var opts options

func main() {
	args, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	sshConf, err := initSSH(opts)
	if err != nil {
		color.Red("Unable to setup SSH client configuration - %s", err)
		os.Exit(1)
	}

	tasks, err := Parse(opts.Efs2file)
	if err != nil {
		color.Red("Error parsing Efs2file - %s", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for _, h := range args {
		x := h + ":" + opts.Port
		if opts.Parallel {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				Exec(tasks, x, sshConf)
				wg.Done()
				return
			}(&wg)
		} else {
			Exec(tasks, x, sshConf)
		}
	}
	wg.Wait()
}

// Exec will loop through a slice of tasks, kicking off whichever instruction is required
func Exec(tasks []*task, h string, sshConf *ssh.ClientConfig) {
	for n, t := range tasks {
		color.Blue("%s: Executing task %d - %s", h, n, t.task)
		if opts.Dryrun {
			continue
		}
		if t.file.active {
			err := Put(t.file, h, sshConf)
			if err != nil {
				color.Red("%s: Error placing remote file - %s", h, err)
				color.Red("%s: Stopping execution due to too many errors", h)
				return
			}
		}
		if t.command.active {
			err := Run(t.command, h, sshConf)
			if err != nil {
				color.Red("%s: Error executing command - %s", h, err)
				color.Red("%s: Stopping execution due to too many errors", h)
				return
			}
		}
	}
}
