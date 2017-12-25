package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
)

// Executes a remote command over ssh
func Run(task *command, h string, conf *ssh.ClientConfig) error {
	c, err := ssh.Dial("tcp", h, conf)
	if err != nil {
		fmt.Printf("%s: Error connecting to host - %s\n", h, err)
		return err
	}
	s, err := c.NewSession()
	if err != nil {
		fmt.Printf("%s: Error creating new ssh session to host - %s\n", h, err)
		return err
	}
	defer s.Close()

	r, err := s.CombinedOutput(task.cmd)
	if err != nil {
		fmt.Printf("%s Failed to execute task on host - %s\n", h, err)
		return err
	}
	if opts.Verbose {
		fmt.Printf("%s: Task Output\n", h)
		fmt.Printf("%s: ------------------------\n", h)
		for _, x := range strings.Split(string(r), "\n") {
			fmt.Printf("%s: %s\n", h, x)
		}
		fmt.Printf("%s: ------------------------\n", h)
	}
	return nil
}
