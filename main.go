package main

import (
	"os"
	"log"
	"fmt"
	"flag"
	"strings"
	"net/http"
	"text/template"
	"github.com/howeyc/fsnotify"
	"code.google.com/p/go.net/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("home.html"))

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}

func MonitorMultiFiles(names []string, out chan []string,
	h hub) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	for _, name := range names {
		result, _ := ReadLastNLines(name, 10)
		PrintMultiLines(result)
		MonitorFile(name, out, watcher)
	}
	for{
		select{
		case lines := <-out:
			content := strings.Join(lines, "\n")
			fmt.Print(content)
			h.broadcast <- content
		}
	}
	watcher.Close()
}

var usage = func() {
	fmt.Fprintf(os.Stderr,
		"Usage: gowebtail [FILE]...\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if len(flag.Args()) < 1{
		usage()
	}
	out := make(chan []string)
	go h.run()
	go MonitorMultiFiles(flag.Args(), out, h)
	http.HandleFunc("/", homeHandler)
	http.Handle("/ws", websocket.Handler(wsHandler))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
