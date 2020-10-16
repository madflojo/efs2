package ssh

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
)

type Conn struct {
	client     *ssh.Client
	sftpClient *sftp.Client
}

type Config struct {
	Host          string
	User          string
	Passphrase    string
	IgnoreHostKey bool
	Key           ssh.Signer
}

type Task struct {
	Task    string
	Command Command
	File    File
}

type Command struct {
	Cmd string
}

type File struct {
	Source      string
	Destination string
	Mode        os.FileMode
}

// Dial will open a new SSH session to the specified host. This session can be used to perform
// actions or upload files.
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

func (c *Conn) Close() {
	defer c.sftpClient.Close()
	defer c.client.Close()
}

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

// ReadKeyFile will open an SSH keyfile and return an SSH Config for use
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
