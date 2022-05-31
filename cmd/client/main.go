package main

import (
	"fmt"
	"github.com/enginkuzu/go-client-server/cmd"
	"runtime"
	"time"
)

func main() {

	// process parameters
	id, server, end := processCliParameters()
	if end {
		return
	}

	// startup log messages
	cmd.LogInfoWithStdout("Starting " + ClientAppNameAndVersion + " (build with " + runtime.Version() + ")")
	cmd.LogInfoWithStdout(fmt.Sprintf("With Parameters : id=%v, server=%q", id, server))

	// signal handler
	cmd.SignalHandlerSetup()

	// logger
	go cmd.LogBegin(LogFileName)

	// client socket
	go clientSocketBegin(id, server)

	// disk scanner
	go diskScannerBegin()

	// wait for application shutdown
	for !cmd.EndProgram {
		time.Sleep(200 * time.Millisecond)
	}

	// terminate all resources and threads
	cmd.LogInfoWithStdout("App ending ...")
	clientSocketEnd()
	cmd.LogEnd()
	fmt.Println("App end")
}
