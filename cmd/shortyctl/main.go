package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		usage()
		return
	}

	switch os.Args[1] {
	case "seed":
		switch os.Args[2] {
		case "generate":
			generateUniqueIds()
		case "import":
			if err := importUniqueIds(); err != nil {
				fmt.Println(err)
			}
		default:
			usage()
		}
	case "login":
		if err := loginCmd(os.Args[2]); err != nil {
			fmt.Println(err)
		}
	case "list":
		listCmd()
	case "create":
		if err := createCmd(os.Args[2]); err != nil {
			fmt.Println(err)
		}
	default:
		usage()
	}
}

func usage() {
	fmt.Println(`shortyctl <command> [<options>]

Commands
  seed <subcommand>
  add <options>

Seed Subcommands
  generate 
  import [--redis <url>] [--file <path>]
  login  <url>
  create <url>
  list

Options
  --redis <url>
  --file  <path>
  --host  <url>
  --url   <url>
	`)
}
