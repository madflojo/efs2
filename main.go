package main

import (
	"fmt"
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
		fmt.Printf("Unable to setup SSH client configuration - %s\n", err)
		os.Exit(1)
	}

	tasks, err := parseFile(opts.Efs2file)
	if err != nil {
		fmt.Printf("Error parsing Efs2file - %s\n", err)
		os.Exit(1)
	}

	for _, h := range args {
		if opts.Parallel {
			var wg sync.WaitGroup
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				execTasks(tasks, h, sshConf)
				return
			}(&wg)
			wg.Wait()
		} else {
			execTasks(tasks, h, sshConf)
		}
	}
}

func execTasks(tasks []*task, h string, sshConf *ssh.ClientConfig) {
	for n, t := range tasks {
		fmt.Printf("%s: Executing task %d - %s\n", h, n, t.task)
		if opts.Dryrun {
			continue
		}
		if t.file.active {
			err := Put(t.file, h, sshConf)
			if err != nil {
				fmt.Printf("%s: Error placing remote file - %s\n", h, err)
				fmt.Printf("%s: Stopping execution due to too many errors\n", h)
				return
			}
		}
		if t.command.active {
			err := Run(t.command, h, sshConf)
			if err != nil {
				fmt.Printf("%s: Error executing command - %s\n", h, err)
				fmt.Printf("%s: Stopping execution due to too many errors\n", h)
				return
			}
		}
	}
}
