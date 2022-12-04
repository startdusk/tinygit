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
		tinygit.PrintHelp()
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
		if err := tinygit.Initail(repo); err != nil {
			fmt.Println("can't init this repository")
			os.Exit(0)
		}
		if repo == "." {
			dir, _ := os.Getwd()
			_, repo = filepath.Split(dir)
		}
		_ = repo // TODO: print something ...
	case "add":
		if len(os.Args) != 3 {
			// TODO: print somthing...
			os.Exit(1)
		}
		path := os.Args[2]
		tinygit.Add(path)
	case "help", "h":
		tinygit.PrintHelp()
	default:
		fmt.Printf("Unsupported `%s` command\n", cmd)
	}
}
