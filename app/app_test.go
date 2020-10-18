package app

import (
	"github.com/madflojo/efs2/config"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestEfs2(t *testing.T) {

	cfg := config.Config{}
  cfg.Quiet = true
  cfg.Verbose = true
	cfg.KeyFile = "/go/src/github.com/madflojo/efs2/testdata/testkey"
	cfg.User = "test"
	cfg.Hosts = []string{"openssh-server:2222"}

	t.Run("Simple Efs2file", func(t *testing.T) {
		f, _ := ioutil.TempFile("/tmp/", "testing.*.txt")
		b := []byte("RUN ls -la\nRUN CMD ls -la\nPUT " + f.Name() + " /tmp/somefile 0644")
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

	t.Run("Simple Dryrun Efs2file", func(t *testing.T) {
		f, _ := ioutil.TempFile("/tmp/", "testing.*.txt")
		b := []byte("RUN ls -la\nRUN CMD ls -la\nPUT " + f.Name() + " /tmp/somefile 0644")
		defer os.Remove(f.Name())
		_, err := f.Write(b)
		if err != nil {
			t.Errorf("Error when creating test file - %s", err)
		}
		_ = f.Close()
		cfg.Efs2File = f.Name()
		cfg.DryRun = true

		err = Run(cfg)
		if err != nil {
			t.Errorf("Unexpected Error from App Execution - %s", err)
		}
	})
	cfg.DryRun = false

	t.Run("Broken Command", func(t *testing.T) {
		f, _ := ioutil.TempFile("/tmp/", "testing.*.txt")
		b := []byte("RUN thisaintnocommand")
		defer os.Remove(f.Name())
		_, err := f.Write(b)
		if err != nil {
			t.Errorf("Error when creating test file - %s", err)
		}
		_ = f.Close()
		cfg.Efs2File = f.Name()

		err = Run(cfg)
		if err == nil {
			t.Errorf("Unexpected success from App Execution")
		}
	})

	t.Run("Broken Put", func(t *testing.T) {
		f, _ := ioutil.TempFile("/tmp/", "testing.*.txt")
		b := []byte("PUT /path/to/no/where /tmp/here 0644")
		defer os.Remove(f.Name())
		_, err := f.Write(b)
		if err != nil {
			t.Errorf("Error when creating test file - %s", err)
		}
		_ = f.Close()
		cfg.Efs2File = f.Name()

		err = Run(cfg)
		if err == nil {
			t.Errorf("Unexpected success from App Execution")
		}
	})

	t.Run("Invalid Host", func(t *testing.T) {
		f, _ := ioutil.TempFile("/tmp/", "testing.*.txt")
		b := []byte("PUT /path/to/no/where /tmp/here 0644")
		defer os.Remove(f.Name())
		_, err := f.Write(b)
		if err != nil {
			t.Errorf("Error when creating test file - %s", err)
		}
		_ = f.Close()
		cfg.Efs2File = f.Name()
		cfg.Hosts = []string{"nope:99090"}

		err = Run(cfg)
		if err == nil {
			t.Errorf("Unexpected success from App Execution")
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

func TestHostFixer(t *testing.T) {
	t.Run("Already Has Ports", func(t *testing.T) {
		hosts := []string{"testing:22", "nope:22"}
		hh := fixUpHosts(hosts, "22")
		if !reflect.DeepEqual(hosts, hh) {
			t.Errorf("Unexpected values returned should be equal %+v != %+v", hosts, hh)
		}
	})

	t.Run("Some have Ports", func(t *testing.T) {
		good := []string{"testing:22", "nope:22"}
		bad := []string{"testing", "nope:22"}
		hh := fixUpHosts(bad, "22")
		if !reflect.DeepEqual(good, hh) {
			t.Errorf("Unexpected values returned should be equal %+v != %+v", good, hh)
		}
	})

	t.Run("Some have Ports Default Port", func(t *testing.T) {
		good := []string{"testing:22", "nope:22"}
		bad := []string{"testing", "nope:22"}
		hh := fixUpHosts(bad, "")
		if !reflect.DeepEqual(good, hh) {
			t.Errorf("Unexpected values returned should be equal %+v != %+v", good, hh)
		}
	})
}
