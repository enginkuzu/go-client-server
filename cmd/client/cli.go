package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

const ParamId string = "--id"
const ParamServer string = "--server"
const ParamVersion string = "--version"
const ParamHelp string = "--help"

const OutputHelp string = `Usage: go-client [options]

Starts TCP/IP client application.

Options:
--id        specifies client ID (uint32)
--server    specifies server ipv4/ipv6 address (string)
--version   output version information and exit
--help      display this help and exit`

func processCliParameters() (id uint32, server string, end bool) {

	args := os.Args

	if len(args) == 1 {
		fmt.Println("go-client: mandatory parameters are missing, --help for more information")
		return 0, "", true
	} else if len(args) == 2 {
		if args[1] == ParamVersion {
			fmt.Println(ClientAppNameAndVersion + " (build with " + runtime.Version() + ")")
			return 0, "", true
		} else if args[1] == ParamHelp {
			fmt.Println(OutputHelp)
			return 0, "", true
		}
	} else if len(args) == 5 {
		if (args[1] == ParamServer && args[3] == ParamId) || (args[1] == ParamId && args[3] == ParamServer) {
			server := func() string {
				if args[1] == ParamServer {
					return args[2]
				} else {
					return args[4]
				}
			}()
			id := func() string {
				if args[1] == ParamId {
					return args[2]
				} else {
					return args[4]
				}
			}()
			i, err := strconv.ParseInt(id, 10, 33)
			if i < 0 || err != nil {
				fmt.Println("go-client: id parameter incompatible for type uint32")
				return 0, "", true
			}
			return uint32(i), server, false
		}
	}

	fmt.Println("go-client: check parameters or you can use --help parameter for more information")
	return 0, "", true
}
