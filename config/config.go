/*
Package config is the internal configuration used for Efs2. This configuration is for the internal application execution. It exists to pave the way for non-cli instances of Efs2 in the future.
*/
package config

import (
	"fmt"
	"os/user"
)

// Config provides a configuration structure used within the Efs2 application.
type Config struct {

	// Verbose if set to true, will enable verbose output of command execution.
	Verbose bool

	// Quiet if set to true, will silence all output and focus only on exit codes.
	Quiet bool

	// Efs2File specifies the location of the Efs2File to execute. The default value is the current directory.
	Efs2File string

	// KeyFile specifies the user's SSH key location. The default value is the .ssh directory within the users home directory.
	KeyFile string

	// Parallel if set to true, will trigger Efs2 to execute the Efs2file against all hosts in parallel.
	Parallel bool

	// DryRun if set to true, will prevent the execution of tasks against hosts. The end-user will see instructions logged to stdout.
	DryRun bool

	// Port allows users to specify a default port when hosts do not have a port already specified.
	Port string

	// User is the remote user name to use for authentication.
	User string

	// Passphrase holds the user-provided SSH key passphrase.
	Passphrase []byte

	// Hosts will contain a list of remote servers to use for execution.
	Hosts []string
}

// New will return Config populated with pre-defined defaults.
func New() Config {
	c := Config{}
	c.Efs2File = "./Efs2file"
	c.Port = "22"
	usr, homeDir, err := UserDetails("")
	if err != nil {
		return c
	}
	c.User = usr
	c.KeyFile = homeDir + "/.ssh/id_rsa"
	return c
}

// UserDetails is a helper function to fetch the username and home directory of a specified user. If no user is specified, UserDetails will return the current runtime user.
func UserDetails(u string) (string, string, error) {
	// Runtime User
	if u == "" {
		usr, err := user.Current()
		if err != nil {
			return "", "", fmt.Errorf("Unable to determine user details - %s", err)
		}
		return usr.Username, usr.HomeDir, nil
	}

	// Specified User
	if u != "" {
		usr, err := user.Lookup(u)
		if err != nil {
			return "", "", fmt.Errorf("Unable to determine user details - %s", err)
		}
		return usr.Username, usr.HomeDir, nil
	}

	// Oh no...
	return "", "", fmt.Errorf("User not found")
}
