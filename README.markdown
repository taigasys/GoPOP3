Introduction
============

GoPOP3 implements the POP3 protocol specified in RFC 1939


POP3 Commands
=============

Currently missing POP3 optional commands:

- TOP
- UIDL
- APOP

Installation & Usag
===================


In order to install the package, run the following command:

	go get github.com/d3xter/GoPOP3



Usage
=====

First of all, you have to import the package

	import (
		"github.com/d3xter/GoPOP3"
	)


Then you can create a connection and try to authenticate:

	plainAuth := pop3.CreatePlainAuthentication("user", "pass")
	client, dialErr := pop3.Dial("127.0.01:3412")
	authErr := client.Authenticate(plainAuth)
