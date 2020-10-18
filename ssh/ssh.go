/*
Package ssh provides users with a simple interface for remote execution over SSH. While this package is for Efs2 it could easily be imported and used by others.

  // Read a Key File
  key, err := ssh.ReadKeyFile("~/.ssh/id_rsa", "aPass")
  if err != nil {
    // do something
  }

  // Dial the remote host
  conn, err := ssh.Dial(ssh.Config{
    Host: "example.com:22",
    User: "example",
    Key:  key,
  })
  if err != nil {
    // do something
  }

  // Upload a File
  err = conn.Put(ssh.File{...})
  if err != nil {
    // do something
  }

  // Run a Command
  output, err := conn.Run(ssh.Command{...})
  if err != nil {
    // do something
  }

*/
package ssh

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
)

// Conn is an SSH connection the internally holds both an SSH and SFTP session.
type Conn struct {
	// client is an SSH client connection. This connection is the underlying SSH and SFTP TCP Client.
	client *ssh.Client

	// sftpClient is an SFTP upgraded SSH client. This client is used to perform file-based operations.
	sftpClient *sftp.Client
}

// Config provides SSH client configuration details to the Dial function. This configuration offers items such as hostname, port, username, and SSH Keys to use for authentication.
type Config struct {
	// Host contains the remote SSH host address in a hostname:port format. If no port is specified, the Dial function will return an error.
	Host string

	// User contains the remote host's username to use for authentication and must not be left blank.
	User string

	// Key provides an SSH key to use for authentication. This key is the private key and requires the public key to be pre-pushed to the remote host for authentication to work.
	Key ssh.Signer
}

// Task is a wrapper structure used during parsing of Efs2 files. Within this structure contains SSH command and file structures.
type Task struct {
	// Task is the original instruction used to create the task.
	Task string

	// Command is a command structure that provides command execution instructions.
	Command Command

	// File is a file structure that provides file upload instructions.
	File File
}

// Command is a structure that holds details for individual commands that execute remotely.
type Command struct {
	// Cmd is the single line shell command for execution.
	Cmd string
}

// File is a structure that holds the details required for uploading files.
type File struct {
	// Source is the full path to the source file to be uploaded.
	Source string

	// Destination is the full path location of the destination filename.
	Destination string

	// Mode is the file permissions to be set upon upload.
	Mode os.FileMode
}

// Dial will open an SSH connection to the configured remote host. The returned connection can then upload files or execute commands.
func Dial(c Config) (*Conn, error) {
	var err error
	s := &Conn{}
	cfg := &ssh.ClientConfig{
		User: c.User,
	}

	// Check if ignore host key is set
	cfg.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	// Set SSH Key to use
	cfg.Auth = []ssh.AuthMethod{ssh.PublicKeys(c.Key)}

	// Open a connection
	s.client, err = ssh.Dial("tcp", c.Host, cfg)
	if err != nil {
		return &Conn{}, fmt.Errorf("unable to open SSH connection to %s - %s", c.Host, err)
	}

	// Create an FTP Client
	s.sftpClient, err = sftp.NewClient(s.client)
	if err != nil {
		return &Conn{}, fmt.Errorf("unable to open SFTP Session to %s - %s", c.Host, err)
	}

	return s, nil
}

// Close will close the open SSH connection cleaning up any lingering executions.
func (c *Conn) Close() {
	defer c.sftpClient.Close()
	defer c.client.Close()
}

// Put will upload the specified file with the provided destination path and permissions. Currently, this function does not support directories. Each execution is a single file.
func (c *Conn) Put(f File) error {
	fh, err := ioutil.ReadFile(f.Source)
	if err != nil {
		return fmt.Errorf("unable to open source file %s - %s", f.Source, err)
	}

	r, err := c.sftpClient.Create(f.Destination)
	if err != nil {
		return fmt.Errorf("unable to create remote file %s - %s", f.Destination, err)
	}

	_, err = r.Write(fh)
	defer r.Close()
	if err != nil {
		return fmt.Errorf("error writing to remote file %s - %s", f.Destination, err)
	}

	err = r.Chmod(f.Mode)
	if err != nil {
		return fmt.Errorf("unable to change remote file permissions on %s - %s", f.Destination, err)
	}

	return nil
}

// Run will execute the specified file returning both standard output and standard error as a combined value. Any error in execution will return an error, even if the command is partially successful.
func (c *Conn) Run(t Command) ([]byte, error) {
	// Create an SSH Channel
	s, err := c.client.NewSession()
	if err != nil {
		return []byte(""), fmt.Errorf("unable to open SSH channel - %s", err)
	}
	r, err := s.CombinedOutput(t.Cmd)
	if err != nil {
		return []byte(""), fmt.Errorf("failed to execute command - %s", err)
	}

	return r, nil
}

// ReadKeyFile will open and load the SSH key from the provided file. If included, the password will allow this function to decrypt SSH keys.
func ReadKeyFile(file string, pass []byte) (Config, error) {
	var cfg Config
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return cfg, err
	}
	if len(pass) >= 1 {
		cfg.Key, err = ssh.ParsePrivateKeyWithPassphrase(b, pass)
		if err != nil {
			return cfg, fmt.Errorf("Could not parse key with provided passphrase - %s", err)
		}
		return cfg, nil
	}
	cfg.Key, err = ssh.ParsePrivateKey(b)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
