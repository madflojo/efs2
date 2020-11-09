Don't you wish you could configure a server as easily as creating a Docker image? Meet Efs2, A dead simple configuration management tool that is powered by stupid shell scripts.

Efs2 is an idea to combine the stupid shell scripts philosophy of [fss](https://github.com/brandonhilkert/fucking_shell_scripts) with the simplicity of a `Dockerfile`.

## Getting Started

Let's take a look at how easy it is to use Efs2 to configure NGINX on Ubuntu.

### Installation

Efs2 is simple to install, with the fastest method being to download one of our [binary releases](https://github.com/madflojo/efs2/releases).

It is also possible to install Efs2 with Go (requires v1.14+).

```console
go get -u github.com/madflojo/efs2
```

Once installed, we can start defining our steps to setup NGINX.

### The Efs2file

An `Efs2file` powers Efs2's configuration; much like a `Dockerfile`, this file uses a simple set of instructions to configure our target servers.

```Dockerfile
# Install and Configure NGINX

# Run apt-get update
RUN apt-get update --fix-missing && apt-get -y upgrade

# Install nginx
RUN apt-get install nginx

# Deploy Config files
PUT example.com.conf /etc/nginx/sites-available/example.com.conf 0644

# Create a Symlink
RUN ln -s /etc/nginx/sites-available/example.com.conf /etc/nginx/sites-enabled/example.com.conf

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

#### Command-Line Options

Efs2 offers several additional options, such as parallel execution and various authentication methods.

```console
  -v, --verbose   Enable verbose output
  -q, --quiet     Silence output
  -f, --file=     Specify an alternative Efs2File (default: ./Efs2file)
  -i, --key=      Specify an SSH Private key to use (default: ~/.ssh/id_rsa)
  -p, --parallel  Execute tasks in parallel
  -d, --dryrun    Print tasks to be executed without actually executing any tasks
      --port=     Define an alternate SSH Port (default: 22)
  -u, --user=     Remote host username (default: current user)
      --passwd    Ask for a password to use for authentication
```

## Call to Action

Efs2 is a small project to fit the fine line between complex configuration management and simple shell scripts.  We are always looking for users to share their stories and contribute to our [examples repository](https://github.com/madflojo/efs2-examples).

If you like Efs2, please tell others about it by sharing this project on the social media site of your choice
