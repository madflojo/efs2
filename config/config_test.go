package config

import (
	"os/user"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("Check Defaults", func(t *testing.T) {
		cfg := New()
		if cfg.Efs2File != "./Efs2file" && cfg.Port != "22" {
			t.Errorf("Defaults are not set")
		}
	})

	t.Run("Check User Details", func(t *testing.T) {
		usr, err := user.Current()
		if err != nil {
			t.Errorf("Unable to determine current user - %s", err)
		}

		username, homedir, err := UserDetails("")
		if err != nil {
			t.Errorf("Unable to grab user details - %s", err)
		}

		if username != usr.Username && homedir != usr.HomeDir {
			t.Errorf("User details are not accurate")
		}
	})

	t.Run("Check Specified User Details", func(t *testing.T) {
		usr, err := user.Current()
		if err != nil {
			t.Errorf("Unable to determine current user - %s", err)
		}

		username, homedir, err := UserDetails(usr.Username)
		if err != nil {
			t.Errorf("Unable to grab user details - %s", err)
		}

		if username != usr.Username && homedir != usr.HomeDir {
			t.Errorf("User details are not accurate")
		}
	})
}
