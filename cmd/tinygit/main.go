package main

import (
	"fmt"
	"os"

	"github.com/startdusk/tinygit"
)

var version = tinygit.Version()

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "version":
		fmt.Fprintln(os.Stdout, version)
	case "init":
	default:
		fmt.Fprintf(os.Stdout, "Unsupported `%s` command\n", cmd)
	}
}

func printHelp() {
	const help = `useage: tinygit [-v | --version] [-h | --help]
	        <command> [<args>]
These are common TinyGit commands used in various situations:

start a working area (see also: tinygit help tutorial)
   init      Create an empty TinyGit repository or reinitialize an existing one

work on the current change (see also: tinygit help everyday)
   add       Add file contents to the index
   mv        Move or rename a file, a directory, or a symlink
   rm        Remove files from the working tree and from the index
	`
	fmt.Fprintln(os.Stdout, help)
}
