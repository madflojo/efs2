# Effing Scripts 2

Wish you could go back and just use stupid scripts to setup and configure your servers? Well now you can! **Effing Scripts 2** is a common sense configuration management tool that is loosly inspired by [fss](https://github.com/brandonhilkert/fucking_shell_scripts).

## Getting started

Getting started with Effing Scripts 2 is as easy as creating an `Efs2file`. This file will define what scripts and commands you wish to run on the target hosts. Let's take a look at an example file.

```
# Setup a simple mailserver

# Run apt-get update
RUN CMD apt-get update --fix-missing && apt-get -y upgrade

# PUT the main.cf file to the remote host
PUT files/main.cf /etc/postfix/main.cf 0644

# Copy the setup_postfix.sh script to the remote host and then execute it
RUN SCRIPT setup_postfix.sh

# Execute a single line command on the remote host
RUN CMD ps -elf | grep -q postfix
```

### Available `Efs2file` Instructions

The `Efs2file` is a simple instruction file, below are a list of available instructions and what they do.

  * `PUT` - Copy the specified file or directory to the remote host
  * `RUN` - Execute a script or command on the remote host 
    * `SCRIPT` - Copy the specified script to the remote host and execute it
    * `CMD` - Login to remote host and execute the specified command

### Executing `efs2`

Once the `Efs2file` is defined you can execute it against any number of target hosts.

```sh
$ efs2 host1.example.com host2.example.com
```

### Options

```
Usage:
  app [OPTIONS]

Application Options:
  -v, --verbose   Enable verbose output
  -f, --file=     Specify an alternative Efs2file (default: ./Efs2file)
                  (default: ./Efs2file)
  -i, --key=      Specify an SSH Private key to use (default: ~/.ssh/id_rsa)
  -p, --parallel  Execute tasks in parallel (default: false)
  -d, --dryrun    Print tasks to be executed without actually executing any
                  tasks
      --port=     Define an alternate SSH Port (default: 22) (default: 22)

Help Options:
  -h, --help      Show this help message
```
