run: clean example
	./example
 
include $(GOROOT)/src/Make.inc

TARG=example
GOFILES=\
    main.go\
    hub.go\
    conn.go

include $(GOROOT)/src/Make.cmd
