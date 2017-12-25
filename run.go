package main

import (
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
	"strings"
)

// Executes a remote command over ssh
func Run(task *command, h string, conf *ssh.ClientConfig) error {
	c, err := ssh.Dial("tcp", h, conf)
	if err != nil {
		color.Red("%s: Error connecting to host - %s", h, err)
		return err
	}
	s, err := c.NewSession()
	if err != nil {
		color.Red("%s: Error creating new ssh session to host - %s", h, err)
		return err
	}
	defer s.Close()

	r, err := s.CombinedOutput(task.cmd)
	if err != nil {
		color.Red("%s Failed to execute task on host - %s", h, err)
		return err
	}
	if opts.Verbose {
		color.Blue("%s: Task Output", h)
		color.Blue("%s: ------------------------", h)
		for _, x := range strings.Split(string(r), "\n") {
			color.Cyan("%s: %s\n", h, x)
		}
		color.Blue("%s: ------------------------", h)
	}
	return nil
}
