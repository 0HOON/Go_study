package main

import (
	"os"

	"example.com/gocoin/cli"
)

func main() {
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()

}
