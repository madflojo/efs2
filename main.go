/*
Don't you wish you could configure a server as easily as creating a Docker image? Meet Efs2, A dead simple configuration management tool that is powered by stupid shell scripts.

Efs2 is an idea to combine the stupid shell scripts philosophy of fss with the simplicity of a Dockerfile.

	# Install and Configure NGINX

	# Run apt-get update
	RUN apt-get update --fix-missing && apt-get -y upgrade

	# Install nginx
	RUN apt-get install nginx

	# Deploy Config files
	PUT nginx.conf /etc/nginx/nginx.conf 0644
	PUT example.com /etc/nginx/sites-available/example.com 0644

	# Create a Symlink
	RUN ln -s /etc/nginx/sites-available/example.com /etc/nginx/sites-enabled/example.com

	# Restart NGINX
	RUN systemctl restart nginx

*/
package main

import (
	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/jessevdk/go-flags"
	"github.com/madflojo/efs2/app"
	"github.com/madflojo/efs2/config"
	"os"
)

// options are command-line options that are provided by the user.
type options struct {
	Verbose  bool   `short:"v" long:"verbose" description:"Enable verbose output"`
	Quiet    bool   `short:"q" long:"quiet" description:"Silence output"`
	Efs2File string `short:"f" long:"file" description:"Specify an alternative Efs2File" default:"./Efs2file"`
	KeyFile  string `short:"i" long:"key" description:"Specify a SSH Private key to use (default: ~/.ssh/id_rsa)"`
	Parallel bool   `short:"p" long:"parallel" description:"Execute tasks across multiple hosts in parallel"`
	DryRun   bool   `short:"d" long:"dryrun" description:"Print tasks to be executed without actually executing any tasks"`
	Port     string `long:"port" description:"Define an alternate SSH Port" default:"22"`
	User     string `short:"u" long:"user" description:"Remote host username (default: current user)"`
	Pass     bool   `long:"passwd" description:"Ask for a password to use for host authentication"`
}

// main runs the command line parsing and validations. This function will also start the application logic execution.
func main() {
	// Parse command line arguments
	var opts options
	args, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	// Convert to internal config
	cfg := config.New()
	cfg.Verbose = opts.Verbose
	cfg.Quiet = opts.Quiet
	if opts.Efs2File != "" {
		cfg.Efs2File = opts.Efs2File
	}
	if opts.KeyFile != "" {
		cfg.KeyFile = opts.KeyFile
	}
	if opts.User != "" {
		cfg.User = opts.User
	}
	cfg.Parallel = opts.Parallel
	cfg.DryRun = opts.DryRun
	cfg.Port = opts.Port
	cfg.Hosts = args

	if opts.Pass {
		color.White("Enter Private Key Passphrase: ")
		p, err := gopass.GetPasswd()
		if err != nil {
			color.Red("Unable to obtain SSH Password: %s", err)
		}
		cfg.Password = string(p)
	}

	// Run the App
	err = app.Run(cfg)
	if err != nil {
		if !cfg.Quiet {
			color.Red("Error executing: %s", err)
		}
		os.Exit(1)
	}
	if !cfg.Quiet {
		color.Green("Execution completed successfully")
	}
}
