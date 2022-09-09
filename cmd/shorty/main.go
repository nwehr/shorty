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
	case "gen-seed":
		generateUniqueIds()
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
  seed   <subcommand>
  add    <options>
  login  <url>
  create <url>
  list

Seed Subcommands
  generate 
  import [--redis <url>] [--file <path>]

Options
  --redis <url>
  --file  <path>
  --host  <url>
  --url   <url>
	`)
}
