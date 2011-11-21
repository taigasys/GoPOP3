//Copyright 2011, Andreas Sinz
// Use of this source code is governed by the GPLv2
// license that can be found in the LICENSE file.


//Implements the Post Office Protocol 3 as defined in RFC 1939
package pop3


import (
	"net"
	"errors"
	"bufio"
	"strings"
	"strconv"
)

const (
	//Carriage Return + Line Feed
	//CRLF is appended at the end of each commands
	CRLF = "\r\n"

	//POP3-Commands
	USER = "USER"
	PASSWORD = "PASS"
	NOOP = "NOOP"
	RESET = "RSET"
	DELETE = "DELE"
	QUIT = "QUIT"
	STATUS = "STAT"

	//Error messages
	IndexERR = "Index must be greater than zero"
)

type Client struct {
	conn       net.Conn
	stream     *bufio.ReadWriter
	ServerName string
	Greeting  string
}

//Returns a new Client connected to a POP3 server at addr.
//The format of addr is "ip:port"
func Dial(addr string) (client *Client, err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	host := addr[:strings.Index(addr, ":")]
	return NewClient(conn, host)

}

//NewClient returns a new Client using an existing connection
//name is used as the Servername
func NewClient(conn net.Conn, name string) (*Client, error) {
	client := new(Client)

	//Create a new ReadWriter and store it
	stream := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	client.stream = stream

	//Download the greeting from the POP3 server
	msg, err := client.ReadMessage(false)

	if err != nil {
		return nil, err
	}

	client.ServerName = name
	client.Greeting = msg

	return client, nil
}

//WriteMessage sends the message to the POP3 server
func (client *Client) WriteMessage(message string) error {
	if client == nil {
		return errors.New("Connection hasn't been established")
	}

	tmp := message + CRLF
	_, err1 := client.stream.WriteString(tmp)
	client.stream.Flush()

	return err1
}

//ReadMessage reads a single or multiline response from the POP3 server
//It doesnt finish, until it has received a message
func (client *Client) ReadMessage(multiLine bool) (string, error) {

	//Check, whether the client connection has already been
	if client == nil {
		return "", errors.New("Connection hasn't been established")
	}

	//Get first line of the response
	msg, err := client.stream.ReadString('\n')

	if err != nil {
		return "", err
	}

	//Check, whether the response starts with "+OK" or "-ERR", otherwise return an error
	if strings.HasPrefix(msg, "+OK") {
		msg = msg[4:]

		if multiLine {

			for true {
				line, err1 := client.stream.ReadString('\n')

				if err1 != nil {
					return "", err1
				}

				if line == "." + CRLF {
					break;
				}

				msg += line
			}
		}

	} else if strings.HasPrefix(msg, "-ERR") {
		return "", errors.New(msg[5:])
	} else {
		return "", errors.New("Unkown Response received")
	}

	return msg, nil
}

//Authenticates with the Server
//It uses the provided Authenticationtype auth
func (client *Client) Authenticate(auth Auth) (string, error) {
	return "", auth.Authenticate(client)
}

//Sends a "NOOP" command and the server will just reply with a positive repsonse
func (client *Client) Ping() (string, error) {
	client.WriteMessage(NOOP)
	return client.ReadMessage(false)
}

//Messages that have been marked as "deleted" will be unmarked after this command
func (client *Client) Reset() (string, error) {
	client.WriteMessage(RESET)
	return client.ReadMessage(false)
}

//Mark a mail as "deleted"
//All marked mails will be deleted, when you close the connection with "QUIT"
func (client *Client) Delete(index int) (string, error) {
	if index < 0 {
		return "", errors.New(IndexERR)
	}

	client.WriteMessage(DELETE + " " + string(index))
	return client.ReadMessage(false)
}

//Issues the Quit-Command, so the POP3 session enters the UPDATE state
//All mails, which are marked as "deleted", are going to be removed now
func (client *Client) Quit() (string, error) {
	client.WriteMessage(QUIT)
	return client.ReadMessage(false)
}

//Retrieves the count of mails and the size of all those mails in the mailbox
//Mails, which are marked as "deleted", won't show up
func (client *Client) Status() (mailCount, mailBoxSize int, err error) {
	client.WriteMessage(STATUS)

	var response string
	response, err = client.ReadMessage(false)

	if err != nil {
		return
	}

	responseParts := strings.Split(response, " ")

	if len(responseParts) != 2 {
		return -1, -1, errors.New("Unkown Response Received")
	}


	if mailCount, err = strconv.Atoi(responseParts[0]); err != nil {
		return
	}

	if mailBoxSize, err = strconv.Atoi(responseParts[1]); err != nil {
		return
	}

	return
}
