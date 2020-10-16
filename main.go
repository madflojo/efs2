/*
Don't you wish you could configure a server as easily as creating a Docker image? Meet Efs2, A dead simple configuration management tool that is powered by stupid shell scripts.

Efs2 is an idea to combine the stupid shell scripts philosophy of fss with the simplicity of a Dockerfile.
*/
package main

import (
	"github.com/madflojo/efs2/cmd/cli"
)

func main() {
	cli.Exec()
}
