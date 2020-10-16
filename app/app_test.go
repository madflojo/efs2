package app

import (
	"github.com/madflojo/efs2/config"
	"io/ioutil"
	"os"
	"testing"
)

func TestEfs2(t *testing.T) {

	cfg := config.Config{}
	cfg.KeyFile = "/go/src/github.com/madflojo/efs2/testdata/testkey"
	cfg.User = "test"
	cfg.Hosts = []string{"openssh-server:2222"}

	t.Run("Simple Non-Default Efs2file", func(t *testing.T) {
		b := []byte("RUN ls -la\nRUN CMD ls -la")
		f, _ := ioutil.TempFile("/tmp/", "testing.*.txt")
		defer os.Remove(f.Name())
		_, err := f.Write(b)
		if err != nil {
			t.Errorf("Error when creating test file - %s", err)
		}
		_ = f.Close()
		cfg.Efs2File = f.Name()

		err = Run(cfg)
		if err != nil {
			t.Errorf("Unexpected Error from App Execution - %s", err)
		}
	})

	t.Run("Non-existent Efs2file", func(t *testing.T) {
		cfg.Efs2File = "/tmp/doesntexist"
		err := Run(cfg)
		if err == nil {
			t.Errorf("Unexpected success from App Execution")
		}
	})

	t.Run("No Efs2file", func(t *testing.T) {
		cfg.Efs2File = ""
		err := Run(cfg)
		if err == nil {
			t.Errorf("Unexpected success from App Execution")
		}
	})

	t.Run("No Keyfile", func(t *testing.T) {
		b := []byte("RUN ls -la\nRUN CMD ls -la")
		f, _ := ioutil.TempFile("/tmp/", "testing.*.txt")
		defer os.Remove(f.Name())
		_, err := f.Write(b)
		if err != nil {
			t.Errorf("Error when creating test file - %s", err)
		}
		_ = f.Close()
		cfg.Efs2File = f.Name()
		cfg.KeyFile = ""

		err = Run(cfg)
		if err == nil {
			t.Errorf("Expected Error from App Execution got nil")
		}
	})
}
