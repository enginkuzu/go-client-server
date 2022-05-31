package main

import (
	"github.com/enginkuzu/go-client-server/cmd"
	"net"
)

var conn net.Conn

func clientSocketBegin(id uint32, server string) {

	// set client id
	clientId = id

	// connect to server
	var err error
	conn, err = net.Dial("tcp", server+":"+cmd.ServerListenPortNumber)
	if err != nil {
		cmd.LogErrorWithStderr("Net dial to " + server + ":" + cmd.ServerListenPortNumber + " error : " + err.Error())
		cmd.EndProgram = true
		return
	}
	cmd.LogInfoWithStdout("Client connected from " + conn.LocalAddr().String() + " to " + conn.RemoteAddr().String())

	// infinite loop
	for !cmd.EndProgram {

		// read message from queue
		message, more := <-diskToSocketChannelQueue
		if !more {
			cmd.LogInfoWithStdout("Queue closed")
			cmd.EndProgram = true
			continue
		}

		// send message to server
		_, err := conn.Write(message)
		if err != nil {
			cmd.LogError("Socket write error : " + err.Error())
			cmd.EndProgram = true
			continue
		}
	}
}

func clientSocketEnd() {
	// close socket
	if conn != nil {
		conn.Close()
	}
}
