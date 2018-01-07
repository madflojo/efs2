package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os/user"
)

// readKeyfile will read a private key from file and return an ssh.Signer object
func readKeyfile(file string, pass []byte) (ssh.Signer, error) {
	var key ssh.Signer
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return key, err
	}
	if len(pass) >= 1 {
		return ssh.ParsePrivateKeyWithPassphrase(b, pass)
	}
	return ssh.ParsePrivateKey(b)
}

// initSSH will initialize an SSH client config object
func initSSH(opts options) (*ssh.ClientConfig, error) {
	var sshConf *ssh.ClientConfig

	usr, err := user.Current()
	if err != nil {
		return sshConf, fmt.Errorf("Unable to grab execution user details - %s", err)
	}
	username := usr.Username

	if opts.User != "" {
		username = opts.User
	}

	keyfile := usr.HomeDir + "/.ssh/id_rsa"
	if opts.Keyfile != "" {
		keyfile = opts.Keyfile
	}

	key, err := readKeyfile(keyfile, opts.Passphrase)
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
