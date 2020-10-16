package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/madflojo/efs2/config"
	"github.com/madflojo/efs2/parser"
	"github.com/madflojo/efs2/ssh"
	"os/user"
	"regexp"
	"sync"
)

// Encrypted Key Error
var isPassErr = regexp.MustCompile(`.*decode encrypted private keys$`)

// Has a port defined
var hasPort = regexp.MustCompile(`.*:\d*`)

func Run(cfg config.Config) error {
	var clientCfg ssh.Config

	// Check if user is defined
	if cfg.User == "" {
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("Unable to determine user details - %s", err)
		}
		cfg.User = usr.Username
	}

	if cfg.Verbose {
		color.Yellow("SSH User: %s", cfg.User)
	}

	// Check if Keyfile is defined
	if cfg.KeyFile == "" {
		usr, err := user.Lookup(cfg.User)
		if err != nil {
			return fmt.Errorf("Unable to determine key file location - %s", err)
		}

		cfg.KeyFile = usr.HomeDir + "/.ssh/id_rsa"
	}

	if cfg.Verbose {
		color.Yellow("Key Path: %s", cfg.KeyFile)
	}

	// Setup SSH Config
	clientCfg, err := ssh.ReadKeyFile(cfg.KeyFile, cfg.Passphrase)
	if err != nil {
		if !isPassErr.MatchString(err.Error()) {
			return fmt.Errorf("Unable to obtain Key Passphrase - %s", err)
		}
		cfg.Passphrase, err = gopass.GetPasswd()
		if err != nil {
			return fmt.Errorf("Unable to obtain Key Passphrase - %s", err)
		}
		clientCfg, err = ssh.ReadKeyFile(cfg.KeyFile, cfg.Passphrase)
		if err != nil {
			return fmt.Errorf("Unable to read keyfile - %s", err)
		}
	}
	clientCfg.User = cfg.User

	// Check if Efs2file is defined
	if cfg.Efs2File == "" {
		cfg.Efs2File = "./Efs2file"
	}

	if cfg.Verbose {
		color.Yellow("Efs2file Path: %s", cfg.Efs2File)
	}

	// Parse Efs2file
	tasks, err := parser.Parse(cfg.Efs2File)
	if err != nil {
		return fmt.Errorf("Unable to parse Efs2file - %s", err)
	}

	// Fixup Hosts
	cfg.Hosts = fixUpHosts(cfg.Hosts, cfg.Port)

	// Execute
	var wg sync.WaitGroup
	var errCount int
	for _, h := range cfg.Hosts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := clientCfg
			c.Host = h
			sh, err := ssh.Dial(c)
			if err != nil {
				errCount = errCount + 1
				color.Red("%s: Error connecting to host - %s", h, err)
				return
			}
			for i, t := range tasks {
				color.Blue("%s: Executing Task %d - %s", h, i, t.Task)
				if cfg.DryRun {
					continue
				}
				if t.File.Source != "" {
					err := sh.Put(t.File)
					if err != nil {
						errCount = errCount + 1
						color.Red("%s: Error uploading file - %s", h, err)
						return
					}
					color.Blue("%s: File upload successful", h)
				}
				if t.Command.Cmd != "" {
					r, err := sh.Run(t.Command)
					if err != nil {
						errCount = errCount + 1
						color.Red("%s: Error executing command - %s", h, err)
						return
					}
					color.Blue("%s: %s", h, r)
				}
			}

		}()
		if !cfg.Parallel {
			wg.Wait()
		}
	}
	wg.Wait()

	if errCount > 0 {
		return fmt.Errorf("Execution failed with %d errors", errCount)
	}
	return nil
}

func fixUpHosts(hosts []string, port string) []string {
	// Fixup Hosts
	var hh []string
	for _, h := range hosts {
		if hasPort.MatchString(h) {
			hh = append(hh, h)
			continue
		}
		if port == "" {
			hh = append(hh, h+":22")
			continue
		}
		hh = append(hh, h+":"+port)
	}
	return hh
}
