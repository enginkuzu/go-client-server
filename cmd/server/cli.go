package main

import (
	"fmt"
	"os"
	"runtime"
)

const ParamVersion string = "--version"
const ParamHelp string = "--help"

const OutputHelp string = `Usage: go-server [options]

Starts TCP/IP server application.

Options:
--version   output version information and exit
--help      display this help and exit`

func processCliParameters() (end bool) {

	args := os.Args

	if len(args) == 1 {
		return false
	} else if len(args) == 2 {
		if args[1] == ParamVersion {
			fmt.Println(ServerAppNameAndVersion + " (build with " + runtime.Version() + ")")
			return true
		} else if args[1] == ParamHelp {
			fmt.Println(OutputHelp)
			return true
		}
	}

	fmt.Println("go-server: unrecognized option '" + args[1] + "'")
	return true
}
