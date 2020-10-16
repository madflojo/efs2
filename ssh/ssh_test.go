package ssh

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestReadKeyFile(t *testing.T) {

	t.Run("No File", func(t *testing.T) {
		_, err := ReadKeyFile("/nothing", []byte(""))
		if err == nil {
			t.Errorf("Did not receive an error with none existant key file")
		}
	})

	t.Run("Empty File", func(t *testing.T) {
		defer os.Remove("/tmp/empty-file-test.key")
		fh, err := os.Create("/tmp/empty-file-test.key")
		if err != nil {
			t.Errorf("Error creating empty file")
		}
		defer fh.Close()

		_, err = ReadKeyFile("/tmp/empty-file-test.key", []byte(""))
		if err == nil {
			t.Errorf("Did not receive an error with empty key file")
		}
	})

	t.Run("Regular Key", func(t *testing.T) {
		_, err := ReadKeyFile("/go/src/github.com/madflojo/efs2/testdata/testkey", []byte(""))
		if err != nil {
			t.Errorf("Unexpected error reading key file - %s", err)
		}

	})

	t.Run("Encrypted Key", func(t *testing.T) {
		_, err := ReadKeyFile("/go/src/github.com/madflojo/efs2/testdata/testkey-passphrase", []byte("testing"))
		if err != nil {
			t.Errorf("Unexpected error reading key file - %s", err)
		}
	})
}

func TestBadConnections(t *testing.T) {
	t.Run("Invalid Host", func(t *testing.T) {
		c, err := ReadKeyFile("/go/src/github.com/madflojo/efs2/testdata/testkey-passphrase", []byte("testing"))
		if err != nil {
			t.Errorf("Unexpected error reading key file - %s", err)
		}
		c.User = "test"

		_, err = Dial(c)
		if err == nil {
			t.Errorf("Unexpected success dialing an invalid host")
		}
	})

	t.Run("No User", func(t *testing.T) {
		c, err := ReadKeyFile("/go/src/github.com/madflojo/efs2/testdata/testkey-passphrase", []byte("testing"))
		if err != nil {
			t.Errorf("Unexpected error reading key file - %s", err)
		}
		c.Host = "openssh-server:2222"

		_, err = Dial(c)
		if err == nil {
			t.Errorf("Unexpected success dialing an invalid host")
		}

	})
}

func TestDial(t *testing.T) {
	c, err := ReadKeyFile("/go/src/github.com/madflojo/efs2/testdata/testkey", []byte(""))
	if err != nil {
		t.Errorf("Unexpected error reading key file - %s", err)
	}

	c.Host = "openssh-server:2222"
	c.User = "test"

	s, err := Dial(c)
	if err != nil {
		t.Errorf("Unexpected error dialing host - %s", err)
		t.FailNow()
	}
	defer s.Close()

	t.Run("Upload a script", func(t *testing.T) {
		script := "/tmp/test-script.sh"
		defer os.Remove(script)
		b := []byte("#!/bin/bash\necho great success")
		err := ioutil.WriteFile(script, b, 0644)
		if err != nil {
			t.Errorf("Error creating test script")
		}

		f := File{
			Source:      "/tmp/test-script.sh",
			Destination: "/tmp/test-script.sh",
			Mode:        os.FileMode(int(0755)),
		}

		err = s.Put(f)
		if err != nil {
			t.Errorf("Failed to put file - %s", err)
		}

	})

	t.Run("Upload missing file", func(t *testing.T) {
		f := File{
			Source:      "/tmp/notafile",
			Destination: "/tmp/afile",
			Mode:        os.FileMode(int(0644)),
		}

		err = s.Put(f)
		if err == nil {
			t.Errorf("Expected file upload failure, got success")
		}
	})

	t.Run("Upload file to nowhere", func(t *testing.T) {
		script := "/tmp/test-script.sh"
		defer os.Remove(script)
		b := []byte("#!/bin/bash\necho great success")
		err := ioutil.WriteFile(script, b, 0644)
		if err != nil {
			t.Errorf("Error creating test script")
		}

		f := File{
			Source:      "/tmp/test-script.sh",
			Destination: "/path/does/not/exist",
			Mode:        os.FileMode(int(0644)),
		}

		err = s.Put(f)
		if err == nil {
			t.Errorf("Expected file upload failure, got success")
		}
	})

	t.Run("Execute that script", func(t *testing.T) {
		cmd := Command{
			Cmd: "/tmp/test-script.sh && rm /tmp/test-script.sh",
		}
		_, err := s.Run(cmd)
		if err != nil {
			t.Errorf("Failed to execute command - %s", err)
		}
	})

	t.Run("Execute invalid command", func(t *testing.T) {
		cmd := Command{
			Cmd: "/tmp/test-script.sh",
		}
		_, err := s.Run(cmd)
		if err == nil {
			t.Errorf("Expected failure got success")
		}
	})
}
