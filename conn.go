package main

import (
	"code.google.com/p/go.net/websocket"
	"time"
	"fmt"
)

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan string
}

func (c *connection) writer() {
	for message := range c.send {
		err := websocket.Message.Send(c.ws, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

func (c *connection) timeTeller() {
	t := time.Tick(1 * time.Second)
	for now := range t{
		fmt.Println(now)
		h.broadcast <- now.String()
	}
}

func wsHandler(ws *websocket.Conn) {
	c := &connection{send: make(chan string, 256), ws: ws}
	h.register <- c
	defer func() { h.unregister <- c }()
	go c.writer()
	go c.timeTeller()
}
