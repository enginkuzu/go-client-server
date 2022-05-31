package main

import (
	"bufio"
	"fmt"
	"github.com/enginkuzu/go-client-server/cmd"
	"os"
	"sync"
	"time"
)

type dump struct {
	mutex  sync.Mutex
	list   []string
	file   *os.File
	writer *bufio.Writer
}

var dmp = dump{list: make([]string, 0)}

func csvAppendHeader() {
	logLine := "Message Size,Protocol,Client Id,File Size,File Permissions,File Path Length,File Path,Checksum\n"
	dmp.mutex.Lock()
	dmp.list = append(dmp.list, logLine)
	dmp.mutex.Unlock()
}

func csvAppend(messageSize uint64, protocol uint8, clientId uint32, fileSize uint64, filePermissions uint32, filePathLength uint32, filePath string, checksum uint32) {
	logLine := fmt.Sprintf("%d,%d,%d,%d,%d,%d,%s,%d\n", messageSize, protocol, clientId, fileSize, filePermissions, filePathLength, filePath, checksum)
	dmp.mutex.Lock()
	dmp.list = append(dmp.list, logLine)
	dmp.mutex.Unlock()
}

func processTasks() {

	// variable for write tasks
	var writeList []string

	// get write tasks
	dmp.mutex.Lock()
	if len(dmp.list) > 0 {
		writeList = dmp.list
		dmp.list = make([]string, 0)
	}
	dmp.mutex.Unlock()

	// process write tasks
	if writeList != nil {

		// write all
		for _, logLine := range writeList {
			_, err := dmp.writer.WriteString(logLine)
			if err != nil {
				cmd.LogError("CSV file WriteString error : " + err.Error())
				cmd.EndProgram = true
				return
			}
		}

		// flush
		err := dmp.writer.Flush()
		if err != nil {
			cmd.LogError("CSV file flush error : " + err.Error())
			cmd.EndProgram = true
			return
		}
	}
}

func csvBegin(fileName string) {

	// append header
	csvAppendHeader()

	// open log file
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		cmd.LogError("CSV file create error : " + err.Error())
		cmd.EndProgram = true
		return
	}

	// save file variables
	dmp.file = file
	dmp.writer = bufio.NewWriter(file)

	// infinite loop
	for !cmd.EndProgram {
		processTasks()
		time.Sleep(200 * time.Millisecond)
	}
}

func csvEnd() {

	// process last tasks
	processTasks()

	// close log file
	if dmp.file != nil {
		dmp.file.Close()
	}
}
