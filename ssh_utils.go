package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os/user"
)

// readPrivateKey will read a private key from file and return parsed version
func readPrivateKey(file string) (ssh.Signer, error) {
	var key ssh.Signer
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return key, err
	}
	return ssh.ParsePrivateKey(b)
}

// initSSH will initialize an SSH client config object to pass to task executors
func initSSH(opts options) (*ssh.ClientConfig, error) {
	var sshConf *ssh.ClientConfig

	usr, err := user.Current()
	if err != nil {
		return sshConf, fmt.Errorf("Unable to grab execution user details - %s", err)
	}
	username := usr.Username

	keyfile := usr.HomeDir + "/.ssh/id_rsa"
	if opts.Keyfile != "" {
		keyfile = opts.Keyfile
	}

	key, err := readPrivateKey(keyfile)
	if err != nil {
		return sshConf, fmt.Errorf("Unable to read private key file - %s", err)
	}

	sshConf = &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(key)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return sshConf, nil
}
