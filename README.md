# Efs2

Don't you wish you could configure a server as easily as creating a Docker image? Meet Efs2, A dead simple configuration management tool that is powered by stupid shell scripts.

Efs2 is an idea to combine the stupid shell scripts philosophy of [fss](https://github.com/brandonhilkert/fucking_shell_scripts) with the simplicity of a `Dockerfile`.

[![PkgGoDev](https://pkg.go.dev/badge/github.com/madflojo/efs2)](https://pkg.go.dev/github.com/madflojo/efs2) [![Go Report Card](https://goreportcard.com/badge/github.com/madflojo/efs2)](https://goreportcard.com/report/github.com/madflojo/efs2) [![Build Status](https://travis-ci.com/madflojo/efs2.svg?branch=master)](https://travis-ci.com/madflojo/efs2)

## Efs2 by Example: NGINX

Let's take a look at how easy it is to configure an NGINX server.

### The Efs2file

An `Efs2file` powers efs2's configuration; much like a `Dockerfile`, this file uses a simple set of instructions to configure our target servers.

```Dockerfile
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
```

The above `Efs2file` showcases how simple the Efs2 instructions are. Our NGINX server is configured with two simple instructions `RUN` and `PUT`.

The `RUN` instruction is simple; it executes whatever command you provide. The `PUT` instruction uploads files. That's it, that's all the instructions included with Efs2. Simple but effective.

### Remote Execution

Efs2 uses SSH to execute the instructions specified within the `Efs2file`. Just run the Efs2 command, followed by the target hosts.

```console
$ efs2 host1.example.com host2.example.com
```

#### Command Line Options

```
-v, --verbose   Enable verbose output
-f, --file=     Specify an alternative Efs2file (default: ./Efs2file)
-i, --key=      Specify an SSH Private key to use (default: ~/.ssh/id_rsa)
-p, --parallel  Execute tasks in parallel (default: false)
-d, --dryrun    Print tasks to be executed without actually executing any tasks
    --port=     Define an alternate SSH Port (default: 22)
-u, --user=     Remote host username (default: current user)
```

## Installation

Efs2 is simple to install, with the fastest method being to download one of our [binary releases](https://github.com/madflojo/efs2/releases).

It is also possible to install Efs2 with Go.

```console
go get -u github.com/madflojo/efs2
```

## Efs2file's In the wild

* [madflojo/masterless-salt-base](https://github.com/madflojo/masterless-salt-base/blob/master/Efs2file) - Installs and Configures a Masterless Salt Minion server

Add your examples above!

## TODO

* Recursive Directory support for `PUT`
* Password Authentication support
* Packaging for brew
* Templating for uploads

## Contributing

Thank you for your interest in helping develop Efs2. The time, skills, and perspectives you contribute to this project are valued.

Please reference our [Contributing Guide](CONTRIBUTING.md) for details.
