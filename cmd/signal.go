package cmd

import (
	"os"
	"os/signal"
	"syscall"
)

func SignalHandlerSetup() {

	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-chanSignal
		LogInfoWithStdout("\"" + sig.String() + "\" signal received")
		EndProgram = true
	}()
}
