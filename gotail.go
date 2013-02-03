package main

import (
	"bytes"
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"strings"
)

func LineCount(lines []byte) int {
	return bytes.Count(lines, []byte("\n"))
}

func GetFileSize(filename string) int {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	return int(fileInfo.Size())
}

func ByteArrayToMultiLines(bytes []byte) []string {
	lines := string(bytes)
	return strings.Split(lines, "\n")
}

func ReadNBytes(filename string, start int, end int) []byte {
	fh, _ := os.Open(filename)
	defer fh.Close()
	fh.Seek(int64(start), 0)
	size := end - start + 1
	buff := make([]byte, size)
	fh.Read(buff)
	return buff
}

func ReadLastNLines(name string, n int) ([]string, error) {
	curr := GetFileSize(name)
	var end int
	count := n
	result := make([]byte, n)
	for count > 0 && curr != 0 {
		curr -= n
		end = curr + n - 1
		if curr < 0 {
			curr = 0
		}
		buff := ReadNBytes(name, curr, end)
		result = append(buff, result...)
		count -= LineCount(buff)
	}
	return ByteArrayToMultiLines(result), nil
}

func MonitorFile(filename string, out chan []string,
	watcher *fsnotify.Watcher) {
	size := GetFileSize(filename)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsModify() {
					NewSize := GetFileSize(ev.Name)
					if NewSize <= size {
						MonitorFile(ev.Name, out, watcher)
						return
					}
					content := ReadNBytes(ev.Name, size,
						NewSize-1)
					size = NewSize
					out <- ByteArrayToMultiLines(content)
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()
	err := watcher.Watch(filename)
	if err != nil {
		log.Fatal(err)
	}
}

func PrintMultiLines(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}
