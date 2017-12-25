package main

import (
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"golang.org/x/crypto/ssh"
	"os"
	"sync"
)

type options struct {
	Verbose  bool   `short:"v" long:"verbose" description:"Enable verbose output"`
	Efs2file string `short:"f" long:"file" description:"Specify an alternative Efs2file (default: ./Efs2file)" default:"./Efs2file"`
	Keyfile  string `short:"i" long:"key" description:"Specify an SSH Private key to use (default: ~/.ssh/id_rsa)"`
	Parallel bool   `short:"p" long:"parallel" description:"Execute tasks in parallel (default: false)"`
	Dryrun   bool   `short:"d" long:"dryrun" description:"Print tasks to be executed without actually executing any tasks"`
	Port     string `long:"port" description:"Define an alternate SSH Port (default: 22)" default:"22"`
	User     string `short:"u" long:"user" description:"Remote host username (default: current user)"`
}

var opts options

type task struct {
	task    string
	command *command
	file    *file
}

type command struct {
	active bool
	cmd    string
}

type file struct {
	active bool
	source string
	dest   string
	mode   os.FileMode
}

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

	tasks, err := parseFile(opts.Efs2file)
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
				execTasks(tasks, x, sshConf)
				wg.Done()
				return
			}(&wg)
		} else {
			execTasks(tasks, x, sshConf)
		}
	}
	wg.Wait()
}

func execTasks(tasks []*task, h string, sshConf *ssh.ClientConfig) {
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
