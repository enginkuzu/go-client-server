package main

import (
	"encoding/binary"
	"fmt"
	"github.com/enginkuzu/go-client-server/cmd"
	"hash/crc32"
	"io/fs"
	"os"
	"runtime"
	"time"
)

var fileCounter, fileErrorCounter uint32 = 0, 0
var folderCounter, folderErrorCounter uint32 = 0, 0

// |--------------------------------------------------------------------------------------------------------------|
// | Message Size | Protocol | Client ID | File Size | File Permissions | File Path Length | File Path | Checksum |
// |--------------------------------------------------------------------------------------------------------------|
// | uint64       | uint8    | uint32    | uint64    | uint32           | uint32           | string    | uint32   |
// |--------------------------------------------------------------------------------------------------------------|

func prepareMessage(currentFolderPath string, info fs.FileInfo) []byte {

	// prepare variables
	fileSize := uint64(info.Size())
	filePermissions := uint32(info.Mode())
	filePath := currentFolderPath + info.Name()
	filePathLength := uint32(len(filePath))
	messageSize := uint64(8 + 1 + 4 + 8 + 4 + 4 + filePathLength + 4)

	// make byte slice
	message := make([]byte, messageSize)

	// fill byte slice
	binary.BigEndian.PutUint64(message[0:], messageSize)
	message[8] = cmd.Protocol
	binary.BigEndian.PutUint32(message[9:], clientId)
	binary.BigEndian.PutUint64(message[13:], fileSize)
	binary.BigEndian.PutUint32(message[21:], filePermissions)
	binary.BigEndian.PutUint32(message[25:], filePathLength)
	copy(message[29:], filePath)

	// calculate CRC-32 IEEE polynomial checksum
	crcTable := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum(message[:29+filePathLength], crcTable)

	// fill byte slice
	binary.BigEndian.PutUint32(message[29+filePathLength:], checksum)

	// return filled byte slice
	return message
}

func folderScan(currentFolderPath string) {

	// read directory contents
	dirEntries, err := os.ReadDir(currentFolderPath)
	if err != nil {
		cmd.LogError("Folder read error : " + err.Error())
		folderErrorCounter++
		return
	}
	folderCounter++

	// process files
	for _, entry := range dirEntries {
		if cmd.EndProgram {
			return
		}
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				cmd.LogError("File info error : " + err.Error())
				fileErrorCounter++
				continue
			}
			message := prepareMessage(currentFolderPath, info)
			// send to channel, current thread is blocked when queue is full
			diskToSocketChannelQueue <- message
			fileCounter++
		}
	}

	// scan sub folders
	for _, entry := range dirEntries {
		if cmd.EndProgram {
			return
		}
		if entry.IsDir() {
			subFolder := currentFolderPath + entry.Name() + string(os.PathSeparator)
			folderScan(subFolder)
		}
	}
}

func diskScannerBegin() {

	// start operations
	t1 := time.Now()
	cmd.LogInfoWithStdout("Disk scan starting ...")

	// scan root folders
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		folderScan(string(os.PathSeparator))
	} else if runtime.GOOS == "windows" {
		for r := 'A'; r <= 'Z'; r++ {
			folderScan(string(r) + ":" + string(os.PathSeparator))
		}
	} else {
		cmd.LogErrorWithStderr("Unknown OS : " + runtime.GOOS)
	}

	// end operations
	t2 := time.Now()
	durationInSeconds := t2.Sub(t1).Milliseconds() / 1000
	cmd.LogInfoWithStdout(fmt.Sprintf("Disk scan found %v files in %v folders (%v files and %v folders can't read) (total %v seconds)", fileCounter, folderCounter, fileErrorCounter, folderErrorCounter, durationInSeconds))
	close(diskToSocketChannelQueue)
}
