package main

import "os"

type options struct {
	Verbose  bool   `short:"v" long:"verbose" description:"Enable verbose output"`
	Efs2file string `short:"f" long:"file" description:"Specify an alternative Efs2file" default:"./Efs2file"`
	Keyfile  string `short:"i" long:"key" description:"Specify an SSH Private key to use (default: ~/.ssh/id_rsa)"`
	Parallel bool   `short:"p" long:"parallel" description:"Execute tasks in parallel"`
	Dryrun   bool   `short:"d" long:"dryrun" description:"Print tasks to be executed without actually executing any tasks"`
	Port     string `long:"port" description:"Define an alternate SSH Port" default:"22"`
	User     string `short:"u" long:"user" description:"Remote host username (default: current user)"`
}

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
