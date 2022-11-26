package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/startdusk/tinygit"
)

var version = tinygit.Version()

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}
	fmt.Println(os.Args)
	cmd := os.Args[1]
	switch cmd {
	case "version":
		fmt.Println(version)
	case "init":
		var repo string
		if len(os.Args) > 2 {
			repo = os.Args[2]
		}
		if repo == "" {
			repo = "."
		}
		if err := initail(repo); err != nil {
			fmt.Println("can't init this repository")
			os.Exit(0)
		}
		dir, _ := os.Getwd()
		_, repo = filepath.Split(dir)
		fmt.Printf("initialized empty repository: %s", repo)
	case "help", "h":
		printHelp()
	default:
		fmt.Printf("Unsupported `%s` command\n", cmd)
	}
}

// Create directory for repo and initialize .git directory.
func initail(repo string) error {
	if err := os.MkdirAll(repo, os.ModePerm); err != nil {
		return err
	}
	tinygitPath := filepath.Join(repo, ".tinygit")
	if err := os.MkdirAll(tinygitPath, os.ModePerm); err != nil {
		return err
	}
	for _, name := range [3]string{"objects", "refs", "refs/heads"} {
		if err := os.MkdirAll(filepath.Join(tinygitPath, name), os.ModePerm); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(tinygitPath, "HEAD"), []byte("ref: refs/heads/master"), os.ModePerm); err != nil {
		return err
	}
	return nil
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
	fmt.Println(help)
}
