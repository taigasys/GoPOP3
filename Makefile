include $(GOROOT)/src/Make.inc

TARG=pop3
GOFILES=\
	pop3.go\
	auth.go\

include $(GOROOT)/src/Make.pkg
