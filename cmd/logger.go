package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

type logger struct {
	mutex  sync.Mutex
	list   []string
	file   *os.File
	writer *bufio.Writer
}

var log = logger{list: make([]string, 0)}

func getCurrentTimeString() string {
	now := time.Now()
	zoneInfo, _ := now.Zone()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%03d %s", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000000, zoneInfo)
}

func appendToList(category string, msg string) {
	logLine := getCurrentTimeString() + " : " + category + " : " + msg + "\n"
	log.mutex.Lock()
	log.list = append(log.list, logLine)
	log.mutex.Unlock()
}

func LogDebug(msg string) {
	appendToList("DEBUG", msg)
}

func LogInfo(msg string) {
	appendToList("INFO ", msg)
}

func LogInfoWithStdout(msg string) {
	LogInfo(msg)
	fmt.Println(msg)
}

func LogWarn(msg string) {
	appendToList("WARN ", msg)
}

func LogError(msg string) {
	appendToList("ERROR", msg)
}

func LogErrorWithStderr(msg string) {
	LogError(msg)
	fmt.Fprintln(os.Stderr, msg)
}

func LogFatal(msg string) {
	appendToList("FATAL", msg)
}

func processTasks() {

	// variable for write tasks
	var writeList []string

	// get write tasks
	log.mutex.Lock()
	if len(log.list) > 0 {
		writeList = log.list
		log.list = make([]string, 0)
	}
	log.mutex.Unlock()

	// process write tasks
	if writeList != nil {

		// write all
		for _, logLine := range writeList {
			_, err := log.writer.WriteString(logLine)
			if err != nil {
				LogError("LOG file WriteString error : " + err.Error())
				EndProgram = true
				return
			}
		}

		// flush
		err := log.writer.Flush()
		if err != nil {
			LogError("LOG file flush error : " + err.Error())
			EndProgram = true
			return
		}
	}
}

func LogBegin(fileName string) {

	// open log file
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		LogError("LOG file create error : " + err.Error())
		EndProgram = true
		return
	}

	// save file variables
	log.file = file
	log.writer = bufio.NewWriter(file)

	// infinite loop
	for !EndProgram {
		processTasks()
		time.Sleep(200 * time.Millisecond)
	}
}

func LogEnd() {

	// process last tasks
	processTasks()

	// close file
	if log.file != nil {
		log.file.Close()
	}
}
