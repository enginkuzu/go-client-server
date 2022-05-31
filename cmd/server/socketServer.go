package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/enginkuzu/go-client-server/cmd"
	"hash/crc32"
	"net"
	"time"
)

func readFromSocket(conn net.Conn, limit uint64) ([]byte, error) {
	var readed uint64 = 0
	buffer := make([]byte, limit)
	for {
		n, err := conn.Read(buffer[readed:limit])
		if err != nil {
			return nil, err
		}
		readed += uint64(n)
		if readed == limit {
			break
		}
	}
	return buffer, nil
}

// |--------------------------------------------------------------------------------------------------------------|
// | Message Size | Protocol | Client ID | File Size | File Permissions | File Path Length | File Path | Checksum |
// |--------------------------------------------------------------------------------------------------------------|
// | uint64       | uint8    | uint32    | uint64    | uint32           | uint32           | string    | uint32   |
// |--------------------------------------------------------------------------------------------------------------|

func handleClientConnection(conn net.Conn) (fileCounter uint32) {

	// release resources
	defer conn.Close()

	// infinite loop
	for !cmd.EndProgram {

		// read first 8 byte
		buffer1, err1 := readFromSocket(conn, 8)
		if err1 != nil {
			return
		}

		// variable of message size
		messageSize := binary.BigEndian.Uint64(buffer1)

		// read the rest of the message
		buffer2, err2 := readFromSocket(conn, messageSize-8)
		if err2 != nil {
			return
		}

		// merge buffers
		message := append(buffer1, buffer2...)

		// get variables
		var protocol = message[8]
		var clientId = binary.BigEndian.Uint32(message[9:])
		var fileSize = binary.BigEndian.Uint64(message[13:])
		var filePermissions = binary.BigEndian.Uint32(message[21:])
		var filePathLength = binary.BigEndian.Uint32(message[25:])
		var filePath = bytes.NewBuffer(message[29 : 29+filePathLength]).String()
		var checksum = binary.BigEndian.Uint32(message[29+filePathLength:])

		// calculate CRC-32 IEEE polynomial checksum
		crcTable := crc32.MakeTable(crc32.IEEE)
		calculatedChecksum := crc32.Checksum(message[:29+filePathLength], crcTable)

		// check message content
		if protocol != cmd.Protocol {
			cmd.LogError(fmt.Sprintf("Message protocol variable mismatch! (connection=%s, clientId=%d)\n", conn.RemoteAddr().String(), clientId))
			return
		} else if checksum != calculatedChecksum {
			cmd.LogError(fmt.Sprintf("Message checksum variable mismatch! (connection=%s, clientId=%d)\n", conn.RemoteAddr().String(), clientId))
			return
		}

		// append to csv file
		csvAppend(messageSize, protocol, clientId, fileSize, filePermissions, filePathLength, filePath, checksum)

		// counter
		fileCounter++
	}

	return fileCounter
}

var listener net.Listener

func serverSocketBegin() {

	// bind port
	var err error
	listener, err = net.Listen("tcp", ServerListenIpAddress+":"+cmd.ServerListenPortNumber)
	if err != nil {
		cmd.LogErrorWithStderr("Port bind error : " + err.Error())
		cmd.EndProgram = true
		return
	}
	cmd.LogInfoWithStdout("Bind on " + ServerListenIpAddress + ":" + cmd.ServerListenPortNumber + " success")

	// infinite loop
	for !cmd.EndProgram {

		// accept new connections
		conn, err := listener.Accept()
		if cmd.EndProgram {
			continue
		}
		if err != nil {
			cmd.LogErrorWithStderr("Port accept error : " + err.Error())
			cmd.EndProgram = true
			continue
		}
		cmd.LogInfoWithStdout("New connection from " + conn.RemoteAddr().String() + " to " + conn.LocalAddr().String())

		// new thread for new connection
		go func() {

			// start operations
			t1 := time.Now()

			// handle connection
			fileCounter := handleClientConnection(conn)

			// end operations
			t2 := time.Now()
			durationInSeconds := t2.Sub(t1).Milliseconds() / 1000
			cmd.LogInfoWithStdout(fmt.Sprintf("Client send %v messages in %v seconds", fileCounter, durationInSeconds))
			cmd.LogInfoWithStdout("Connection closed from " + conn.RemoteAddr().String())
		}()
	}
}

func serverSocketEnd() {
	// close listener
	if listener != nil {
		listener.Close()
	}
}
