package main

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"fmt"
)

// Puts a local file to a remote destination over SSH
func Put(file *file, h string, conf *ssh.ClientConfig) error {
	c, err := ssh.Dial("tcp", h+":22", conf)
	if err != nil {
		fmt.Printf("%s: Error connecting to host - %s\n", h, err)
		return err
	}

	s, err := sftp.NewClient(c)
	if err != nil {
		fmt.Printf("%s: Error connecting to SFTP host - %s\n", h, err)
		return err
	}
	defer s.Close()

	f, err := ioutil.ReadFile(file.source)
	if err != nil {
		fmt.Printf("%s: Unable to open source file - %s\n", h, err)
		return err
	}

	sf, err := s.Create(file.dest)
	if err != nil {
		fmt.Printf("%s: Unable to create remote file - %s\n", h, err)
		return err
	}

	_, err = sf.Write(f)
	if err != nil {
		fmt.Printf("%s: Error writing remote file - %s\n", h, err)
    sf.Close()
		return err
	}
	sf.Close()

	err = s.Chmod(file.dest, file.mode)
	if err != nil {
		fmt.Printf("%s: Unable to modify remote file permissions - %s\n", h, err)
		return err
	}

	if opts.Verbose {
		fmt.Printf("%s: File copied from %s to %s\n", h, file.source, file.dest)
	}

	return nil
}
