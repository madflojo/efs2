# Effing Shell Scripts 2

Are you hating Puppet right now? Wish you could dump these complex tools like CFEngine, Chef, Puppet, Ansible and Salt? Want to go back to using stupid scripts!? Now you can!

**Effing Shell Scripts 2** is a common sense remote command execution tool inspired by [fss](https://github.com/brandonhilkert/fucking_shell_scripts) and written in **Go**.

Let's get started!

## Installation

Installation with Go is as easy as running `go get`.

```sh
go get -u github.com/madflojo/efs2
```

Binary releases are [available](https://github.com/madflojo/efs2/releases).

## Creating an `Efs2file`

Rather than using a `yaml` based desired state structure. Effing Shell Scripts 2 tries to keep things simple. Just create a `Efs2file` and define what files to copy and scripts/commands to run.

Let's take a look at an example file.

```
# Setup a simple mailserver

# Run apt-get update
RUN CMD apt-get update --fix-missing && apt-get -y upgrade

# PUT the main.cf file to the remote host
PUT files/main.cf /etc/postfix/main.cf 0644

# Copy the setup_postfix.sh script to the remote host and then execute it
RUN SCRIPT setup_postfix.sh

# Execute a one liner command on the remote host
RUN CMD ps -elf | grep -q postfix
```

The order of this file, is the order these instructions are executed. Allowing you to skip the hassle of the complex dependencies other tools need.

### Available `Efs2file` Instructions

Effing Shell Scripts 2 is a simple tool for simple tasks. As such, the `Efs2file` has only two instructions; `PUT` and `RUN`.

- `PUT` - Copy the specified file to the remote host
- `RUN` - Execute a script or command on the remote host
  - `SCRIPT` - Copy the specified script to the remote host and execute it
  - `CMD` - Execute the specified command on the remote host

## Executing `efs2`

Once defined, the `Efs2file` can be executed against any number of target hosts.

```sh
$ efs2 host1.example.com host2.example.com
```

**Available command line options:**
```
-v, --verbose   Enable verbose output
-f, --file=     Specify an alternative Efs2file (default: ./Efs2file)
-i, --key=      Specify an SSH Private key to use (default: ~/.ssh/id_rsa)
-p, --parallel  Execute tasks in parallel (default: false)
-d, --dryrun    Print tasks to be executed without actually executing any tasks
    --port=     Define an alternate SSH Port (default: 22)
-u, --user=     Remote host username (default: current user)
```

## TODO

* Directory support for `PUT`
* Example `Efs2file`s
* Packaging for common OS distribution channels
