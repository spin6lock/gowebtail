This demo is based on [fsnotify](https://github.com/howeyc/fsnotify)
and [Implementing Chat with WebSockets and Go](http://gary.beagledreams.com/page/go-websocket-chat.html).
It's the web version of my [gotail](https://github.com/spin6lock/gotail).
The home.html.go is packed by [go-bindata](https://github.com/jteeuwen/go-bindata).
I really enjoy log.io but it's a little bit heavy for me.
So I write my own version with golang. Hope you like it.

Usage:
======
$gowebtail -addr=":8081" test.log

Then you can read it in your browser through http://localhost:8081
