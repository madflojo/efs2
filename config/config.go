package config

// Config is a configuration structure used to configure efs2.
type Config struct {

	// Verbose controls the verbosity of command and execution output.
	Verbose bool

	// Efs2File is the location of the Efs2file used during command execution.
	Efs2File string

	// KeyFile is the user ssh keyfile to use for authentication.
	KeyFile string

	// Parellel will denote whether execution on multiple hosts should be performed in parallel or not.
	Parallel bool

	// DryRun will denote whether to actually perform the executions steps or perform a dry run.
	DryRun bool

	// Port allows users to specify a non-standard SSH port.
	Port string

	// User is the remote user name to use for SSH authentication.
	User string

	// Passphrase is used to hold user passphrase for SSH keys
	Passphrase []byte

	// Hosts is a list of hosts to execute against
	Hosts []string
}

// New will return a Config object with defaults pre-loaded.
func New() Config {
	c := Config{}
	c.Efs2File = "./Efs2file"
	c.Port = "22"
	return c
}
