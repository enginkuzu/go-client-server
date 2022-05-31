package main

import (
	"fmt"
	"github.com/enginkuzu/go-client-server/cmd"
	"runtime"
	"time"
)

func main() {

	// process parameters
	end := processCliParameters()
	if end {
		return
	}

	// startup log message
	cmd.LogInfoWithStdout("Starting " + ServerAppNameAndVersion + " (build with " + runtime.Version() + ")")

	// signal handler
	cmd.SignalHandlerSetup()

	// logger
	go cmd.LogBegin(LogFileName)

	// csv file management
	go csvBegin(CsvFileName)

	// server socket
	go serverSocketBegin()

	// wait for application shutdown
	for !cmd.EndProgram {
		time.Sleep(200 * time.Millisecond)
	}

	// terminate all resources and threads
	cmd.LogInfoWithStdout("App ending ...")
	serverSocketEnd()
	csvEnd()
	cmd.LogEnd()
	fmt.Println("App end")
}
