package main

import (
	"flag"
	"http"
	"log"
	"template"
	"websocket"
)

type message struct {
	text []byte
}

type client struct {
	messages chan message
	ws       *websocket.Conn
}

type hub struct {
	messages     chan message
	subscribes   chan *client
	unsubscribes chan *client
}

func (h *hub) run() {
	clients := make(map[*client]bool)
	for {
		select {
		case c := <-h.subscribes:
			clients[c] = true
		case c := <-h.unsubscribes:
			if clients[c] {
				delete(clients, c)
				close(c.messages)
			}
		case m := <-h.messages:
			for c := range clients {
				select {
				case c.messages <- m:
				default:
					delete(clients, c)
					close(c.messages)
					go c.ws.Close()
				}
			}
		}
	}
}

func (c *client) shutdown() {
	h.unsubscribes <- c
	c.ws.Close()
}

func (c *client) reader() {
	defer c.shutdown()
	for {
		text := make([]byte, 256)
		n, err := c.ws.Read(text)
		if err != nil {
			break
		}
		h.messages <- message{text[:n]}
	}
}

func (c *client) writer() {
	defer c.shutdown()
	for m := range c.messages {
		_, err := c.ws.Write(m.text)
		if err != nil {
			break
		}
	}
}

var addr = flag.String("addr", ":8080", "http service address")

var h = hub{make(chan message), make(chan *client), make(chan *client)}

func main() {
	flag.Parse()
	go h.run()
	http.HandleFunc("/", homeHandler)
	http.Handle("/ws", websocket.Handler(wsHandler))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func wsHandler(ws *websocket.Conn) {
	c := &client{make(chan message, 256), ws}
	h.subscribes <- c
	go c.writer()
	c.reader()
}

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}

var homeTempl = template.Must(template.New("home").Parse(`
<html>
<head>
<title>Chat Example</title>
<script type="text/javascript" src="http://ajax.googleapis.com/ajax/libs/jquery/1.4.2/jquery.min.js"></script>
<script type="text/javascript">
    $(function() {

    var conn;
    var msg = $("#msg");
    var log = $("#log");

    function appendLog(msg) {
        var d = log[0]
        var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
        msg.appendTo(log)
        if (doScroll) {
            d.scrollTop = d.scrollHeight - d.clientHeight;
        }
    }

    $("#form").submit(function() {
        if (!conn) {
            return false;
        }
        if (!msg.val()) {
            return false;
        }
        conn.send(msg.val());
        msg.val("");
        return false
    });

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://{{$}}/ws");
        conn.onclose = function(evt) {
            appendLog($("<div><b>Connection closed.</b></div>"))
        }
        conn.onmessage = function(evt) {
            appendLog($("<div/>").text(evt.data))
        }
    } else {
        appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
    }
    });
</script>
<style type="text/css">
html {
    overflow: hidden;
}

body {
    overflow: hidden;
    padding: 0;
    margin: 0;
    width: 100%;
    height: 100%;
    background: gray;
}

#log {
    background: white;
    margin: 0;
    padding: 0.5em 0.5em 0.5em 0.5em;
    position: absolute;
    top: 0.5em;
    left: 0.5em;
    right: 0.5em;
    bottom: 3em;
    overflow: auto;
}

#form {
    padding: 0 0.5em 0 0.5em;
    margin: 0;
    position: absolute;
    bottom: 1em;
    left: 0px;
    width: 100%;
    overflow: hidden;
}

</style>
</head>
<body>
<div id="log"></div>
<form id="form">
    <input type="submit" value="Send" />
    <input type="text" id="msg" size="64"/>
</form>
</body>
</html> `))
