package main

import (
	"github.com/fatih/color"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

// Put will put a local file to a remote destination over SSH
func Put(file *file, h string, conf *ssh.ClientConfig) error {
	c, err := ssh.Dial("tcp", h, conf)
	if err != nil {
		color.Red("%s: Error connecting to host - %s", h, err)
		return err
	}

	s, err := sftp.NewClient(c)
	if err != nil {
		color.Red("%s: Error connecting to SFTP host - %s", h, err)
		return err
	}
	defer s.Close()

	f, err := ioutil.ReadFile(file.source)
	if err != nil {
		color.Red("%s: Unable to open source file - %s", h, err)
		return err
	}

	sf, err := s.Create(file.dest)
	if err != nil {
		color.Red("%s: Unable to create remote file - %s", h, err)
		return err
	}

	_, err = sf.Write(f)
	if err != nil {
		color.Red("%s: Error writing remote file - %s", h, err)
		sf.Close()
		return err
	}
	sf.Close()

	err = s.Chmod(file.dest, file.mode)
	if err != nil {
		color.Red("%s: Unable to modify remote file permissions - %s", h, err)
		return err
	}

	if opts.Verbose {
		color.Blue("%s: File copied from %s to %s", h, file.source, file.dest)
	}

	return nil
}
